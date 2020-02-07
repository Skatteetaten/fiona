package config

import (
	"os"
)

// Config for the S3 access
type Config struct {
	S3Host string
	S3Port string
}

// Reader interface
type Reader interface {
	ReadConfig() (*Config, error)
}

// S3ConfigReader is an Reader receiver
type S3ConfigReader struct {
}

// NewConfigReader factory method
func NewConfigReader() Reader {
	return &S3ConfigReader{}
}

// ReadConfig implementation
func (m *S3ConfigReader) ReadConfig() (*Config, error) {
	return &Config{
		S3Host: getEnvOrDefault("FIONA_S3_HOST", "minio"),
		S3Port: getEnvOrDefault("FIONA_S3_PORT", "6377"),
	}, nil
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
