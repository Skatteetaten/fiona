package apis

import (
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/s3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testAmw struct {
}

func (ta testAmw) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TestApis(t *testing.T) {
	t.Run("Should initialize web router without failing", func(t *testing.T) {
		InitAPI(getTestAppConfig())
	})

	t.Run("Should return correct welcome string on root request to router", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
		request.Header.Set("Authorization", "aurora-token ")
		response := httptest.NewRecorder()

		routerHandler, _ := createRouter(getTestAppConfig(), &testAmw{})
		routerHandler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code, "OK response is expected")
		assert.Equal(t, "Fiona says hi at localhost:8080!", response.Body.String())
	})
}

func getTestAppConfig() *config.Config {
	return &config.Config{
		S3Config: s3.Config{
			S3Host:          "minio",
			S3Port:          "9000",
			S3UseSSL:        false,
			S3Region:        "us-east-1",
			RandomUserpass:  false,
			DefaultUserpass: "S3userpass",
			AccessKey:       "minio",
			SecretKey:       "minio",
			DefaultBucket:   "utv",
		},
		AuroraTokenLocation: "testdata/token",
	}
}
