package main

import (
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/apis"
	"github.com/skatteetaten/fiona/pkg/config"
	"github.com/skatteetaten/fiona/pkg/s3"
	"net/http"
)

func main() {
	appConfigReader := config.NewConfigReader()

	initWebServer(appConfigReader)

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}

func initWebServer(appConfigReader config.Reader) {
	appConfig, err := appConfigReader.ReadConfig()
	if err != nil {
		logrus.Fatalf("Fatal error: Failed to read application config: %s", err)
	}
	if appConfig.DebugLog {
		logrus.SetLevel(logrus.DebugLevel)
	}

	minioAdmClient, err := s3.NewAdmClient(&appConfig.S3Config)
	if err != nil {
		logrus.Fatalf("Fatal error: Failed to create s3 adminclient: %s", err)
	}
	minioclient, err := s3.NewClient(&appConfig.S3Config)
	if err != nil {
		logrus.Fatalf("Fatal error: Failed to create s3 client: %s", err)
	}

	logrus.Info("Starting the webserver")
	err = apis.InitAPI(appConfig, minioAdmClient, minioclient)
	if err != nil {
		logrus.Fatal(err)
	}

	managementInterfaceHandler := apis.InitManagementHandler(minioAdmClient)
	managementInterfaceHandler.StartHTTPListener()
}
