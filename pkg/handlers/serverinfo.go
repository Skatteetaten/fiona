package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/minio/minio/pkg/madmin"
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/s3"
	"net/http"
)

type serverInfoRetriever interface {
	ServerInfo() (madmin.InfoMessage, error)
}

// ServerInfoHandler fetches information for minio S3 servers
type ServerInfoHandler struct {
	serverInfoRetriever
}

// NewServerInfoHandler is a factory for ServerInfoHandler
func NewServerInfoHandler(config *config.Config) (*ServerInfoHandler, error) {
	admClient, err := s3.NewAdmClient(&config.S3Config)
	if err != nil {
		return nil, err
	}

	return &ServerInfoHandler{admClient}, nil
}

// ServeHTTP handles the requests for ServerInfoHandler
func (serverinfo *ServerInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	infoMessage, err := serverinfo.ServerInfo()
	if err != nil {
		failLogAndResponse(w, "Error calling ServerInfo on S3AdmClient", http.StatusNoContent, err)
		return
	}

	infoJSON, err := json.Marshal(infoMessage)
	if err != nil {
		failLogAndResponse(w, "Unable to parse server info", http.StatusInternalServerError, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", infoJSON)

	logrus.Info("StatusOK: serverinfo")
}
