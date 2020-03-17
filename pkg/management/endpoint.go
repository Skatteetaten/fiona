package management

import (
	"fmt"
	"net/http"
)

// EndPointType is a type for management endpoints
type EndPointType string

// Valid EndPointTypes used for management interface
const (
	Management EndPointType = "management"
	Health     EndPointType = "health"
	Info       EndPointType = "info"
	Env        EndPointType = "env"
)

const (
	defaultManagementPath = "/management"
	defaultHealthPath     = "/health"
	defaultEnvPath        = "/env"
	defaultInfoPath       = "/info"
)

type endpoint struct {
	endpointid  EndPointType
	path        string
	handlerfunc func(http.ResponseWriter, *http.Request)
}

func (ept EndPointType) isValid() error {
	switch ept {
	case Management, Health, Info, Env:
		return nil
	}
	return fmt.Errorf("Invalid endpointid type: %s", ept)
}

func newEndPoint(ept EndPointType, handler func(http.ResponseWriter, *http.Request)) (*endpoint, error) {
	if err := ept.isValid(); err != nil {
		return nil, err
	}

	var ep = &endpoint{
		endpointid:  ept,
		handlerfunc: handler,
	}
	switch ept {
	case Management:
		ep.path = defaultManagementPath
	case Health:
		ep.path = defaultHealthPath
	case Info:
		ep.path = defaultInfoPath
	case Env:
		ep.path = defaultEnvPath
	}
	return ep, nil
}

func (endpoint endpoint) setPath(path string) {
	endpoint.path = path
}

func (endpoint endpoint) getEndpointURL(hostWithPort string) string {
	return fmt.Sprintf("http://%s%s", hostWithPort, endpoint.path)
}
