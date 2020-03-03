package s3

import (
	"github.com/minio/minio/pkg/madmin"
	"github.com/sirupsen/logrus"
)

// NewAdmClient creates an S3 admin client
func NewAdmClient(s3config *Config) (*madmin.AdminClient, error) {

	endpoint := endpoint(s3config)
	minioadmin, err := madmin.New(endpoint, s3config.AccessKey, s3config.SecretKey, s3config.S3UseSSL)
	logrus.Infof("Creating minio adminclient to %s using ssl=%t", endpoint, s3config.S3UseSSL)
	if err != nil {
		logrus.Errorf("Could not create S3 admin client %v", err)
		return nil, err
	}

	return minioadmin, nil
}
