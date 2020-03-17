package s3

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestS3endpoint(t *testing.T) {
	t.Run("Should create endpoint", func(t *testing.T) {

		conf := getTestAppConfig()
		endpoint := endpoint(conf)
		serviceEndpoint := conf.getServiceEndpoint()

		assert.Equal(t, "minio:9000", endpoint)
		assert.Equal(t, "http://minio:9000", serviceEndpoint)
	})
}

func TestS3client(t *testing.T) {
	t.Run("Should create client", func(t *testing.T) {

		newclient, err := NewClient(getTestAppConfig())

		assert.Nil(t, err)
		assert.Contains(t, newclient.EndpointURL().String(), "minio:9000")
	})
}

func TestS3admclient(t *testing.T) {
	t.Run("Should create minio admclient", func(t *testing.T) {

		newadmclient, err := NewAdmClient(getTestAppConfig())
		assert.Nil(t, err)
		assert.NotNil(t, newadmclient)
		// Hard to assert anything in a unit test, since most methods on the admclient require a working backend
	})
}

type testBucketClient struct {
}

func (tbc testBucketClient) BucketExists(bucketName string) (bool, error) {
	return bucketName == "utv", nil
}
func (tbc testBucketClient) MakeBucket(bucketName string, location string) (err error) {
	return nil
}
func (tbc testBucketClient) SetBucketPolicy(bucketName, policy string) error {
	return nil
}

func TestS3bucketmanager(t *testing.T) {
	t.Run("Should create new bucketmanager", func(t *testing.T) {
		dummyClient, _ := NewClient(getTestAppConfig())
		newbucketmanager := NewMinioBucketManager(getTestAppConfig(), dummyClient)
		assert.NotNil(t, newbucketmanager)
	})

	t.Run("Should detect existing bucket", func(t *testing.T) {
		bucketclient := testBucketClient{}
		bucketmanager := MinioBucketManager{
			bucketclient,
			*getTestAppConfig(),
		}

		hook := test.NewGlobal()

		bucketmanager.MakeSureBucketExists()

		assert.Equal(t, 2, len(hook.Entries))
		assert.Equal(t, logrus.InfoLevel, hook.Entries[0].Level)
		assert.Contains(t, hook.Entries[0].Message, "Found existing bucket utv")
		assert.Contains(t, hook.Entries[1].Message, "General bucket policy is set on bucket")

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("Should create bucket if not found", func(t *testing.T) {
		bucketclient := testBucketClient{}
		bucketmanager := MinioBucketManager{
			bucketclient,
			*getTestAppConfigNewbucket(),
		}
		hook := test.NewGlobal()

		bucketmanager.MakeSureBucketExists()

		assert.Equal(t, 2, len(hook.Entries))
		assert.Equal(t, logrus.InfoLevel, hook.Entries[0].Level)
		assert.Contains(t, hook.Entries[0].Message, "Created bucket newbucket")
		assert.Contains(t, hook.Entries[1].Message, "General bucket policy is set on bucket")

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	})

	t.Run("Should successfully create random strings", func(t *testing.T) {
		randomString := createRandomString()
		assert.True(t, len(randomString) > 9)
		assert.True(t, len(randomString) < 17)

		// Testing for randomness is an exercise in probability.
		// This test should never fail, but if the universe choose an instance of improbable non-randomness...
		generatedDifferent := false

		i := 1
		for i <= 100 {
			newRandomString := createRandomString()
			if newRandomString != randomString {
				generatedDifferent = true
				break
			}
			i++
		}

		assert.True(t, generatedDifferent)
		assert.True(t, i <= 100)
		t.Logf("Generated different random string in %d attempts", i)
	})
}

func getTestAppConfig() *Config {
	return &Config{
		S3Host:          "minio",
		S3Port:          "9000",
		S3UseSSL:        false,
		S3Region:        "us-east-1",
		RandomUserpass:  false,
		DefaultUserpass: "S3userpass",
		AccessKey:       "minio",
		SecretKey:       "minio",
		DefaultBucket:   "utv",
	}
}

func getTestAppConfigNewbucket() *Config {
	return &Config{
		S3Host:          "minio",
		S3Port:          "9000",
		S3UseSSL:        false,
		S3Region:        "us-east-1",
		RandomUserpass:  false,
		DefaultUserpass: "S3userpass",
		AccessKey:       "minio",
		SecretKey:       "minio",
		DefaultBucket:   "newbucket",
	}
}
