package main

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/apis"
	"github.com/skatteetaten/fiona/pkg/config"
	"net/http"
)

func main() {
	// Remove this when certificates on minio servers are fixed
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	appConfigReader := config.NewConfigReader()
	initWebServer(appConfigReader)

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func initWebServer(appConfigReader config.Reader) {
	appConfig, err := appConfigReader.ReadConfig()
	if err != nil {
		logrus.Fatal("Failed to read application config")
	}
	if appConfig.DebugLog {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Info("Starting the webserver")
	err = apis.InitAPI(appConfig)
	if err != nil {
		logrus.Fatal(err)
	}
}
