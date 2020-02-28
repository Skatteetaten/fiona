package s3

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

// UserManager is a manager for managing users
type UserManager interface {
	CreateUser(userName string, path string) (*CreateUserResult, error)
}

const cannedPolicyTemplate = `
{
  "Version": "2012-10-17",
  "Statement": [
      {
         "Effect":"Allow",
         "Action":[
            "s3:ListAllMyBuckets",
            "s3:GetBucketLocation"
         ],
         "Resource":"arn:aws:s3:::*"
      },
      {
         "Effect":"Allow",
         "Action":[
            "s3:ListBucket",
            "s3:GetBucketLocation"
         ],
         "Resource":"arn:aws:s3:::%s"
      },
      {
         "Effect":"Allow",
         "Action":[
            "s3:PutObject",
            "s3:GetObject",
            "s3:DeleteObject"
         ],
         "Resource":"arn:aws:s3:::%s/%s/*"
      }
  ]
}`

type userClient interface {
	AddUser(accessKey, secretKey string) error
	AddCannedPolicy(policyName, policy string) error
	SetPolicy(policyName, entityName string, isGroup bool) error
}

// MinioUserManager provides methods to manage a users
type MinioUserManager struct {
	userClient
	randomUserpass  bool
	defaultUserpass string
	defaultBucket   string
	serviceEndpoint string
	bucketRegion    string
}

// CreateUserResult provides a map of return values after creating user
type CreateUserResult struct {
	SecretKey       string
	ServiceEndpoint string
	Bucket          string
	BucketRegion    string
}

// NewMinioUserManager is a factory for MinioUserManager
func NewMinioUserManager(s3config *Config) (*MinioUserManager, error) {
	minioAdmClient, err := NewAdmClient(s3config)
	if err != nil {
		return nil, err
	}

	return &MinioUserManager{
		userClient:      minioAdmClient,
		randomUserpass:  s3config.RandomUserpass,
		defaultUserpass: s3config.DefaultUserpass,
		defaultBucket:   s3config.DefaultBucket,
		serviceEndpoint: s3config.getServiceEndpoint(),
		bucketRegion:    s3config.S3Region,
	}, nil
}

// CreateUser creates a user with access policy for a folder path
func (userman *MinioUserManager) CreateUser(userName string, path string) (*CreateUserResult, error) {
	secret := userman.getUserSecret()
	if err := userman.AddUser(userName, secret); err != nil {
		logrus.Error("Could not create new user")
		return nil, err
	}

	if err := userman.createCannedPolicyForUser(userName, path); err != nil {
		logrus.Error("Could not create access policy for user")
		return nil, err
	}
	return &CreateUserResult{
		SecretKey:       secret,
		ServiceEndpoint: userman.serviceEndpoint,
		Bucket:          userman.defaultBucket,
		BucketRegion:    userman.bucketRegion,
	}, nil
}

func (userman *MinioUserManager) getUserSecret() string {
	if userman.randomUserpass {
		return createRandomString()
	}
	return userman.defaultUserpass
}

func (userman *MinioUserManager) createCannedPolicyForUser(username string, path string) error {
	bucket := userman.defaultBucket
	policyName := fmt.Sprintf("RWD%s%s_%d", bucket, path, rand.Intn(1000))
	policy := fmt.Sprintf(cannedPolicyTemplate, bucket, bucket, path)
	if err := userman.AddCannedPolicy(policyName, policy); err != nil {
		logrus.Errorf("Failed to create canned policy %s: %s", policyName, err)
		return err
	}
	if err := userman.SetPolicy(policyName, username, false); err != nil {
		logrus.Errorf("Failed to set policy %s for user %s: %s", policyName, username, err)
		return err
	}

	logrus.Infof("Success: Created policy %s and assigned to user %s.", policyName, username)
	return nil
}

func createRandomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 10 + rand.Intn(6)
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}
