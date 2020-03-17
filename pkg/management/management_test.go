package management

import (
	"fmt"
	"github.com/skatteetaten/fiona/pkg/management/env"
	"github.com/skatteetaten/fiona/pkg/management/health"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRoutingHandler(t *testing.T) {
	t.Run("Should create management routing handler", func(t *testing.T) {
		managementHandler := CreateRoutingHandler()
		assert.NotNil(t, managementHandler)
		assert.NotNil(t, managementHandler.managementSpec)
		assert.NotNil(t, managementHandler.managementMux)
		assert.Equal(t, DefaultPort, managementHandler.port)
	})
	t.Run("Should set up routing", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8081/health", nil)

		testhandlerfunc := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/vnd.spring-boot.actuator.v3+json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprintf(w, "%s", `{"status": "UNKNOWN"}`)
		}

		managementHandler := CreateRoutingHandler()
		assert.NotNil(t, managementHandler)
		managementHandler.RouteEndPointToHandlerFunc(Health, testhandlerfunc)

		// Verify management spec
		assert.Equal(t, 1, len(managementHandler.managementSpec.endpoints))
		assert.Equal(t, Health, managementHandler.managementSpec.endpoints[Health].endpointid)
		handlerfunc := managementHandler.managementSpec.endpoints[Health].handlerfunc
		response := httptest.NewRecorder()
		handlerfunc(response, request)
		assert.Contains(t, response.Body.String(), "UNKNOWN")

		// Verify mapping in mux
		handler, pattern := managementHandler.managementMux.Handler(request)
		assert.Equal(t, defaultHealthPath, pattern)
		response = httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		assert.Contains(t, response.Body.String(), "UNKNOWN")
	})
	t.Run("Should set up routing with default retrievers for health and env", func(t *testing.T) {
		managementHandler := CreateRoutingHandler()
		defaultHealthRetriever := health.GetDefaultHealthRetriever()
		defaultEnvRetriever := env.GetDefaultEnvRetriever()

		managementHandler.RouteApplicationHealthRetriever(defaultHealthRetriever)
		managementHandler.RouteApplicationEnvRetriever(defaultEnvRetriever)

		// Verify management spec
		assert.Equal(t, 2, len(managementHandler.managementSpec.endpoints))
		assert.Equal(t, Health, managementHandler.managementSpec.endpoints[Health].endpointid)
		assert.Equal(t, Env, managementHandler.managementSpec.endpoints[Env].endpointid)

		// Verify mapping in mux
		healthrequest, _ := http.NewRequest("GET", "http://localhost:8081/health", nil)
		_, healthpattern := managementHandler.managementMux.Handler(healthrequest)
		assert.Equal(t, defaultHealthPath, healthpattern)

		envrequest, _ := http.NewRequest("GET", "http://localhost:8081/env", nil)
		_, envpattern := managementHandler.managementMux.Handler(envrequest)
		assert.Equal(t, defaultEnvPath, envpattern)
	})
}
