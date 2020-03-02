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

type userLister interface {
	ListUsers() (map[string]madmin.UserInfo, error)
}

// ListUsersHandler handles listusers
type ListUsersHandler struct {
	userLister
}

// NewListUsersHandler is a factory for ListUsersHandlers
func NewListUsersHandler(config *config.Config) (*ListUsersHandler, error) {
	admClient, err := s3.NewAdmClient(&config.S3Config)
	if err != nil {
		return nil, err
	}

	return &ListUsersHandler{admClient}, nil
}

// ServeHTTP handles the requests for ListUsersHandler
func (listusers *ListUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	users, err := listusers.ListUsers()
	if err != nil {
		failLogAndResponse(w, "Error calling ListUsers on S3AdmClient", http.StatusInternalServerError, err)
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		failLogAndResponse(w, "Could not create json for users", http.StatusInternalServerError, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", usersJSON)

	logrus.Info("StatusOK: listusers")
}
