package health

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Standard health status responses
const (
	Up           = "UP"
	Down         = "DOWN"
	OutOfService = "OUT_OF_SERVICE"
	Unknown      = "UNKNOWN"
	Observe      = "OBSERVE"
)

// HTTPStatus holds default mapping between health status and http response code
var HTTPStatus = map[string]int{
	Up:           http.StatusOK,
	Down:         http.StatusServiceUnavailable,
	OutOfService: http.StatusServiceUnavailable,
	Unknown:      http.StatusOK,
	Observe:      http.StatusOK,
}

// ApplicationHealthRetriever is an interface for health check methods at application level
type ApplicationHealthRetriever interface {
	GetApplicationHealth() *ApplicationHealth
}

// ApplicationHealth is a structure for standardized health check response from applications
type ApplicationHealth struct {
	Status     string                     `json:"status"`
	Components map[string]ComponentHealth `json:"components,omitempty"`
}

// ComponentHealth is a structure for standardized health check response from components
type ComponentHealth struct {
	Status     string                     `json:"status"`
	Components map[string]ComponentHealth `json:"components,omitempty"`
	Details    map[string]interface{}     `json:"details,omitempty"`
}

// ApplicationHealthHandler fetches an ApplicationHealth structure from the application and parse it for a proper http response
type ApplicationHealthHandler struct {
	ApplicationHealthRetriever ApplicationHealthRetriever
}

func (ahh *ApplicationHealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	healthResponse := ahh.ApplicationHealthRetriever.GetApplicationHealth()
	if healthResponse == nil {
		errorResponse(w, "Error getting health response", fmt.Errorf("healthResponse was nil"))
		return
	}

	httpStatus := getHTTPStatus(healthResponse.Status)
	healthResponseJSON, err := json.Marshal(healthResponse)
	if err != nil {
		errorResponse(w, "Error while creating JSON for health response", err)
		return
	}
	logrus.Debugf("Health response output:\n%s\n", string(healthResponseJSON))

	w.Header().Set("Content-Type", "application/vnd.spring-boot.actuator.v3+json; charset=UTF-8")
	w.WriteHeader(httpStatus)
	_, _ = fmt.Fprintf(w, "%s", healthResponseJSON)
}

func errorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	logrus.Errorf("%s: %s", message, err)
	w.WriteHeader(http.StatusInternalServerError)
	responseJSON, _ := json.Marshal(map[string]string{
		"error": message,
		"cause": fmt.Sprintf("%v", err),
	})
	_, _ = fmt.Fprintf(w, "%s", responseJSON)
}

// GetDefaultHealthRetriever gets a DefaultHealthRetriever with default initialization
func GetDefaultHealthRetriever() DefaultHealthRetriever {
	return DefaultHealthRetriever{}
}

// DefaultHealthRetriever is a very shallow default health check handler, simply replying status UP
type DefaultHealthRetriever struct{}

// GetApplicationHealth returns Status UP for the DefaultHealthRetriever
func (dhh DefaultHealthRetriever) GetApplicationHealth() *ApplicationHealth {
	var defaultHealthResponse = ApplicationHealth{Status: Up}
	return &defaultHealthResponse
}

func getHTTPStatus(healthStatus string) int {
	status := HTTPStatus[healthStatus]
	if status == 0 {
		logrus.Errorf("Error: Found no HTTPStatus for health status: %s. Returning 500", healthStatus)
		return http.StatusInternalServerError
	}
	return HTTPStatus[healthStatus]
}
