package apis

import (
	"fmt"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/handlers"
	"net/http"
)

// InitAPI initializes API with routing
func InitAPI(config *config.Config) error {

	auroraTokenAuthenticator, err := NewAuroraTokenAuthenticator(config.AuroraTokenLocation)
	if err != nil {
		logrus.Errorf("Error while creating initializing AuroraTokenAuthenticator : %s", err)
		return err
	}

	routeHandler, err := createRouter(config, auroraTokenAuthenticator)
	if err != nil {
		logrus.Errorf("Error while creating router: %s", err)
		return err
	}

	http.Handle("/", routeHandler)

	return nil
}

func createRouter(config *config.Config, amw AuthMiddleware) (http.Handler, error) {

	router := mux.NewRouter()

	if err := addRoutes(router, amw, config); err != nil {
		return nil, err
	}

	logger := logrus.New()
	logwriter := logger.Writer()
	loggedRouter := ghandlers.LoggingHandler(logwriter, router)

	return loggedRouter, nil
}

func addRoutes(router *mux.Router, amw AuthMiddleware, config *config.Config) error {
	listusersHandler, err := handlers.NewListUsersHandler(config)
	if err != nil {
		return err
	}
	serverinfoHandler, err := handlers.NewServerInfoHandler(config)
	if err != nil {
		return err
	}
	createuserHandler, err := handlers.NewCreateUserHandler(config)
	if err != nil {
		return err
	}

	router.HandleFunc("/", roothandler)
	router.HandleFunc("/healthcheck", handlers.HealthCheckHandler).Methods("GET")
	router.Handle("/listusers", amw.Authenticate(listusersHandler)).Methods("GET")
	router.Handle("/serverinfo", amw.Authenticate(serverinfoHandler)).Methods("GET")
	router.Handle("/createuser", amw.Authenticate(createuserHandler)).Methods("POST")

	return nil
}

func roothandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fiona says hi at %s!", r.Host)
}
