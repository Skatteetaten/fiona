package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v6"
	"github.com/minio/minio/pkg/madmin"
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/s3"
	"net/http"
)

// CreateAppUserHandler adds an application user for a specified path
type CreateAppUserHandler struct {
	BucketManager s3.BucketManager
	UserManager   s3.UserManager
}

// NewCreateAppUserHandler is a factory for CreateUserHandler
func NewCreateAppUserHandler(config *s3.Config, adminClient *madmin.AdminClient, minioClient *minio.Client) (*CreateAppUserHandler, error) {
	bucketManager := s3.NewMinioBucketManager(config, minioClient)
	userManager := s3.NewMinioUserManager(config, adminClient)
	return &CreateAppUserHandler{
		BucketManager: bucketManager,
		UserManager:   userManager,
	}, nil
}

// ServeHTTP handles the requests for CreateAppUserHandler
func (createappuser *CreateAppUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Header based access info not needed at present
	/*	minioAccess, err := readMinioAccessHeader(r.Header)
		if err != nil {
			failLogAndResponse(w, "Could not read access header", http.StatusBadRequest, err)
			return
		}
		if minioAccess != nil {
			adminClient, err = s3.NewAdmClientForExternalAccess(minioAccess)
			if err != nil {
				failLogAndResponse(w, fmt.Sprintf("Could create minio admin client for incoming %s", MinioAccess), http.StatusBadRequest, err)
				return
			}
		}
		userManager := s3.NewMinioUserManager(createappuser.S3Config, adminClient) */

	createAppUserInput, doneWithError := getCreateAppUserInput(w, r)
	if doneWithError {
		return
	}
	logrus.Debugf("createAppUserInput: %+v", *createAppUserInput)

	bucketExists, err := createappuser.BucketManager.BucketNameExists(createAppUserInput.Bucketname)
	if err != nil {
		failLogAndResponse(w, "Error creating user. Could not verify existing bucket", http.StatusInternalServerError, err)
		return
	}
	if !bucketExists {
		failLogAndResponse(w, "Error creating user", http.StatusUnprocessableEntity, errors.New("Bucket does not exist"))
		return
	}

	createAppUserResult, err := createappuser.UserManager.CreateAppUser(createAppUserInput)
	if err != nil {
		failLogAndResponse(w, fmt.Sprintf("Error creating user for input: %+v", *createAppUserInput), http.StatusInternalServerError, err)
		return
	}
	responseJSON, err := json.Marshal(createAppUserResult)
	if err != nil {
		failLogAndResponse(w, "Failed marshalling result for return, aborted", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "%s", responseJSON)
	logrus.Infof("StatusCreated: createuser %s", createAppUserInput.Username)
}

func getCreateAppUserInput(w http.ResponseWriter, r *http.Request) (*s3.CreateAppUserInput, bool) {
	params := mux.Vars(r)
	var createAppUserInput s3.CreateAppUserInput
	body, err := readRequestBody(r.Body)
	if err != nil {
		failLogAndResponse(w, "Could not read request body", http.StatusBadRequest, err)
		return nil, true
	}
	if err := json.Unmarshal(body, &createAppUserInput); err != nil {
		failLogAndResponse(w, "Could not unmarshal body", http.StatusUnprocessableEntity, err)
		return nil, true
	}
	if bucketname, ok := params["bucketname"]; ok {
		createAppUserInput.Bucketname = bucketname
	}
	if path, ok := params["path"]; ok {
		createAppUserInput.Path = path
	}
	if username, ok := params["username"]; ok {
		createAppUserInput.Username = username
	}
	if createAppUserInput.Path == "" || createAppUserInput.Username == "" || createAppUserInput.Bucketname == "" || len(createAppUserInput.Access) <= 0 {
		failLogAndResponse(w, "Missing required input to create user.", http.StatusBadRequest, err)
		return nil, true
	}
	return &createAppUserInput, false
}
