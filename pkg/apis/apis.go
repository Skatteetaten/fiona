package apis

import (
	"fmt"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v6"
	"github.com/minio/minio/pkg/madmin"
	"github.com/sirupsen/logrus"
	management "github.com/skatteetaten/aurora-management-interface-go"
	"github.com/skatteetaten/aurora-management-interface-go/env"
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/handlers"
	"github.com/skatteetaten/fiona/pkg/handlers/healthcheck"
	"net/http"
)

// InitAPI initializes API with routing
func InitAPI(config *config.Config, adminClient *madmin.AdminClient, minioClient *minio.Client) error {

	auroraTokenAuthenticator, err := NewAuroraTokenAuthenticator(config.AuroraTokenLocation)
	if err != nil {
		return err
	}

	routeHandler, err := createRouter(config, auroraTokenAuthenticator, adminClient, minioClient)
	if err != nil {
		logrus.Errorf("Error while creating router: %s", err)
		return err
	}

	http.Handle("/", routeHandler)

	return nil
}

func createRouter(config *config.Config, amw AuthMiddleware, adminClient *madmin.AdminClient, minioClient *minio.Client) (http.Handler, error) {

	router := mux.NewRouter()

	if err := addRoutes(router, amw, config, adminClient, minioClient); err != nil {
		return nil, err
	}

	logger := logrus.New()
	logwriter := logger.Writer()
	loggedRouter := ghandlers.LoggingHandler(logwriter, router)

	return loggedRouter, nil
}

func addRoutes(router *mux.Router, amw AuthMiddleware, config *config.Config, adminClient *madmin.AdminClient, minioClient *minio.Client) error {
	listusersHandler := handlers.NewListUsersHandler(adminClient)
	serverinfoHandler := handlers.NewServerInfoHandler(adminClient)
	createuserHandler, err := handlers.NewCreateUserHandler(config, adminClient, minioClient)
	if err != nil {
		return err
	}

	router.HandleFunc("/", roothandler)
	router.Handle("/listusers", amw.Authenticate(listusersHandler)).Methods("GET")
	router.Handle("/serverinfo", amw.Authenticate(serverinfoHandler)).Methods("GET")
	router.Handle("/createuser", amw.Authenticate(createuserHandler)).Methods("POST")

	return nil
}

func roothandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fiona says hi at %s!", r.Host)
}

// InitManagementHandler initializes the management interface with /health and /env endpoints
func InitManagementHandler(admClient *madmin.AdminClient) *management.RoutingHandler {
	managementHandler := management.CreateRoutingHandler()

	fionaHealthRetriever := healthcheck.NewFionaHealthRetriever(admClient)

	fionaEnvRetriever := env.GetDefaultEnvRetriever()
	fionaEnvRetriever.SetKeysToMask([]string{config.FionaDefaultPassword, config.FionaSecretKey, config.FionaAccessKey})

	managementHandler.RouteApplicationHealthRetriever(fionaHealthRetriever)
	managementHandler.RouteApplicationEnvRetriever(fionaEnvRetriever)

	return managementHandler
}
