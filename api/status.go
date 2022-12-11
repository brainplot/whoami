package api

import (
	"net/http"
	"net/url"

	"github.com/desotech-it/whoami/status"
	"github.com/gorilla/handlers"
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

func statusHandler(get, put http.Handler) http.Handler {
	handler := handlers.MethodHandler{}
	handler[http.MethodGet] = get
	handler[http.MethodHead] = get
	handler[http.MethodPut] = put
	return handler
}
