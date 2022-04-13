package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desotech-it/whoami/handlers"
)

func TestReadHandler(t *testing.T) {
	handler := handlers.ReadHandler(http.NotFoundHandler())
	t.Run("GET=allowed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		handler.ServeHTTP(recorder, request)
		got := recorder.Code
		want := http.StatusNotFound
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("HEAD=allowed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodHead, "/", http.NoBody)
		handler.ServeHTTP(recorder, request)
		got := recorder.Code
		want := http.StatusNotFound
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("OPTIONS=allowed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodOptions, "/", http.NoBody)
		handler.ServeHTTP(recorder, request)
		t.Run("Code=200", func(t *testing.T) {
			got := recorder.Code
			want := http.StatusOK
			if got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
		})
		t.Run("Allow=GET,HEAD", func(t *testing.T) {
			got := recorder.Header().Get("Allow")
			want := "GET, HEAD"
			if got != want {
				t.Errorf("got = %q; want = %q", got, want)
			}
		})
	})
	t.Run("POST=disallowed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		handler.ServeHTTP(recorder, request)
		t.Run("Code=405", func(t *testing.T) {
			got := recorder.Code
			want := http.StatusMethodNotAllowed
			if got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
		})
	})
}
