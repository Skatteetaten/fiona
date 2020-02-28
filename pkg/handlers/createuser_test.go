package handlers

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testUserCreator struct {
}

func (tuc testUserCreator) CreateUser(userName string, path string) (string, error) {
	return "S3userpass", nil
}
func (tuc testUserCreator) MakeSureBucketExists() error {
	return nil
}

func TestCreateUser(t *testing.T) {
	t.Run("Should create new CreateUserHandler", func(t *testing.T) {
		createUserHandler, err := NewCreateUserHandler(getTestAppConfig())
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
		assert.True(t, isJSONString(response.Body.String()))
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

func createTestUserHandler(tuc testUserCreator) CreateUserHandler {
	return CreateUserHandler{
		tuc,
		tuc,
	}
}
