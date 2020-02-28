package handlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func failLogAndResponse(w http.ResponseWriter, message string, status int, err error) {
	logrus.Errorf("%s: %s", message, err)
	w.Header().Set("Content-Type", "application/text; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, "%s.", message)
}
