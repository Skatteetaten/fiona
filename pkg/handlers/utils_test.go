package handlers

import (
	"encoding/json"
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/s3"
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func isJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil
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
	}
}
