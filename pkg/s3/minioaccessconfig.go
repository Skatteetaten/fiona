package s3

// MinioAccessConfig for external specification of access parameters
type MinioAccessConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	UseSsl    bool   `json:"useSsl"` // Default true
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}
