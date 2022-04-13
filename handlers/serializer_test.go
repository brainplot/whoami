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
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	handler := handlers.SerializerHandler{
		Payload:     nil,
		Serializer:  successfulSerializer{},
		ContentType: "test/success",
	}
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	t.Run("Code=OK", func(t *testing.T) {
		got := response.StatusCode
		want := http.StatusOK
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("Content-Type=test/success", func(t *testing.T) {
		got := response.Header.Get("Content-Type")
		want := "test/success"
		if got != want {
			t.Errorf("got = %q; want = %q", got, want)
		}
	})
	t.Run("Body=test", func(t *testing.T) {
		got, err := io.ReadAll(response.Body)
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
	recoder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	handler := handlers.SerializerHandler{
		Payload:    nil,
		Serializer: errorfulSerializer{},
	}
	handler.ServeHTTP(recoder, request)
	response := recoder.Result()
	t.Run("Code=InternalServerError", func(t *testing.T) {
		got := response.StatusCode
		want := http.StatusInternalServerError
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("Content-Type=text/plain", func(t *testing.T) {
		got := response.Header.Get("Content-Type")
		want := "text/plain; charset=utf-8"
		if got != want {
			t.Errorf("got = %q; want = %q", got, want)
		}
	})
	t.Run("Body=error", func(t *testing.T) {
		got, err := io.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}
		want := []byte("serialization failed\n")
		if !bytes.Equal(got, want) {
			t.Errorf("got = %s; want = %s", got, want)
		}
	})
}
