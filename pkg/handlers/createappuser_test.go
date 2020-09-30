package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/skatteetaten/fiona/pkg/s3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const validtestbucketname = "testbucketname"

type testAppUserCreator struct {
}

func (tuc testAppUserCreator) CreateAppUser(createAppUserInput *s3.CreateAppUserInput) (*s3.CreateAppUserResult, error) {
	mockCreateAppUserResult := s3.CreateAppUserResult{
		AccessKey: "testuser",
		SecretKey: "S3userpass",
		HostURL:   "http://localhost:9000",
	}
	return &mockCreateAppUserResult, nil
}
func (tuc testAppUserCreator) CreateUser(userName string, path string) (*s3.CreateUserResult, error) {
	return nil, nil
}
func (tuc testAppUserCreator) BucketNameExists(bucketName string) (bool, error) {
	return (bucketName == validtestbucketname), nil
}
func (tuc testAppUserCreator) MakeSureBucketExists() error {
	return nil
}

func TestCreateAppUser(t *testing.T) {
	t.Run("Should create new CreateAppUserHandler", func(t *testing.T) {
		dummyAdmClient, _ := s3.NewAdmClient(&getTestAppConfig().S3Config)
		dummyClient, _ := s3.NewClient(&getTestAppConfig().S3Config)
		createAppUserHandler, err := NewCreateAppUserHandler(&getTestAppConfig().S3Config, dummyAdmClient, dummyClient)
		assert.Nil(t, err)
		assert.NotNil(t, createAppUserHandler)
	})

	t.Run("Should 'create app user' without failing (happy test)", func(t *testing.T) {
		reader := strings.NewReader("{\"username\":\"testuser\", \"access\":[\"READ\", \"WRITE\", \"DELETE\"]}")
		request, _ := http.NewRequest("POST", "http://localhost:8080/buckets/testbucketname/paths/testpath/userprofiles/", reader)
		request = mux.SetURLVars(request, map[string]string{"bucketname": validtestbucketname, "path": "testpath"})
		response := httptest.NewRecorder()
		testAppUserCreator := testAppUserCreator{}
		createAppUserHandler := createTestAppUserHandler(testAppUserCreator)

		createAppUserHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusCreated, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), getTestAppConfig().S3Config.DefaultUserpass)
	})

	t.Run("Should fail to create user when body is not valid JSON", func(t *testing.T) {
		reader := strings.NewReader("{\"Not valid JSON\"}")
		request, _ := http.NewRequest("POST", "http://localhost:8080/buckets/testbucketname/paths/testpath/userprofiles/", reader)
		request = mux.SetURLVars(request, map[string]string{"bucketname": validtestbucketname, "path": "testpath"})
		response := httptest.NewRecorder()
		testAppUserCreator := testAppUserCreator{}
		createUserHandler := createTestAppUserHandler(testAppUserCreator)

		hook := test.NewGlobal()

		createUserHandler.ServeHTTP(response, request)

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		assert.Contains(t, hook.LastEntry().Message, "Could not unmarshal body")
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("Should fail to create user when bucket does not exist", func(t *testing.T) {
		reader := strings.NewReader("{\"username\":\"testuser\", \"access\":[\"READ\", \"WRITE\", \"DELETE\"]}")
		request, _ := http.NewRequest("POST", "http://localhost:8080/buckets/testbucketname/paths/testpath/userprofiles/", reader)
		request = mux.SetURLVars(request, map[string]string{"bucketname": "nonexistingbucket", "path": "testpath"})
		response := httptest.NewRecorder()
		testAppUserCreator := testAppUserCreator{}
		createUserHandler := createTestAppUserHandler(testAppUserCreator)

		hook := test.NewGlobal()

		createUserHandler.ServeHTTP(response, request)

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		assert.Contains(t, hook.LastEntry().Message, "Bucket does not exist")
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})
}

func createTestAppUserHandler(testAppUserCreator testAppUserCreator) CreateAppUserHandler {
	return CreateAppUserHandler{
		BucketManager: testAppUserCreator,
		UserManager:   testAppUserCreator,
	}
}
