package handlers

import (
	"github.com/minio/minio/pkg/madmin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServerInfo struct {
}

func (tsi testServerInfo) ServerInfo() (madmin.InfoMessage, error) {
	infoMessage := madmin.InfoMessage{
		Mode:         "online",
		DeploymentID: "some-deployment-id",
		Buckets:      madmin.Buckets{},
		Objects:      madmin.Objects{},
		Usage:        madmin.Usage{},
		Services: madmin.Services{
			Vault: madmin.Vault{
				Status: "disabled",
			},
			LDAP: madmin.LDAP{},
		},
	}
	return infoMessage, nil
}

func TestServerInfo(t *testing.T) {
	t.Run("Should create new ServerInfoHandler", func(t *testing.T) {
		serverInfoHandler, err := NewServerInfoHandler(getTestAppConfig())
		assert.Nil(t, err)
		assert.NotNil(t, serverInfoHandler)
	})

	t.Run("Should return JSON with serverInfo", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8080/serverinfo", nil)
		response := httptest.NewRecorder()
		testServerInfo := testServerInfo{}
		serverInfoHandler := ServerInfoHandler{
			serverInfoRetriever: testServerInfo,
		}

		serverInfoHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), "online")
		assert.Contains(t, response.Body.String(), "some-deployment-id")
		assert.Contains(t, response.Body.String(), "disabled")
		assert.Contains(t, response.Body.String(), "buckets")
	})

}
