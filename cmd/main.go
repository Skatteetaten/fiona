package main

import (
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/fiona/pkg/apis"
	"github.com/skatteetaten/fiona/pkg/config"
	"net/http"
)

func main() {

	appConfigReader := config.NewConfigReader()
	appConfig, err := appConfigReader.ReadConfig()
	if err != nil {
		logrus.Fatal("Failed to read application config")
	}

	web := apis.Web{
		Config: appConfig,
	}

	logrus.Info("Starting the webserver")

	web.InitAPI()

	logrus.Fatal(http.ListenAndServe(":8080", nil))

}
