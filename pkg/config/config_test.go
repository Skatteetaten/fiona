package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("Should return config with default values", func(t *testing.T) {
		if getEnvBoolOrDefault("FIONA_DEBUG", false) {
			t.Skip("Skipping default config test when FIONA_DEBUG is true")
		}
		confreader := ConfReader{}

		config, err := confreader.ReadConfig()
		assert.Nil(t, err)

		assert.Equal(t, "localhost", config.S3Config.S3Host)
		assert.Equal(t, "9000", config.S3Config.S3Port)
		assert.Equal(t, false, config.S3Config.S3UseSSL)
		assert.Equal(t, "us-east-1", config.S3Config.S3Region)
		assert.Equal(t, true, config.S3Config.RandomUserpass)
		assert.Equal(t, "S3userpass", config.S3Config.DefaultUserpass)
		assert.Equal(t, "aurora", config.S3Config.AccessKey)
		assert.Equal(t, "fragleberget", config.S3Config.SecretKey)
		assert.Equal(t, "utv", config.S3Config.DefaultBucket)
		assert.Equal(t, false, config.DebugLog)
	})

}
