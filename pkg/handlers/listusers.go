package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/minio/minio/pkg/madmin"
	"github.com/secure-io/sio-go"
	"github.com/sirupsen/logrus"
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
func NewListUsersHandler(admClient *madmin.AdminClient) *ListUsersHandler {
	return &ListUsersHandler{admClient}
}

// ServeHTTP handles the requests for ListUsersHandler
func (listusers *ListUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	users, err := listusers.ListUsers()
	if err != nil {
		message := "Error calling ListUsers on S3AdmClient"
		if err == sio.NotAuthentic {
			message += ". Decryption failed. Is the secret key incorrect?"
		}
		failLogAndResponse(w, message, http.StatusInternalServerError, err)
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
