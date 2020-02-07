package apis

import (
	"github.com/gorilla/mux"
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/handlers"
	"net/http"
)

// Web server
type Web struct {
	Config *config.Config
}

// InitAPI initializes API with routing
func (w *Web) InitAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", handlers.HealthCheckHandler).Methods("GET")
	http.Handle("/", router)
}
