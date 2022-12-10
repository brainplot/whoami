package handlers_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/desotech-it/whoami/handlers"
	"github.com/desotech-it/whoami/serialize"
)

var (
	errSerializer = errors.New("serialization failed")
	testBody      = []byte("test")
)

type successfulSerializer struct{}

func (s successfulSerializer) Serialize(v any) ([]byte, error) {
	return testBody, nil
}

type errorfulSerializer struct{}

func (s errorfulSerializer) Serialize(v any) ([]byte, error) {
	return nil, errSerializer
}

func TestSerializerHandler(t *testing.T) {
	testCases := []struct {
		name        string
		serializer  serialize.Serializer
		contentType string
		code        int
		body        []byte
	}{
		{
			name:        "Success",
			serializer:  successfulSerializer{},
			contentType: "test/success",
			code:        http.StatusOK,
			body:        testBody,
		},
		{
			name:        "Failure",
			serializer:  errorfulSerializer{},
			contentType: "text/plain",
			code:        http.StatusInternalServerError,
			body:        []byte(errSerializer.Error()),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			handler := handlers.SerializerHandler{
				StatusCode:  tC.code,
				Serializer:  tC.serializer,
				ContentType: tC.contentType,
			}
			var response *http.Response = nil
			{
				request := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
				response = sendRequestToHandler(request, handler)
			}
			defer response.Body.Close()
			{
				got := response.StatusCode
				want := tC.code
				if got != want {
					t.Errorf("got = %d; want = %d", got, want)
				}
			}
			{
				got := response.Header.Get("Content-Type")
				want := tC.contentType
				if !strings.HasPrefix(got, want) {
					t.Errorf("got = %v; want = %v", got, want)
				}
			}
			{
				got, err := io.ReadAll(response.Body)
				if err != nil {
					t.Error(err)
				} else {
					got = bytes.TrimSpace(got)
				}
				want := tC.body
				if !bytes.Equal(got, want) {
					t.Errorf("got = %v; want = %v", string(got), string(want))
				}
			}
		})
	}
}
