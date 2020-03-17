package health

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHealthHandlerFunc_ServeHTTP(t *testing.T) {
	t.Run("Should work without error (happytest)", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "http://localhost:8081/health", nil)
		response := httptest.NewRecorder()

		appHealthHandler := ApplicationHealthHandler{ApplicationHealthRetriever: DefaultHealthRetriever{}}
		appHealthHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), "status")
		assert.Contains(t, response.Body.String(), "UP")
	})
	t.Run("Should work with status DOWN", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "http://localhost:8081/health", nil)
		response := httptest.NewRecorder()

		appHealthHandler := ApplicationHealthHandler{
			ApplicationHealthRetriever: testApplicationHealthRetriever{
				response: &ApplicationHealth{Status: Down},
			},
		}

		appHealthHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusServiceUnavailable, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), "status")
		assert.Contains(t, response.Body.String(), "DOWN")
	})
	t.Run("Should return 500 (Internal Server Error) with empty status", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "http://localhost:8081/health", nil)
		response := httptest.NewRecorder()

		appHealthHandler := ApplicationHealthHandler{
			ApplicationHealthRetriever: testApplicationHealthRetriever{
				response: &ApplicationHealth{},
			},
		}

		appHealthHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, isJSON(response.Body.String()))
	})
	t.Run("Should fail on nil", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8081/health", nil)
		response := httptest.NewRecorder()

		appHealthHandler := ApplicationHealthHandler{
			ApplicationHealthRetriever: testApplicationHealthRetriever{
				response: nil,
			},
		}

		appHealthHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), "cause")
		assert.Contains(t, response.Body.String(), "error")
	})

}

type testApplicationHealthRetriever struct {
	response *ApplicationHealth
}

func (tahr testApplicationHealthRetriever) GetApplicationHealth() *ApplicationHealth {
	return tahr.response
}
