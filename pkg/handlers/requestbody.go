package handlers

import (
	"io"
	"io/ioutil"
)

func readRequestBody(bodyReadCloser io.ReadCloser) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(bodyReadCloser, 1048576))
	if err != nil {
		return nil, err
	}
	if err := bodyReadCloser.Close(); err != nil {
		return nil, err
	}
	return body, nil
}
