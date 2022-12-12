package api

import (
	"net/url"

	"github.com/desotech-it/whoami/status"
)

const (
	HTTPParamStatus = "status"
)

type StatusInfo struct {
	Status string `json:"status"`
}

func ParseStatusInValues(values url.Values) (status.Status, error) {
	statusText := values.Get(HTTPParamStatus)
	if status, err := status.Parse(statusText); err != nil {
		return "", err
	} else {
		return status, nil
	}
}
