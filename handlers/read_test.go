package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desotech-it/whoami/handlers"
)

func newEmptyBodyRequest(method, target string) *http.Request {
	return httptest.NewRequest(method, target, http.NoBody)
}

func newGetRequest(target string) *http.Request {
	return newEmptyBodyRequest(http.MethodGet, target)
}

func newHeadRequest(target string) *http.Request {
	return newEmptyBodyRequest(http.MethodHead, target)
}

func newOptionsRequest(target string) *http.Request {
	return newEmptyBodyRequest(http.MethodOptions, target)
}

func TestReadHandler(t *testing.T) {
	handler := handlers.ReadHandler(http.NotFoundHandler())
	t.Run("GET=allowed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := newGetRequest("/")
		handler.ServeHTTP(recorder, request)
		got := recorder.Code
		want := http.StatusNotFound
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("HEAD=allowed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := newHeadRequest("/")
		handler.ServeHTTP(recorder, request)
		got := recorder.Code
		want := http.StatusNotFound
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("OPTIONS=allowed", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := newOptionsRequest("/")
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
		request := newEmptyBodyRequest(http.MethodPost, "/")
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