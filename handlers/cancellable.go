package handlers

import (
	"context"
	"net/http"
)

type cancellableHandler struct {
	http.Handler
	err error
}

func (h cancellableHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.err == context.Canceled {
		panic(http.ErrAbortHandler)
	}
	h.Handler.ServeHTTP(w, r)
}

func CancellableHandler(err error, handler http.Handler) http.Handler {
	return cancellableHandler{handler, err}
}
