package s3

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minio/minio/pkg/madmin"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

// UserManager is a manager for managing users
type UserManager interface {
	CreateUser(userName string, path string) (*CreateUserResult, error)
	CreateAppUser(createAppUserInput *CreateAppUserInput) (*CreateAppUserResult, error)
}

const cannedPolicyTemplateForOldUser = `
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

// CreateAppUserInput provides input for creating an application user
type CreateAppUserInput struct {
	Bucketname string   `json:"bucketname"`
	Path       string   `json:"path"`
	Username   string   `json:"username"`
	Access     []string `json:"access"`
}

// CreateAppUserResult provides information after for creating an application user
type CreateAppUserResult struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	HostURL   string `json:"host"`
}

// NewMinioUserManager is a factory for MinioUserManager
func NewMinioUserManager(s3config *Config, adminClient *madmin.AdminClient) *MinioUserManager {
	return &MinioUserManager{
		userClient:      adminClient,
		randomUserpass:  s3config.RandomUserpass,
		defaultUserpass: s3config.DefaultUserpass,
		defaultBucket:   s3config.DefaultBucket,
		serviceEndpoint: s3config.getServiceEndpoint(),
		bucketRegion:    s3config.S3Region,
	}
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

// CreateAppUser creates a user with access policy for a folder path
func (userman *MinioUserManager) CreateAppUser(createAppUserInput *CreateAppUserInput) (*CreateAppUserResult, error) {
	secret := userman.getUserSecret()
	if err := userman.AddUser(createAppUserInput.Username, secret); err != nil {
		logrus.Errorf("Could not create new user: %s", createAppUserInput.Username)
		return nil, err
	}

	if err := userman.createCannedPolicyForAppUser(createAppUserInput); err != nil {
		logrus.Error("Could not create access policy for user")
		return nil, err
	}
	return &CreateAppUserResult{
		AccessKey: createAppUserInput.Username,
		SecretKey: secret,
		HostURL:   userman.serviceEndpoint,
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
	policy := fmt.Sprintf(cannedPolicyTemplateForOldUser, bucket, bucket, path)
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

func (userman *MinioUserManager) createCannedPolicyForAppUser(createAppUserInput *CreateAppUserInput) error {
	bucket := createAppUserInput.Bucketname
	path := createAppUserInput.Path
	username := createAppUserInput.Username
	policyNamePostfix := ""
	for _, s := range createAppUserInput.Access {
		policyNamePostfix += s[0:1]
	}
	policyName := fmt.Sprintf("%s%s_%s_%s", bucket, path, username, policyNamePostfix)

	generatedAppUserPolicy, err := generateAppUserPolicy(createAppUserInput)
	if err != nil {
		return err
	}

	policy, err := json.Marshal(generatedAppUserPolicy)
	if err != nil {
		logrus.Errorf("Failed to create canned policy %s: %s", policyName, err)
		return err
	}

	if err := userman.AddCannedPolicy(policyName, string(policy)); err != nil {
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

func generateAppUserPolicy(createAppUserInput *CreateAppUserInput) (map[string]interface{}, error) {
	bucket := createAppUserInput.Bucketname
	path := createAppUserInput.Path

	s3ObjectActions, err := getS3ObjectActions(createAppUserInput.Access)
	if err != nil {
		return nil, err
	}
	generatedAppUserPolicy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"s3:ListAllMyBuckets",
					"s3:GetBucketLocation",
				},
				"Resource": "arn:aws:s3:::*",
			},
			{
				"Effect": "Allow",
				"Action": []string{
					"s3:ListBucket",
					"s3:GetBucketLocation",
				},
				"Resource": fmt.Sprintf("arn:aws:s3:::%s/*", bucket),
			},
			{
				"Effect":   "Allow",
				"Action":   s3ObjectActions,
				"Resource": fmt.Sprintf("arn:aws:s3:::%s/%s/*", bucket, path),
			},
		},
	}
	return generatedAppUserPolicy, nil
}

func getS3ObjectActions(access []string) ([]string, error) {
	if len(access) <= 0 {
		return []string{
			"s3:PutObject",
			"s3:GetObject",
			"s3:DeleteObject",
		}, nil
	}
	var actions []string
	for _, s := range access {
		switch strings.ToUpper(s) {
		case "READ":
			actions = append(actions, "s3:GetObject")
		case "WRITE":
			actions = append(actions, "s3:PutObject")
		case "DELETE":
			actions = append(actions, "s3:DeleteObject")
		default:
			return nil, errors.New("Got illegal access parameter")
		}

	}
	return actions, nil
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
