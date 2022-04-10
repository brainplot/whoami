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

type successfulSerializer struct{}

func (s successfulSerializer) Serialize(v any) ([]byte, error) {
	return []byte("test"), nil
}

type errorfulSerializer struct{}

func (s errorfulSerializer) Serialize(v any) ([]byte, error) {
	return nil, errors.New("serialization failed")
}

func TestSuccessfulSerializationReturnsSuccessResponse(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := newGetRequest("/")
	handler := handlers.SerializeHandler{
		Payload:     nil,
		Serializer:  successfulSerializer{},
		ContentType: "test/success",
	}
	handler.ServeHTTP(responseRecorder, request)
	t.Run("Code=OK", func(t *testing.T) {
		got := responseRecorder.Code
		want := http.StatusOK
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("Content-Type=test/success", func(t *testing.T) {
		got := responseRecorder.Header().Get("Content-Type")
		want := "test/success"
		if got != want {
			t.Errorf("got = %q; want = %q", got, want)
		}
	})
	t.Run("Body=test", func(t *testing.T) {
		got, err := io.ReadAll(responseRecorder.Body)
		if err != nil {
			t.Error(err)
		}
		want := []byte("test")
		if !bytes.Equal(got, want) {
			t.Errorf("got = %q; want = %q", got, want)
		}
	})
}

func TestErrorfulSerializationReturnsErrorResponse(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	request := newGetRequest("/")
	handler := handlers.SerializeHandler{
		Payload:    nil,
		Serializer: errorfulSerializer{},
	}
	handler.ServeHTTP(responseRecorder, request)
	t.Run("Code=InternalServerError", func(t *testing.T) {
		got := responseRecorder.Code
		want := http.StatusInternalServerError
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("Content-Type=text/plain", func(t *testing.T) {
		got := responseRecorder.Header().Get("Content-Type")
		want := "text/plain; charset=utf-8"
		if got != want {
			t.Errorf("got = %q; want = %q", got, want)
		}
	})
	t.Run("Body=error", func(t *testing.T) {
		got, err := io.ReadAll(responseRecorder.Body)
		if err != nil {
			t.Error(err)
		}
		want := []byte("serialization failed\n")
		if !bytes.Equal(got, want) {
			t.Errorf("got = %s; want = %s", got, want)
		}
	})
}
