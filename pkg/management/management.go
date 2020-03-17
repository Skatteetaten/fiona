package management

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/management/env"
	"github.com/skatteetaten/fiona/pkg/management/health"
	"net/http"
)

// DefaultPort is the default port for the management interface
const DefaultPort = "8081"

// RoutingHandler is a router for the management interface
type RoutingHandler struct {
	port           string
	managementMux  *http.ServeMux
	managementSpec *ManagementinterfaceSpec
}

// CreateRoutingHandler creates a router to handle the management interface requests on default port (8081)
func CreateRoutingHandler() *RoutingHandler {
	return CreateRoutingHandlerForPort(DefaultPort)
}

// CreateRoutingHandlerForPort creates a router to handle the management interface requests on a specific port
func CreateRoutingHandlerForPort(port string) *RoutingHandler {
	managementSpec := NewManagementinterfaceSpec()
	managementMux := http.NewServeMux()

	mrh := RoutingHandler{
		port:           port,
		managementMux:  managementMux,
		managementSpec: managementSpec,
	}
	mrh.RouteEndPointToHandlerFunc(Management, mrh.managementHandler)

	return &mrh
}

func (mrh RoutingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mrh.managementMux.ServeHTTP(w, r)
}

// StartHTTPListener starts a http.ListenAndServe process for the router on the configured port
func (mrh RoutingHandler) StartHTTPListener() {
	managementPort := fmt.Sprintf(":%s", mrh.port)
	go http.ListenAndServe(managementPort, mrh)
}

// RouteApplicationHealthRetriever routes the Health endpoint to get status from the specified ApplicationHealthRetriever
func (mrh RoutingHandler) RouteApplicationHealthRetriever(healthRetriever health.ApplicationHealthRetriever) {
	appHealthHandler := health.ApplicationHealthHandler{healthRetriever}
	mrh.RouteEndPointToHandlerFunc(Health, appHealthHandler.ServeHTTP)
}

// RouteApplicationEnvRetriever routes the Health endpoint to get environment variables from the specified ApplicationHealthRetriever
func (mrh RoutingHandler) RouteApplicationEnvRetriever(envRetriever env.ApplicationEnvRetriever) {
	appEnvHandler := env.ApplicationEnvHandler{ApplicationEnvRetriever: envRetriever}
	mrh.RouteEndPointToHandlerFunc(Env, appEnvHandler.ServeHTTP)
}

// RouteEndPointToHandlerFunc routes endpoint to a handlerfunc
func (mrh RoutingHandler) RouteEndPointToHandlerFunc(endPointType EndPointType, handlerfunc func(http.ResponseWriter, *http.Request)) error {
	usedefaultpathstring := ""
	return mrh.RouteEndPointToHandlerFuncWithPath(endPointType, usedefaultpathstring, handlerfunc)
}

// RouteEndPointToHandlerFuncWithPath routes endpoint to a handlerfunc on a specified, non-default path
func (mrh RoutingHandler) RouteEndPointToHandlerFuncWithPath(endPointType EndPointType, path string, handlerfunc func(http.ResponseWriter, *http.Request)) error {
	endpoint, err := newEndPoint(endPointType, handlerfunc)
	if err != nil {
		return err
	}

	if path != "" {
		endpoint.setPath(path)
	}
	mrh.managementSpec.mapEndpoint(*endpoint)
	mrh.managementMux.HandleFunc(endpoint.path, handlerfunc)

	return nil
}

func (mrh RoutingHandler) managementHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	managementJSON, err := mrh.managementSpec.createManagementJSON(r.Host)

	if err != nil {
		message := "Error while creating JSON for management interface"
		logrus.Errorf("%s: %s", message, err)
		w.WriteHeader(http.StatusInternalServerError)
		responseJSON, _ := json.Marshal(map[string]string{
			"error": message,
			"cause": fmt.Sprintf("%v", err),
		})
		_, _ = fmt.Fprintf(w, "%s", responseJSON)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", managementJSON)
}
