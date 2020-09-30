package s3

import (
	"fmt"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
)

// BucketManager is an interfacce for bucket management
type BucketManager interface {
	MakeSureBucketExists() error
	BucketNameExists(bucketName string) (bool, error)
}

type bucketClient interface {
	BucketExists(bucketName string) (bool, error)
	MakeBucket(bucketName string, location string) (err error)
	SetBucketPolicy(bucketName, policy string) error
}

// MinioBucketManager provides methods to manage a bucket
type MinioBucketManager struct {
	bucketClient
	Config
}

const bucketPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:ListAllMyBuckets",
        "s3:ListBucket",
		"s3:GetBucketLocation"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::%s",
      "Principal": "*"
    }
  ]
}`

// NewMinioBucketManager is a factory for MinioBucketManager
func NewMinioBucketManager(s3config *Config, minioClient *minio.Client) *MinioBucketManager {
	return &MinioBucketManager{minioClient, *s3config}
}

// MakeSureBucketExists checks that the default bucket exists and creates it if not
func (bucketManager *MinioBucketManager) MakeSureBucketExists() error {
	return bucketManager.makeSureNamedBucketExists(bucketManager.DefaultBucket)
}

// BucketNameExists checks that the named bucket exists
func (bucketManager *MinioBucketManager) BucketNameExists(bucketName string) (bool, error) {
	return bucketManager.BucketExists(bucketName)
}

// MakeSureNamedBucketExists checks that the named bucket exists and creates it if not
func (bucketManager *MinioBucketManager) makeSureNamedBucketExists(bucketName string) error {
	found, _ := bucketManager.BucketExists(bucketManager.DefaultBucket)
	if !found {
		err := bucketManager.MakeBucket(bucketManager.DefaultBucket, bucketManager.S3Region)
		if err != nil {
			logrus.Errorf("Could not create missing bucket %s in region %s", bucketManager.DefaultBucket, bucketManager.S3Region)
			return err
		}
		logrus.Infof("Created bucket %s", bucketManager.DefaultBucket)
	} else {
		logrus.Infof("Found existing bucket %s", bucketManager.DefaultBucket)
	}

	if err := bucketManager.setGeneralBucketPolicy(bucketManager.DefaultBucket); err != nil {
		logrus.Errorf("Could not set general bucket policy on bucket %s.", bucketManager.DefaultBucket)
		return err
	}
	return nil
}

func (bucketManager *MinioBucketManager) setGeneralBucketPolicy(bucketName string) error {
	bucketpolicy := fmt.Sprintf(bucketPolicy, bucketName)
	logrus.Debugf("Policy: \n %s", bucketpolicy)
	err := bucketManager.SetBucketPolicy(bucketName, bucketpolicy)
	if err != nil {
		logrus.Errorf("Bucket policy could not be set: %s", err)
		return err
	}
	logrus.Infof("General bucket policy is set on bucket: %s", bucketName)
	return nil
}
