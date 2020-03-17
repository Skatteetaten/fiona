package env

import (
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestApplicationEnvHandler_ServeHTTP(t *testing.T) {
	t.Run("Should handle nil case from Envfunc", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8081/env", nil)
		response := httptest.NewRecorder()
		applicationEnv := ApplicationEnv{
			ActiveProfiles:  nil,
			PropertySources: nil,
		}
		appenvhandler := ApplicationEnvHandler{
			ApplicationEnvRetriever: testApplicationEnvRetriever{
				testFunc: func() *ApplicationEnv {
					return &applicationEnv
				},
			},
		}

		appenvhandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), "activeProfiles")
		assert.Contains(t, response.Body.String(), "propertySources")
	})

	t.Run("Should handle full happy case from Envfunc", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8081/env", nil)
		response := httptest.NewRecorder()

		propertyMasker := PropertyMasker{}

		appenvhandler := ApplicationEnvHandler{
			ApplicationEnvRetriever: testApplicationEnvRetriever{
				testFunc: func() *ApplicationEnv {
					key1 := "property1"
					key2 := "propsecret"
					value1 := "value1"
					value2 := "value2_secret"
					properties := make(map[string]PropertyValue)
					properties[key1] = propertyMasker.GetPropertyValue(key1, value1)
					properties[key2] = propertyMasker.GetPropertyValue(key2, value2)
					propertysource1 := PropertySource{
						Name:       "systemEnvironment",
						Properties: properties,
					}
					propertySources := []PropertySource{propertysource1}
					applicationEnv := ApplicationEnv{
						ActiveProfiles:  []string{"openshift"},
						PropertySources: propertySources,
					}
					return &applicationEnv
				},
			},
		}

		appenvhandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), "openshift")
		assert.Contains(t, response.Body.String(), "systemEnvironment")
		assert.Contains(t, response.Body.String(), "property1")
		assert.Contains(t, response.Body.String(), "propsecret")
		assert.Contains(t, response.Body.String(), "value1")
		assert.Contains(t, response.Body.String(), "***")
		assert.NotContains(t, response.Body.String(), "value2_secret")
	})

}

func TestPropertyMasker_GetPropertyValue(t *testing.T) {
	t.Run("Should handle nil case from Envfunc", func(t *testing.T) {
		propertymasker := PropertyMasker{}
		propertymasker.SetKeysToMask([]string{"Superhemmelig"})

		assert.True(t, propertymasker.GetPropertyValue("Open", "Open").Value == "Open")
		assert.True(t, propertymasker.GetPropertyValue("KEY", "Secret").Value == "***")
		assert.True(t, propertymasker.GetPropertyValue("key", "Secret").Value == "***")
		assert.True(t, propertymasker.GetPropertyValue("superhemmelig", "Secret").Value == "***")
		assert.True(t, propertymasker.GetPropertyValue("asecretproperty", "Secret").Value == "***")
	})
}

func TestDefaultApplicationEnvHandler_ServeHTTP(t *testing.T) {
	t.Run("Should work without error (happytest)", func(t *testing.T) {
		fdp, fdpexists := os.LookupEnv(config.FionaDefaultPassword)
		os.Setenv(config.FionaDefaultPassword, "testdefaultpass")

		request, _ := http.NewRequest("GET", "http://localhost:8081/env", nil)
		response := httptest.NewRecorder()
		appEnvHandler := DefaultApplicationEnvHandler()
		appEnvHandler.SetKeysToMask([]string{config.FionaDefaultPassword, config.FionaSecretKey, config.FionaAccessKey})

		appEnvHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
		assert.True(t, isJSON(response.Body.String()))
		assert.Contains(t, response.Body.String(), "activeProfiles")
		assert.Contains(t, response.Body.String(), "propertySources")
		assert.Contains(t, response.Body.String(), "systemEnvironment")
		assert.Contains(t, response.Body.String(), "***")
		assert.NotContains(t, response.Body.String(), "testdefaultpass")

		if fdpexists {
			os.Setenv(config.FionaDefaultPassword, fdp)
		} else {
			os.Unsetenv(config.FionaDefaultPassword)
		}
	})
}

type testApplicationEnvRetriever struct {
	testFunc func() *ApplicationEnv
}

func (taer testApplicationEnvRetriever) GetApplicationEnv() *ApplicationEnv {
	return taer.testFunc()
}
func (taer testApplicationEnvRetriever) SetKeysToMask(keysToMask []string) {
	/* dummy implementation */
}
