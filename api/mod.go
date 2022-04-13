package api

import (
	"errors"
	"net/http"

	"github.com/desotech-it/whoami/handlers"
	"github.com/desotech-it/whoami/version"
	"github.com/gorilla/mux"
)

var (
	ErrMissingVersion = errors.New("missing version")
)

type Config struct {
}

// GET /version
func getVersionHandler(version *version.Info) http.Handler {
	return handlers.ReadHandler(handlers.JSONSerializerHandler(version))
}

func Handler(version *version.Info, config *Config) http.Handler {
	if version == nil {
		panic(ErrMissingVersion)
	}
	r := mux.NewRouter()
	{
		handler := getVersionHandler(version)
		r.Handle("/version", handler)
		r.Handle("/version/", handler)
	}
	return r
}
