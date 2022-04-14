package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desotech-it/whoami/handlers"
)

func passIfPanicWithError(err error, t *testing.T) {
	got := err
	want := recover()
	if got != want {
		t.Errorf("got = %v; want = %v", got, want)
	}
}

func TestCancellableHandler(t *testing.T) {
	t.Run("Canceled=AbortHandler", func(t *testing.T) {
		handler := handlers.CancellableHandler(context.Canceled, nil)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		defer passIfPanicWithError(http.ErrAbortHandler, t)
		handler.ServeHTTP(recorder, request)
	})

	t.Run("NotCanceled=NextHandler", func(t *testing.T) {
		handler := handlers.CancellableHandler(nil, http.NotFoundHandler())
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		defer passIfPanicWithError(nil, t)
		handler.ServeHTTP(recorder, request)
		response := recorder.Result()
		defer response.Body.Close()
		got := response.StatusCode
		want := http.StatusNotFound
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
}
