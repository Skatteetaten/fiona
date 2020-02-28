package main

import (
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/apis"
	"github.com/skatteetaten/fiona/pkg/config"
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
