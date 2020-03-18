// healthcheck_test.go
package healthcheck

import (
	"github.com/minio/minio/pkg/madmin"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/skatteetaten/fiona/pkg/s3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Run("Should timeout when client has no connected s3 instance", func(t *testing.T) {
		hook := test.NewGlobal()
		s3AdmClient, _ := s3.NewAdmClient(&getTestAppConfig().S3Config)

		fhr := NewFionaHealthRetriever(s3AdmClient)
		healthResult := fhr.GetApplicationHealth()

		assert.NotNil(t, healthResult, "Health check should return something")
		assert.True(t, len(hook.Entries) > 0, "Should have at least one log entry")
		lastLogIdx := len(hook.Entries) - 1
		assert.Equal(t, logrus.ErrorLevel, hook.Entries[lastLogIdx].Level)
		assert.Contains(t, hook.Entries[lastLogIdx].Message, "Timed out after 2 seconds")

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})
	t.Run("Should return OK", func(t *testing.T) {
		fhr := NewFionaHealthRetriever(mockServerInfoRetriever{})

		healthResult := fhr.GetApplicationHealth()

		assert.NotNil(t, healthResult, "Health check should return something")
		assert.Equal(t, "UP", healthResult.Status)
		assert.Equal(t, 1, len(healthResult.Components))
		assert.Equal(t, "UP", healthResult.Components["minio"].Status)
	})
}

type mockServerInfoRetriever struct{}

func (msir mockServerInfoRetriever) ServerInfo() (madmin.InfoMessage, error) {
	return madmin.InfoMessage{
		Mode: "online",
	}, nil
}
