package api

import (
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
	VirtualMemoryProvider memory.VirtualMemoryProvider
}

func NewConfig() *Config {
	return &Config{
		VirtualMemoryProvider: memory.VirtualMemoryProviderFunc(memory.VirtualMemoryWithContext),
	}
}

// GET /version
func getVersionHandler(version *version.Info) http.Handler {
	return handlers.ReadHandler(handlers.JSONSerializerHandler(version))
}

// GET /memory
func getMemoryHandler(provider memory.VirtualMemoryProvider) http.Handler {
	middleware := func(w http.ResponseWriter, r *http.Request) {
		if vm, err := provider.VirtualMemory(r.Context()); err != nil {
			errorHandler := handlers.ErrorHandler(err.Error(), http.StatusInternalServerError)
			handlers.CancellableHandler(err, errorHandler).ServeHTTP(w, r)
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
