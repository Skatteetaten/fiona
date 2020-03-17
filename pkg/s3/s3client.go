package s3

import (
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
)

// NewClient creates the S3 client
func NewClient(s3config *Config) (*minio.Client, error) {

	endpoint := endpoint(s3config)
	minioClient, err := minio.New(endpoint, s3config.AccessKey, s3config.SecretKey, s3config.S3UseSSL)
	logrus.Infof("Creating minio client to %s using ssl=%t", endpoint, s3config.S3UseSSL)
	if err != nil {
		logrus.Errorf("Could not create S3 client %v", err)
		return nil, err
	}

	return minioClient, nil
}
