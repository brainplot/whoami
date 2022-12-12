package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/desotech-it/whoami/api"
)

func TestDecodeRequestBody(t *testing.T) {
	type data struct {
		MyKey string `json:"myKey" form:"myKey"`
	}
	testCases := []struct {
		name        string
		contentType string
		requestBody io.Reader
		decoded     data
		err         error
	}{
		{
			name:        "JSON",
			contentType: "application/json",
			requestBody: strings.NewReader(`{"myKey":"myValue"}`),
			decoded:     data{MyKey: "myValue"},
		},
		{
			name:        "Form",
			contentType: "application/x-www-form-urlencoded",
			requestBody: strings.NewReader("myKey=myValue"),
			decoded:     data{MyKey: "myValue"},
		},
		{
			name:        "unknown",
			contentType: "unknown/unknown",
			requestBody: http.NoBody,
			err:         api.ErrUnknownRequestFormat,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", tC.requestBody)
			request.Header.Set("Content-Type", tC.contentType)
			decoded := data{}
			if err := api.DecodeRequestBody(request, &decoded); err != tC.err {
				t.Errorf("got = %v; want = %v", err, tC.err)
			}
			if got, want := decoded, tC.decoded; got != want {
				t.Errorf("got = %v; want = %v", got, want)
			}
		})
	}
}
