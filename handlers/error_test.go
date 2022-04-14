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
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	t.Run("Code=Input", func(t *testing.T) {
		got := response.StatusCode
		want := errorHandlerCode
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("Body=ErrorMessage", func(t *testing.T) {
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
	})
}
