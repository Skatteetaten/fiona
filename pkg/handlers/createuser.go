package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/s3"
	"net/http"
)

// CreateUserHandler adds a user to the minio S3 server
type CreateUserHandler struct {
	BucketManager s3.BucketManager
	UserManager   s3.UserManager
}

// NewCreateUserHandler is a factory for CreateUserHandler
func NewCreateUserHandler(config *config.Config) (*CreateUserHandler, error) {
	bucketManager, err := s3.NewMinioBucketManager(&config.S3Config)
	if err != nil {
		return nil, err
	}
	userManager, err := s3.NewMinioUserManager(&config.S3Config)
	if err != nil {
		return nil, err
	}

	return &CreateUserHandler{
		BucketManager: bucketManager,
		UserManager:   userManager,
	}, nil
}

// ServeHTTP handles the requests for CreateUserHandler
func (createuser *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var user User
	body, err := readRequestBody(r.Body)
	if err != nil {
		failLogAndResponse(w, "Could not read request body", http.StatusBadRequest, err)
		return
	}
	if err := json.Unmarshal(body, &user); err != nil {
		failLogAndResponse(w, "Could not unmarshal body", http.StatusUnprocessableEntity, err)
		return
	}
	if user.Basepath == "" || user.Username == "" {
		failLogAndResponse(w, "Missing required input", http.StatusForbidden, err)
		return
	}

	if err := createuser.BucketManager.MakeSureBucketExists(); err != nil {
		failLogAndResponse(w, "Error when making sure bucket exists", http.StatusInternalServerError, err)
		return
	}

	result, err := createuser.UserManager.CreateUser(user.Username, user.Basepath)
	if err != nil {
		failLogAndResponse(w, "Error creating user", http.StatusInternalServerError, err)
		return
	}
	responseJson, err := json.Marshal(map[string]string{
		"secretKey":       result.SecretKey,
		"serviceEndpoint": result.ServiceEndpoint,
		"bucket":          result.Bucket,
		"bucketRegion":    result.BucketRegion,
	})
	if err != nil {
		failLogAndResponse(w, "Failed marshalling secret for return, aborted", http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "%s", responseJson)
	logrus.Infof("StatusCreated: createuser %s", user.Username)
}
