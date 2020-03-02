package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func failLogAndResponse(w http.ResponseWriter, message string, status int, err error) {
	logrus.Errorf("%s: %s", message, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	responseJSON, err := json.Marshal(map[string]string{
		"error": message,
		"cause": fmt.Sprintf("%v", err),
	})
	_, _ = fmt.Fprintf(w, "%s", responseJSON)
}
