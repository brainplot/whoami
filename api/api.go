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
	"github.com/desotech-it/whoami/status"
	"github.com/desotech-it/whoami/version"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Api interface {
	// GET /version
	GetVersion(http.ResponseWriter, *http.Request)
	// GET /health
	GetHealth(http.ResponseWriter, *http.Request)
	// PUT /health
	PutHealth(http.ResponseWriter, *http.Request)
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

var (
	defaultVersion               *version.Info                = &version.BuildInfo
	defaultInstanceStatus        *status.InstanceStatus       = &status.Current
	defaultVirtualMemoryProvider memory.VirtualMemoryProvider = memory.VirtualMemoryProviderFunc(memory.VirtualMemoryWithContext)
	defaultInterfacesProvider    net.InterfacesProvider       = net.InterfacesProviderFunc(net.Interfaces)
	defaultHostnameProvider      os.HostnameProvider          = os.HostnameProviderFunc(os.Hostname)
)

type Server struct {
	Version                  *version.Info
	InstanceStatus           *status.InstanceStatus
	VirtualMemoryProvider    memory.VirtualMemoryProvider
	InterfacesProvider       net.InterfacesProvider
	HostnameProvider         os.HostnameProvider
	memoryStressSession      *MemoryStressSession
	memoryStressSessionMutex sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		Version:               defaultVersion,
		InstanceStatus:        defaultInstanceStatus,
		VirtualMemoryProvider: defaultVirtualMemoryProvider,
		InterfacesProvider:    defaultInterfacesProvider,
		HostnameProvider:      defaultHostnameProvider,
	}
}

func NewServerWithVersion(versionInfo *version.Info) *Server {
	if versionInfo == nil {
		panic("versionInfo must not be nil")
	}
	return &Server{
		Version:               versionInfo,
		InstanceStatus:        defaultInstanceStatus,
		VirtualMemoryProvider: defaultVirtualMemoryProvider,
		InterfacesProvider:    defaultInterfacesProvider,
		HostnameProvider:      defaultHostnameProvider,
	}
}

func (s *Server) GetVersion(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, s.Version)
}

func (s *Server) GetHealth(w http.ResponseWriter, r *http.Request) {
	var httpStatus int
	health := s.InstanceStatus.Health
	if health == status.Up {
		httpStatus = http.StatusOK
	} else {
		httpStatus = http.StatusServiceUnavailable
	}
	render.Status(r, httpStatus)
	render.JSON(w, r, Status{health.String()})
}

func (s *Server) PutHealth(w http.ResponseWriter, r *http.Request) {
	form := Status{}
	if err := DecodeRequestBody(r, &form); err != nil {
		renderError(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	if status, err := status.Parse(form.Status); err != nil {
		renderError(w, r, err.Error(), http.StatusBadRequest)
	} else {
		s.InstanceStatus.Health = status
		render.JSON(w, r, Status{status.String()})
	}
}

func (s *Server) GetMemory(w http.ResponseWriter, r *http.Request) {
	if vm, err := s.VirtualMemoryProvider.VirtualMemory(r.Context()); err != nil {
		next := func(w http.ResponseWriter, r *http.Request) {
			renderError(w, r, err.Error(), http.StatusInternalServerError)
		}
		handlers.CancellableHandler(err, http.HandlerFunc(next)).ServeHTTP(w, r)
	} else {
		render.JSON(w, r, vm)
	}
}

func (s *Server) GetInterfaces(w http.ResponseWriter, r *http.Request) {
	if netInterfaces, err := s.InterfacesProvider.Interfaces(); err != nil {
		renderError(w, r, err.Error(), http.StatusInternalServerError)
	} else {
		render.JSON(w, r, netInterfaces)
	}
}

func (s *Server) GetHostname(w http.ResponseWriter, r *http.Request) {
	if hostname, err := s.HostnameProvider.Hostname(); err != nil {
		renderError(w, r, err.Error(), http.StatusInternalServerError)
	} else {
		render.JSON(w, r, hostname)
	}
}

func (s *Server) EchoRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	clientRequest := client.RequestFromStdRequest(r)
	render.JSON(w, r, clientRequest)
}

func (s *Server) GetMemoryStress(w http.ResponseWriter, r *http.Request) {
	s.memoryStressSessionMutex.RLock()
	if s.memoryStressSession == nil {
		s.memoryStressSessionMutex.RUnlock()
		renderError(w, r, "memory stress is not currently running", http.StatusBadRequest)
	} else {
		render.JSON(w, r, MemoryStressSession{
			StressSession: StressSession{
				StartedAt:  s.memoryStressSession.StartedAt,
				FinishedAt: s.memoryStressSession.FinishedAt,
			},
			BytesAllocated: s.memoryStressSession.stresser.BytesAllocated(),
		})
		s.memoryStressSessionMutex.RUnlock()
	}
}

func (s *Server) PostMemoryStress(w http.ResponseWriter, r *http.Request, params memory.StressParameters) {
	s.memoryStressSessionMutex.Lock()
	if s.memoryStressSession == nil {
		if stresser, err := memory.NewStresser(params); err != nil {
			s.memoryStressSessionMutex.Unlock()
			renderError(w, r, err.Error(), http.StatusBadRequest)
		} else {
			ctx, cancelFunc := context.WithCancel(context.Background())
			s.memoryStressSession = &MemoryStressSession{
				StressSession: StressSession{
					StartedAt: time.Now().UTC(),
				},
				cancelFunc: cancelFunc,
				stresser:   stresser,
			}
			render.JSON(w, r, MemoryStressSession{
				StressSession: StressSession{
					StartedAt: s.memoryStressSession.StartedAt,
				},
				BytesAllocated: s.memoryStressSession.stresser.BytesAllocated(),
			})
			s.memoryStressSessionMutex.Unlock()
			go stresser.Stress(ctx)
		}
	} else {
		s.memoryStressSessionMutex.Unlock()
		renderError(w, r, "memory stress is already running", http.StatusBadRequest)
	}
}

func (s *Server) CancelMemoryStress(w http.ResponseWriter, r *http.Request) {
	s.memoryStressSessionMutex.Lock()
	if s.memoryStressSession == nil {
		s.memoryStressSessionMutex.Unlock()
		renderError(w, r, "memory stress is not currently running", http.StatusBadRequest)
	} else {
		cancelFunc := s.memoryStressSession.cancelFunc
		finishedAt := time.Now().UTC()
		s.memoryStressSession.FinishedAt = &finishedAt
		cancelFunc()
		render.JSON(w, r, MemoryStressSession{
			StressSession: StressSession{
				StartedAt:  s.memoryStressSession.StartedAt,
				FinishedAt: s.memoryStressSession.FinishedAt,
			},
			BytesAllocated: s.memoryStressSession.stresser.BytesAllocated(),
		})
		s.memoryStressSession = nil
		s.memoryStressSessionMutex.Unlock()
	}
}

func handleGetHead(r *chi.Mux, pattern string, getHandler http.HandlerFunc) {
	r.Get(pattern, getHandler)
	r.Head(pattern, getHandler)
}

func handleGetHeadPostDelete(r *chi.Mux, pattern string, getHandler, postHandler, deleteHandler http.HandlerFunc) {
	handleGetHead(r, pattern, getHandler)
	r.Post(pattern, postHandler)
	r.Delete(pattern, deleteHandler)
}

func handleGetHeadPut(r *chi.Mux, pattern string, getHandler, putHandler http.HandlerFunc) {
	handleGetHead(r, pattern, getHandler)
	r.Put(pattern, putHandler)
}

func requestHasNoBody(r *http.Request) bool {
	contentLength := r.Header.Get("Content-Length")
	return len(contentLength) == 0 || contentLength == "0"
}

func Handler(api Api) http.Handler {
	r := chi.NewRouter()
	handleGetHead(r, "/version", api.GetVersion)
	handleGetHeadPut(r, "/health", api.GetHealth, api.PutHealth)
	handleGetHead(r, "/memory", api.GetMemory)
	handleGetHead(r, "/interfaces", api.GetInterfaces)
	handleGetHead(r, "/hostname", api.GetHostname)
	r.HandleFunc("/request", api.EchoRequest)
	// TODO test this endpoint better (different parameters etc...)
	handleGetHeadPostDelete(r, "/memory/stresssession",
		api.GetMemoryStress,
		func(w http.ResponseWriter, r *http.Request) {
			form := memory.StressParameters{}
			if !requestHasNoBody(r) {
				if err := DecodeRequestBody(r, &form); err != nil {
					renderError(w, r, err.Error(), http.StatusBadRequest)
					return
				}
			}
			api.PostMemoryStress(w, r, form)
		},
		api.CancelMemoryStress,
	)
	return r
}
