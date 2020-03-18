package healthcheck

import (
	"fmt"
	"github.com/minio/minio/pkg/madmin"
	"github.com/sirupsen/logrus"
	"github.com/skatteetaten/aurora-management-interface-go/health"
	"time"
)

const timeoutSeconds = 2

type serverInfoRetriever interface {
	ServerInfo() (madmin.InfoMessage, error)
}

// FionaHealthRetriever retreives health status information
type FionaHealthRetriever struct {
	serverInfoRetriever
}

// NewFionaHealthRetriever is a factory for HealthHandler
func NewFionaHealthRetriever(admClient serverInfoRetriever) *FionaHealthRetriever {
	return &FionaHealthRetriever{admClient}
}

// GetApplicationHealth returns a health.ApplicationHealth structure for fiona
func (fhh *FionaHealthRetriever) GetApplicationHealth() *health.ApplicationHealth {
	minioHealth := health.ComponentHealth{}
	serverInfoMessage, err := serverInfoWithTimeoutAfter2sec(fhh)
	if err != nil {
		logrus.Errorf("Failed during health check call to minio: %s", err)
		minioHealth.Status = health.Observe
	} else {
		if serverInfoMessage.Mode == "online" {
			minioHealth.Status = health.Up
		} else {
			minioHealth.Status = health.Down
		}
	}

	fionaHealth := health.ApplicationHealth{
		Status:     minioHealth.Status,
		Components: map[string]health.ComponentHealth{"minio": minioHealth},
	}

	return &fionaHealth
}

func serverInfoWithTimeoutAfter2sec(retriever serverInfoRetriever) (madmin.InfoMessage, error) {

	type ChanResult struct {
		InfoMessage madmin.InfoMessage
		Error       error
	}
	serverInfoChan := make(chan ChanResult)
	go func() {
		infoMsg, err := retriever.ServerInfo()
		if err != nil {
			serverInfoChan <- ChanResult{
				InfoMessage: madmin.InfoMessage{},
				Error:       err,
			}
		} else {
			serverInfoChan <- ChanResult{
				InfoMessage: infoMsg,
				Error:       nil,
			}
		}
	}()

	select {
	case res := <-serverInfoChan:
		return res.InfoMessage, res.Error
	case <-time.After(timeoutSeconds * time.Second):
		return madmin.InfoMessage{}, fmt.Errorf("Timed out after %d seconds", timeoutSeconds)
	}
}
