package main

import (
	"log"
	"net/http"
	"os"

	"github.com/desotech-it/whoami/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func apiHandler() http.Handler {
	return http.StripPrefix("/api", api.Handler(api.NewServer()))
}

func main() {
	errorLog := log.New(os.Stderr, "[whoami] ", log.LstdFlags)
	addr := os.Getenv("WHOAMI_ADDRESS")
	if addr == "" {
		addr = "127.0.0.1:3000"
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Mount("/api/", apiHandler())
	server := http.Server{
		Addr:    addr,
		Handler: r,
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
		ErrorLog: errorLog,
	}
	if err := server.ListenAndServe(); err != nil {
		errorLog.Fatalln(err)
	}
}
