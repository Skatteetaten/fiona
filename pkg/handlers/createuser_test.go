package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/skatteetaten/fiona/pkg/s3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testUserCreator struct {
}

func (tuc testUserCreator) CreateUser(userName string, path string) (*s3.CreateUserResult, error) {
	mockCreateUserResult := s3.CreateUserResult{
		SecretKey:       "S3userpass",
		ServiceEndpoint: "http://localhost:9000",
		Bucket:          "utv",
		BucketRegion:    "us-east-1",
	}
	return &mockCreateUserResult, nil
}
func (tuc testUserCreator) CreateAppUser(createAppUserInput *s3.CreateAppUserInput) (*s3.CreateAppUserResult, error) {
	return nil, nil
}
func (tuc testUserCreator) MakeSureBucketExists() error {
	return nil
}
func (tuc testUserCreator) BucketNameExists(bucketName string) (bool, error) {
	return true, nil
}

func TestCreateUser(t *testing.T) {
	t.Run("Should create new CreateUserHandler", func(t *testing.T) {
		dummyAdmClient, _ := s3.NewAdmClient(&getTestAppConfig().S3Config)
		dummyClient, _ := s3.NewClient(&getTestAppConfig().S3Config)
		createUserHandler, err := NewCreateUserHandler(getTestAppConfig(), dummyAdmClient, dummyClient)
		assert.Nil(t, err)
		assert.NotNil(t, createUserHandler)
	})

	t.Run("Should 'create user' without failing (happy test)", func(t *testing.T) {
		reader := strings.NewReader("{\"user\":\"testuser\", \"path\":\"testpathrandom\"}")
		request, _ := http.NewRequest("POST", "http://localhost:8080/createuser", reader)
		response := httptest.NewRecorder()
		testUserCreator := testUserCreator{}
		createUserHandler := createTestUserHandler(testUserCreator)

		createUserHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), getTestAppConfig().S3Config.DefaultUserpass)
	})

	t.Run("Should fail to create user when body is not valid JSON", func(t *testing.T) {
		reader := strings.NewReader("{\"Not valid JSON\"}")
		request, _ := http.NewRequest("POST", "http://localhost:8080/createuser", reader)
		response := httptest.NewRecorder()
		testUserCreator := testUserCreator{}
		createUserHandler := createTestUserHandler(testUserCreator)

		hook := test.NewGlobal()

		createUserHandler.ServeHTTP(response, request)

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		assert.Contains(t, hook.LastEntry().Message, "Could not unmarshal body")
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})
}

func createTestUserHandler(testmanager testUserCreator) CreateUserHandler {
	return CreateUserHandler{
		testmanager,
		testmanager,
	}
}
