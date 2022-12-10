package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/desotech-it/whoami/api"
	"github.com/desotech-it/whoami/api/memory"
)

func postRequestWithBody(target string, body io.Reader) *http.Request {
	request := httptest.NewRequest(http.MethodPost, target, body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request
}

func TestParseMemoryStressParams(t *testing.T) {
	testCases := []struct {
		name   string
		values url.Values
		want   memory.StressParameters
	}{
		{
			name:   "EmptyValues",
			values: url.Values{},
			want:   memory.StressParameters{},
		},
		{
			name: "WithValues",
			values: func() url.Values {
				v := url.Values{}
				v.Set(api.HTTPParamInterval, "42s")
				v.Set(api.HTTPParamAllocSize, "256")
				return v
			}(),
			want: memory.StressParameters{Interval: 42 * time.Second, AllocationSize: 256},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			request := postRequestWithBody("/memory/stresssession", strings.NewReader(tC.values.Encode()))
			if err := request.ParseForm(); err != nil {
				t.Error(err)
			}
			params, err := api.ParseMemoryStressParams(request)
			if err != nil {
				t.Error(err)
			}
			if got, want := *params, tC.want; got != want {
				t.Errorf("got = %v; want = %v", got, want)
			}
		})
	}
}

func TestParseMemoryStressParamsError(t *testing.T) {
	testCases := []struct {
		name   string
		values url.Values
	}{
		{
			name: "InvalidDuration",
			values: func() url.Values {
				v := url.Values{}
				v.Set(api.HTTPParamInterval, "test")
				v.Set(api.HTTPParamAllocSize, "256")
				return v
			}(),
		},
		{
			name: "InvalidAllocSize",
			values: func() url.Values {
				v := url.Values{}
				v.Set(api.HTTPParamInterval, "42s")
				v.Set(api.HTTPParamAllocSize, "test")
				return v
			}(),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			request := postRequestWithBody("/memory/stresssession", strings.NewReader(tC.values.Encode()))
			if err := request.ParseForm(); err != nil {
				t.Error(err)
			}
			params, err := api.ParseMemoryStressParams(request)
			if params != nil {
				t.Error(tC.name + " did not return a nil params")
			}
			if err == nil {
				t.Error(tC.name + " did not raise an error")
			}
		})
	}
}
