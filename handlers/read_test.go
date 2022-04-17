package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desotech-it/whoami/handlers"
)

func TestReadHandler(t *testing.T) {
	testCases := []struct {
		method string
		want   int
	}{
		{
			method: http.MethodGet,
			want:   http.StatusNotFound,
		},
		{
			method: http.MethodHead,
			want:   http.StatusNotFound,
		},
		{
			method: http.MethodOptions,
			want:   http.StatusOK,
		},
		{
			method: http.MethodPost,
			want:   http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodPut,
			want:   http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodDelete,
			want:   http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodConnect,
			want:   http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodTrace,
			want:   http.StatusMethodNotAllowed,
		},
		{
			method: http.MethodPatch,
			want:   http.StatusMethodNotAllowed,
		},
	}
	handler := handlers.ReadHandler(http.NotFoundHandler())
	for _, tC := range testCases {
		t.Run(tC.method, func(t *testing.T) {
			request := httptest.NewRequest(tC.method, "/", http.NoBody)
			response := sendRequestToHandler(request, handler)
			defer response.Body.Close()
			got := response.StatusCode
			want := tC.want
			if got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
		})
	}
}
