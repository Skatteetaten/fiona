package s3

// Config for the minio S3 clients and S3 operations
type Config struct {
	S3Host          string
	S3Port          string
	S3UseSSL        bool // Default true
	S3Region        string
	RandomUserpass  bool   // Default true
	DefaultUserpass string // Only used when RandomUserpass is false, default "S3userpass"
	AccessKey       string
	SecretKey       string
	DefaultBucket   string // Default "utv"
}
