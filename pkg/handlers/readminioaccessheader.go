package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/skatteetaten/fiona/pkg/s3"
	"net/http"
)

// MinioAccess is a key for incoming minio access parameters in header
const MinioAccess = "MINIO_ACCESS"

func readMinioAccessHeader(requestHeader http.Header) (*s3.MinioAccessConfig, error) {
	minioAccessHeader := requestHeader.Get(MinioAccess)
	if len(minioAccessHeader) <= 0 {
		return nil, fmt.Errorf("Found no %s header", MinioAccess)
	}
	b64DecodedStr, err := base64.StdEncoding.DecodeString(minioAccessHeader)
	if err != nil {
		return nil, fmt.Errorf("Could not decode %s: %w", MinioAccess, err)
	}
	var minioAccessConfig s3.MinioAccessConfig
	err = json.Unmarshal([]byte(b64DecodedStr), &minioAccessConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal %s: %w", MinioAccess, err)
	}

	return &minioAccessConfig, nil
}
