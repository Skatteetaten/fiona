package s3

import (
	"crypto/tls"
	"net/http"
)

func createTransportWithInsecureTls() http.RoundTripper {
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	var transport http.RoundTripper = &http.Transport{
		TLSClientConfig:    tlsConfig,
		DisableCompression: true,
	}
	return transport
}
