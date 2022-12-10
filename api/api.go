package api

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/desotech-it/whoami/api/memory"
	"github.com/desotech-it/whoami/api/net"
	"github.com/desotech-it/whoami/api/os"
	"github.com/desotech-it/whoami/client"
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
	// GET /hostname
	GetHostname(http.ResponseWriter, *http.Request)
	// GET,POST,PUT,DELETE,CONNECT,OPTIONS,TRACE,PATCH /request
	EchoRequest(http.ResponseWriter, *http.Request)
	// GET /memory/stresssession
	GetMemoryStress(http.ResponseWriter, *http.Request)
	// POST /memory/stresssession
	PostMemoryStress(http.ResponseWriter, *http.Request, memory.StressParameters)
	// DELETE /memory/stresssession
	CancelMemoryStress(http.ResponseWriter, *http.Request)
}

type Server struct {
	Version                  *version.Info
	VirtualMemoryProvider    memory.VirtualMemoryProvider
	InterfacesProvider       net.InterfacesProvider
	HostnameProvider         os.HostnameProvider
	memoryStressSession      *memoryStressSession
	memoryStressSessionMutex sync.RWMutex
}

func NewServer(versionInfo *version.Info) *Server {
	if versionInfo == nil {
		panic("versionInfo must not be nil")
	}
	return &Server{
		Version:               versionInfo,
		VirtualMemoryProvider: memory.VirtualMemoryProviderFunc(memory.VirtualMemoryWithContext),
		InterfacesProvider:    net.InterfacesProviderFunc(net.Interfaces),
		HostnameProvider:      os.HostnameProviderFunc(os.Hostname),
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

func (s *Server) GetHostname(w http.ResponseWriter, r *http.Request) {
	handlers.ReadHandler(http.HandlerFunc(s.serializeHostname)).ServeHTTP(w, r)
}

func (s *Server) serializeHostname(w http.ResponseWriter, r *http.Request) {
	if hostname, err := s.HostnameProvider.Hostname(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		handlers.JSONSerializerHandler(hostname).ServeHTTP(w, r)
	}
}

func (s *Server) EchoRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clientRequest := client.RequestFromStdRequest(r)
	handlers.JSONSerializerHandler(clientRequest).ServeHTTP(w, r)
}

func (s *Server) GetMemoryStress(w http.ResponseWriter, r *http.Request) {
	s.memoryStressSessionMutex.RLock()
	if s.memoryStressSession == nil {
		s.memoryStressSessionMutex.RUnlock()
		http.Error(w, "memory stress is not currently running", http.StatusBadRequest)
	} else {
		response := buildMemoryStressSessionResponse(s.memoryStressSession)
		s.memoryStressSessionMutex.RUnlock()
		handlers.JSONSerializerHandler(&response).ServeHTTP(w, r)
	}
}

func (s *Server) PostMemoryStress(w http.ResponseWriter, r *http.Request, params memory.StressParameters) {
	s.memoryStressSessionMutex.Lock()
	if s.memoryStressSession == nil {
		if stresser, err := memory.NewStresser(params); err != nil {
			s.memoryStressSessionMutex.Unlock()
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			ctx, cancelFunc := context.WithCancel(context.Background())
			s.memoryStressSession = &memoryStressSession{
				stressSession: stressSession{
					StartedAt:  time.Now().UTC(),
					cancelFunc: cancelFunc,
				},
				stresser: stresser,
			}
			s.memoryStressSessionMutex.Unlock()
			go stresser.Stress(ctx)
			response := buildMemoryStressSessionResponse(s.memoryStressSession)
			handlers.JSONSerializerHandler(response).ServeHTTP(w, r)
		}
	} else {
		s.memoryStressSessionMutex.Unlock()
		http.Error(w, "memory stress is already running", http.StatusBadRequest)
	}
}

func (s *Server) CancelMemoryStress(w http.ResponseWriter, r *http.Request) {
	s.memoryStressSessionMutex.Lock()
	if s.memoryStressSession == nil {
		s.memoryStressSessionMutex.Unlock()
		http.Error(w, "memory stress is not currently running", http.StatusBadRequest)
	} else {
		cancelFunc := s.memoryStressSession.cancelFunc
		finishedAt := time.Now().UTC()
		s.memoryStressSession.FinishedAt = &finishedAt
		cancelFunc()
		response := buildMemoryStressSessionResponse(s.memoryStressSession)
		s.memoryStressSession = nil
		s.memoryStressSessionMutex.Unlock()
		handlers.JSONSerializerHandler(response).ServeHTTP(w, r)
	}
}

func Handler(api Api) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/version", api.GetVersion)
	r.HandleFunc("/memory", api.GetMemory)
	r.HandleFunc("/interfaces", api.GetInterfaces)
	r.HandleFunc("/hostname", api.GetHostname)
	r.HandleFunc("/request", api.EchoRequest)
	r.Handle("/memory/stresssession", stressHandler(
		http.HandlerFunc(api.GetMemoryStress),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if params, err := ParseMemoryStressParams(r); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			} else {
				api.PostMemoryStress(w, r, *params)
			}
		}),
		http.HandlerFunc(api.CancelMemoryStress),
	))
	return r
}
