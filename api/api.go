package api

import (
	"net/http"

	"github.com/desotech-it/whoami/api/memory"
	"github.com/desotech-it/whoami/api/net"
	"github.com/desotech-it/whoami/handlers"
	"github.com/desotech-it/whoami/version"
	"github.com/gorilla/mux"
)

type Api interface {
	// GET /version
	GetVersion(http.ResponseWriter, *http.Request)
	// GET /memory
	GetMemory(http.ResponseWriter, *http.Request)
	// GET /interfaces
	GetInterfaces(http.ResponseWriter, *http.Request)
}

type Server struct {
	Version               *version.Info
	VirtualMemoryProvider memory.VirtualMemoryProvider
	InterfacesProvider    net.InterfacesProvider
}

func NewServer(versionInfo *version.Info) *Server {
	if versionInfo == nil {
		panic("versionInfo must not be nil")
	}
	return &Server{
		Version:               versionInfo,
		VirtualMemoryProvider: memory.VirtualMemoryProviderFunc(memory.VirtualMemoryWithContext),
		InterfacesProvider:    net.InterfacesProviderFunc(net.Interfaces),
	}
}

func (s *Server) GetVersion(w http.ResponseWriter, r *http.Request) {
	handlers.ReadHandler(handlers.JSONSerializerHandler(s.Version)).ServeHTTP(w, r)
}

func (s *Server) GetMemory(w http.ResponseWriter, r *http.Request) {
	handlers.ReadHandler(http.HandlerFunc(s.serializeVirtualMemory)).ServeHTTP(w, r)
}

func (s *Server) serializeVirtualMemory(w http.ResponseWriter, r *http.Request) {
	if vm, err := s.VirtualMemoryProvider.VirtualMemory(r.Context()); err != nil {
		errorHandler := handlers.ErrorHandler(err.Error(), http.StatusInternalServerError)
		handlers.CancellableHandler(err, errorHandler).ServeHTTP(w, r)
	} else {
		handlers.JSONSerializerHandler(vm).ServeHTTP(w, r)
	}
}

func (s *Server) GetInterfaces(w http.ResponseWriter, r *http.Request) {
	handlers.ReadHandler(http.HandlerFunc(s.serializeInterfaces)).ServeHTTP(w, r)
}

func (s *Server) serializeInterfaces(w http.ResponseWriter, r *http.Request) {
	if netInterfaces, err := s.InterfacesProvider.Interfaces(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		handlers.JSONSerializerHandler(netInterfaces).ServeHTTP(w, r)
	}
}

func Handler(api Api) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/version", api.GetVersion)
	r.HandleFunc("/memory", api.GetMemory)
	r.HandleFunc("/interfaces", api.GetInterfaces)
	return r
}
