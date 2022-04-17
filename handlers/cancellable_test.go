package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desotech-it/whoami/handlers"
)

func TestCancellableHandler(t *testing.T) {
	testCases := []struct {
		name        string
		inputError  error
		outputError error
	}{
		{
			name:        "Canceled",
			inputError:  context.Canceled,
			outputError: http.ErrAbortHandler,
		},
		{
			name: "NotCanceled",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			handler := handlers.CancellableHandler(tC.inputError, http.NotFoundHandler())
			request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
			defer passIfPanicWithError(tC.outputError, t)
			response := sendRequestToHandler(request, handler)
			defer response.Body.Close()
		})
	}
}
