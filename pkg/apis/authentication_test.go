package apis

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type dummyHandler struct {
}

func (dh dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func TestAuthentication(t *testing.T) {
	t.Run("Should initialize AuroraTokenAuthenticator without failing", func(t *testing.T) {
		hook := test.NewGlobal()

		ata, err := NewAuroraTokenAuthenticator("testdata/token")

		assert.Nil(t, err)
		assert.NotNil(t, ata)
		assert.Equal(t, 0, len(hook.Entries))
		assert.Equal(t, "testtoken", ata.auroratoken)

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("Should fail to initialize and panic when empty aurora token", func(t *testing.T) {
		assert.Panics(t, func() { NewAuroraTokenAuthenticator("testdata/emptytoken") })
	})

	t.Run("Should fail to initialize and panic when no aurora token file", func(t *testing.T) {
		assert.Panics(t, func() { NewAuroraTokenAuthenticator("testdata/nofile") })
	})

	t.Run("Should successfully authenticate token", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
		request.Header.Set("Authorization", "aurora-token testtoken")
		response := httptest.NewRecorder()
		hook := test.NewGlobal()
		dummyHandler := dummyHandler{}

		ata := AuroraTokenAuthenticator{
			auroratoken: "testtoken",
		}
		handlerToTest := ata.Authenticate(dummyHandler)
		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
		assert.Equal(t, "Authentication OK", hook.LastEntry().Message)
		assert.Equal(t, http.StatusOK, response.Code)

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("Should fail authentication when wrong token", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
		request.Header.Set("Authorization", "aurora-token wrongtoken")
		response := httptest.NewRecorder()
		hook := test.NewGlobal()
		dummyHandler := dummyHandler{}

		ata := AuroraTokenAuthenticator{
			auroratoken: "testtoken",
		}
		handlerToTest := ata.Authenticate(dummyHandler)
		handlerToTest.ServeHTTP(response, request)

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
		assert.Equal(t, "Authentication failed", hook.LastEntry().Message)
		assert.Equal(t, http.StatusUnauthorized, response.Code)

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})

}
