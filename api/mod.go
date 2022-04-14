package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/desotech-it/whoami/api/memory"
	"github.com/desotech-it/whoami/handlers"
	"github.com/desotech-it/whoami/version"
	"github.com/gorilla/mux"
)

var (
	ErrMissingVersion = errors.New("missing version")
)

type Config struct {
	VirtualMemoryProvider VirtualMemoryProvider
}

func NewConfig() *Config {
	return &Config{
		VirtualMemoryProvider: VirtualMemoryProviderFunc(memory.VirtualMemoryWithContext),
	}
}

func handleCancelableOperation(w http.ResponseWriter, err error, code int) {
	if err == context.Canceled {
		panic(http.ErrAbortHandler)
	}
	http.Error(w, err.Error(), code)
}

// GET /version
func getVersionHandler(version *version.Info) http.Handler {
	return handlers.ReadHandler(handlers.JSONSerializerHandler(version))
}

// GET /memory
func getMemoryHandler(provider VirtualMemoryProvider) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		if vm, err := provider.VirtualMemory(r.Context()); err != nil {
			handleCancelableOperation(w, err, http.StatusInternalServerError)
		} else {
			handlers.JSONSerializerHandler(vm).ServeHTTP(w, r)
		}
	}
	return handlers.ReadHandler(http.HandlerFunc(middleware))
}

func Handler(version *version.Info, config *Config) http.Handler {
	if version == nil {
		panic(ErrMissingVersion)
	}
	r := mux.NewRouter()
	r.Handle("/version", getVersionHandler(version))
	r.Handle("/memory", getMemoryHandler(config.VirtualMemoryProvider))
	return r
}
