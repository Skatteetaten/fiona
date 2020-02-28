package main

import (
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/s3"
	"testing"
)

type testConfigReader struct {
}

func (tcr testConfigReader) ReadConfig() (*config.Config, error) {
	return getTestAppConfig(), nil
}

func TestMainProgram(t *testing.T) {
	t.Run("Should run initWebServer() without failing", func(t *testing.T) {
		initWebServer(testConfigReader{})
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
