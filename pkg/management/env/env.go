package env

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

var defaultKeysToSanitize = [...]string{"password", "secret", "key", "token", "credential", "vcap_services", "sun.java.command"}

// ApplicationEnvRetriever is an interface for retrieving a structure of environment variables for management-interface use
type ApplicationEnvRetriever interface {
	GetApplicationEnv() *ApplicationEnv
	SetKeysToMask(keysToMask []string)
}

// ApplicationEnv is a structure for standardized environment variables response from management interface at application level
type ApplicationEnv struct {
	ActiveProfiles  []string         `json:"activeProfiles"`
	PropertySources []PropertySource `json:"propertySources"`
}

// PropertySource is a structure of environment properties from a specific source
type PropertySource struct {
	Name       string                   `json:"name"`
	Properties map[string]PropertyValue `json:"properties"`
}

// PropertyValue is a standard structure for a property value
type PropertyValue struct {
	Value  string `json:"value"`
	Origin string `json:"origin,omitempty"`
}

// ApplicationEnvHandler fetches an ApplicationEnv structure from the application and parses it for a proper http response
type ApplicationEnvHandler struct {
	ApplicationEnvRetriever ApplicationEnvRetriever
}

func (aeh *ApplicationEnvHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	envResponse := aeh.ApplicationEnvRetriever.GetApplicationEnv()

	envResponseJSON, err := json.Marshal(envResponse)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		message := "Error while creating JSON for env response"
		logrus.Errorf("%s: %s", message, err)
		w.WriteHeader(http.StatusInternalServerError)
		responseJSON, _ := json.Marshal(map[string]string{
			"error": message,
			"cause": fmt.Sprintf("%v", err),
		})
		_, _ = fmt.Fprintf(w, "%s", responseJSON)
		return
	}
	logrus.Debugf("Env response output:\n%s\n", string(envResponseJSON))

	w.Header().Set("Content-Type", "application/vnd.spring-boot.actuator.v3+json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", envResponseJSON)
}

// SetKeysToMask specifies an array of environvent variable names that will be masked regardless of standard masking
func (aeh *ApplicationEnvHandler) SetKeysToMask(keysToMask []string) {
	aeh.ApplicationEnvRetriever.SetKeysToMask(keysToMask)
}

// DefaultApplicationEnvHandler provides a simple, standardized handler for the env management endpoint
func DefaultApplicationEnvHandler() *ApplicationEnvHandler {
	applicationEnvHandler := ApplicationEnvHandler{
		ApplicationEnvRetriever: DefaultEnvRetriever{
			PropertyMasker: &PropertyMasker{keysToMask: nil},
		},
	}

	return &applicationEnvHandler
}

// GetDefaultEnvRetriever gets a DefaultEnvRetriever with default initialization
func GetDefaultEnvRetriever() DefaultEnvRetriever {
	return DefaultEnvRetriever{PropertyMasker: &PropertyMasker{}}
}

// DefaultEnvRetriever returns all system environment variables
type DefaultEnvRetriever struct {
	PropertyMasker *PropertyMasker
}

// GetApplicationEnv loads system environment variables and returns them as a standard ApplicationEnv structure
func (der DefaultEnvRetriever) GetApplicationEnv() *ApplicationEnv {
	properties := make(map[string]PropertyValue)
	for _, pair := range os.Environ() {
		variable := strings.Split(pair, "=")
		key := variable[0]
		value := variable[1]
		properties[key] = der.PropertyMasker.GetPropertyValue(key, value)
	}
	propertySource := PropertySource{
		Name:       "systemEnvironment",
		Properties: properties,
	}
	appenv := ApplicationEnv{
		ActiveProfiles:  []string{},
		PropertySources: []PropertySource{propertySource},
	}

	return &appenv
}

// SetKeysToMask sets a set of property keys that will be masked with ***
func (der DefaultEnvRetriever) SetKeysToMask(keysToMask []string) {
	if der.PropertyMasker == nil {
		der.PropertyMasker = &PropertyMasker{keysToMask: keysToMask}
	} else {
		der.PropertyMasker.keysToMask = keysToMask
	}
}

// PropertyMasker masks values for some keys to protect secrets
type PropertyMasker struct {
	keysToMask []string
}

// GetPropertyValue creates a PropertyValue with masking of values were appropriate
func (pm *PropertyMasker) GetPropertyValue(key string, rawvalue string) PropertyValue {
	propvalue := PropertyValue{}

	if pm.maskForKey(key) {
		propvalue.Value = "***"
	} else {
		propvalue.Value = rawvalue
	}
	return propvalue
}

// SetKeysToMask sets a set of property keys that will be masked with ***
func (pm *PropertyMasker) SetKeysToMask(keysToMask []string) {
	pm.keysToMask = keysToMask
}

func (pm *PropertyMasker) maskForKey(key string) bool {
	for _, sankey := range defaultKeysToSanitize {
		lowkey, lowSanKey := strings.ToLower(key), strings.ToLower(sankey)
		if strings.Contains(lowkey, lowSanKey) {
			return true
		}
	}
	if pm.keysToMask != nil {
		for _, sankey := range pm.keysToMask {
			lowkey, lowSanKey := strings.ToLower(key), strings.ToLower(sankey)
			if strings.Contains(lowkey, lowSanKey) {
				return true
			}
		}
	}

	return false
}
