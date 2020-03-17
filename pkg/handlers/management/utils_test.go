package management

import (
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/s3"
)

func getTestAppConfig() *config.Config {
	return &config.Config{
		S3Config: s3.Config{
			S3Host:          "localhost",
			S3Port:          "9000",
			S3UseSSL:        false,
			S3Region:        "us-east-1",
			RandomUserpass:  false,
			DefaultUserpass: "S3userpass",
			AccessKey:       "minio",
			SecretKey:       "minio",
			DefaultBucket:   "utv",
		},
	}
}
