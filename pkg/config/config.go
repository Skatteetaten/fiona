package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/s3"
	"os"
	"strconv"
)

const auroraTokenLocation = "./aurora-token"

// Environment variables for external reference
const (
	FionaDefaultPassword = "FIONA_DEFAULT_PASSWORD"
	FionaSecretKey       = "FIONA_SECRET_KEY"
	FionaAccessKey       = "FIONA_ACCESS_KEY"
)

// Config for the S3 access
type Config struct {
	S3Config            s3.Config
	DebugLog            bool // default "false"
	AuroraTokenLocation string
}

// Reader interface
type Reader interface {
	ReadConfig() (*Config, error)
}

// ConfReader is a Reader receiver
type ConfReader struct {
}

// NewConfigReader factory method
func NewConfigReader() Reader {
	return &ConfReader{}
}

// ReadConfig implementation
func (m *ConfReader) ReadConfig() (*Config, error) {
	useSsl := getEnvBoolOrDefault("FIONA_S3_USESSL", false)
	randomUserpass := getEnvBoolOrDefault("FIONA_RANDOMPASS", true)
	debuglog := getEnvBoolOrDefault("FIONA_DEBUG", false)

	return &Config{
		S3Config: s3.Config{
			S3Host:          getEnvOrDefault("FIONA_S3_HOST", "localhost"),
			S3Port:          getEnvOrDefault("FIONA_S3_PORT", "9000"),
			S3UseSSL:        useSsl,
			S3Region:        getEnvOrDefault("FIONA_S3_REGION", "us-east-1"),
			RandomUserpass:  randomUserpass,
			DefaultUserpass: getEnvOrDefault(FionaDefaultPassword, "S3userpass"),
			AccessKey:       getEnvOrDefault(FionaAccessKey, "aurora"),
			SecretKey:       getEnvOrDefault(FionaSecretKey, "fragleberget"),
			DefaultBucket:   getEnvOrDefault("FIONA_DEFAULTBUCKET", "utv"),
		},
		DebugLog:            debuglog,
		AuroraTokenLocation: getEnvOrDefault("FIONA_AURORATOKENLOCATION", auroraTokenLocation),
	}, nil
}

func getEnvBoolOrDefault(key string, fallback bool) bool {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	valueBool, err := strconv.ParseBool(value)
	if err != nil {
		logrus.Warnf(fmt.Sprintf("%s must be a boolean (true or false), was %s. Using fallback value.", key, value))
		return fallback
	}
	return valueBool
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
