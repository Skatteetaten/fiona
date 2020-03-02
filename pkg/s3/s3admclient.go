package s3

import (
	"github.com/minio/minio/pkg/madmin"
	"github.com/sirupsen/logrus"
)

// NewAdmClient creates an S3 admin client
func NewAdmClient(s3config *Config) (*madmin.AdminClient, error) {

	minioadmin, err := madmin.New(endpoint(s3config), s3config.AccessKey, s3config.SecretKey, s3config.S3UseSSL)
	minioadmin.SetCustomTransport(createTransportWithInsecureTls())
	if err != nil {
		logrus.Errorf("Could not create S3 admin client %v", err)
		return nil, err
	}

	return minioadmin, nil
}
