package s3

import "fmt"

func endpoint(s3config *Config) string {
	return s3config.S3Host + ":" + s3config.S3Port
}

func (c *Config) getServiceEndpoint() string {

	protocol := "http"
	if c.S3UseSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s", protocol, endpoint(c))
}
