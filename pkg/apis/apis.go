package apis

import (
	"fmt"
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
	router.HandleFunc("/", roothandler)
	router.HandleFunc("/healthcheck", handlers.HealthCheckHandler).Methods("GET")

	http.Handle("/", router)
}

func roothandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fiona says hi at %s!", r.Host)
}
