package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/desotech-it/whoami/handlers"
)

var (
	errTestError     = errors.New("handlers: error")
	errorHandlerCode = http.StatusInternalServerError
)

func TestErrorHandler(t *testing.T) {
	handler := handlers.ErrorHandler(errTestError.Error(), errorHandlerCode)
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	response := sendRequestToHandler(request, handler)
	defer response.Body.Close()
	{
		got := response.StatusCode
		want := errorHandlerCode
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	}
	{
		got, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		} else {
			got = bytes.TrimSpace(got)
		}
		want := []byte(errTestError.Error())
		if !bytes.Equal(got, want) {
			t.Errorf("got = %s; want = %s", string(got), string(want))
		}
	}
}
