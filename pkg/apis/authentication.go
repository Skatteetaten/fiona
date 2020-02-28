package apis

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const tokenPrefix = "aurora-token"

// AuthMiddleware is an interface for authentication
type AuthMiddleware interface {
	Authenticate(next http.Handler) http.Handler
}

// AuroraTokenAuthenticator handles authentication for certain routes in api
type AuroraTokenAuthenticator struct {
	auroratoken string
}

// NewAuroraTokenAuthenticator creates and initializes an AuroraTokenAuthenticator
func NewAuroraTokenAuthenticator(auroraTokenLocation string) (*AuroraTokenAuthenticator, error) {
	auroratoken, err := getAuroraToken(auroraTokenLocation)
	if err != nil {
		logrus.Panic("Could not initialize auroratoken. Aborting fiona")
	}

	return &AuroraTokenAuthenticator{auroratoken: auroratoken}, nil

}

// Authenticate verifies that request token is valid
func (amw *AuroraTokenAuthenticator) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		found := amw.equalToAuroraToken(token)
		if found {
			logrus.Info("Authentication OK")
			next.ServeHTTP(w, r)
		} else {
			logrus.Warn("Authentication failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func (amw *AuroraTokenAuthenticator) equalToAuroraToken(token string) bool {
	trimmedToken := strings.TrimSpace(token)

	if strings.Index(trimmedToken, tokenPrefix+" ") == 0 {
		trimmedToken = strings.TrimPrefix(trimmedToken, tokenPrefix+" ")
	}

	found := trimmedToken == amw.auroratoken

	return found
}

func getAuroraToken(auroraTokenLocation string) (string, error) {
	if fileExists(auroraTokenLocation) {
		auroratoken, err := ioutil.ReadFile(auroraTokenLocation)
		if err != nil {
			logrus.Errorf("Could not read file at %s : %s", auroraTokenLocation, err)
			return "", err

		}
		if string(auroratoken) == "" {
			msg := fmt.Sprintf("Found empty auroratoken file on %s. Authorization token required.", auroraTokenLocation)
			logrus.Error("msg")
			return "", errors.New(msg)
		}
		return string(auroratoken), nil
	}
	return "", errors.New("No auroratoken file found")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
