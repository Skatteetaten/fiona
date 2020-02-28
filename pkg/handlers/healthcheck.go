package handlers

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

// HealthCheckHandler checks for application life
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our S3 backend etc.
	io.WriteString(w, `{"alive": true}`)

	logrus.Info("StatusOK: healthcheck")
}
