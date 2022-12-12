package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/desotech-it/whoami/api"
	"github.com/desotech-it/whoami/api/memory"
	"github.com/desotech-it/whoami/api/net"
	"github.com/desotech-it/whoami/api/os"
	"github.com/desotech-it/whoami/status"
	"github.com/desotech-it/whoami/version"
)

var (
	testVersion = version.Info{
		Major:     1,
		Minor:     2,
		Patch:     3,
		Commit:    "feedcoffee",
		BuildDate: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	testVirtualMemory = memory.VirtualMemoryStat{
		Total:       42,
		Available:   42,
		Used:        42,
		UsedPercent: 42.0,
	}

	testInterfaces = []net.Interface{
		{
			Index: 42,
			MTU:   1234,
			Name:  "test",
			Addresses: []net.Addr{
				{
					Network: "test+net",
					Value:   "10.0.0.1",
				},
			},
		},
	}

	testHostname = os.HostnameInfo{
		Hostname: "foobar",
	}

	errTest           = errors.New("test error")
	errUnknownOutcome = errors.New("unknown outcome")
)

func newServer() *api.Server {
	return api.NewServer(&testVersion)
}

type Outcome uint

const (
	OutcomeFailure  Outcome = iota
	OutcomeSuccess  Outcome = iota
	OutcomeCanceled Outcome = iota
)

func mockVirtualMemoryProvider(outcome Outcome) memory.VirtualMemoryProviderFunc {
	switch outcome {
	case OutcomeFailure:
		return func(context.Context) (*memory.VirtualMemoryStat, error) {
			return nil, errTest
		}
	case OutcomeSuccess:
		return func(context.Context) (*memory.VirtualMemoryStat, error) {
			return &testVirtualMemory, nil
		}
	case OutcomeCanceled:
		return func(ctx context.Context) (*memory.VirtualMemoryStat, error) {
			toBeCanceledCtx, cancelFunc := context.WithCancel(ctx)
			cancelFunc()
			return nil, toBeCanceledCtx.Err()
		}
	default:
		panic(errUnknownOutcome)
	}
}

func mockInterfacesProvider(outcome Outcome) net.InterfacesProviderFunc {
	switch outcome {
	case OutcomeFailure:
		return func() ([]net.Interface, error) {
			return nil, errTest
		}
	case OutcomeSuccess:
		return func() ([]net.Interface, error) {
			return testInterfaces, nil
		}
	default:
		panic(errUnknownOutcome)
	}
}

func mockHostnameProvider(outcome Outcome) os.HostnameProviderFunc {
	switch outcome {
	case OutcomeFailure:
		return func() (*os.HostnameInfo, error) {
			return nil, errTest
		}
	case OutcomeSuccess:
		return func() (*os.HostnameInfo, error) {
			return &testHostname, nil
		}
	default:
		panic(errUnknownOutcome)
	}
}

func passIfPanicWithError(err error, t *testing.T) {
	got := err
	want := recover()
	if got != want {
		t.Errorf("got = %v; want = %v", got, want)
	}
}

func sendRequestToHandler(request *http.Request, handler http.Handler) *http.Response {
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder.Result()
}

func TestVersion(t *testing.T) {
	handler := api.Handler(newServer())
	request := httptest.NewRequest(http.MethodGet, "/version", http.NoBody)
	response := sendRequestToHandler(request, handler)
	defer response.Body.Close()
	{
		got := response.StatusCode
		want := http.StatusOK
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	}
	{
		got := version.Info{}
		want := &testVersion
		if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
			t.Error(err)
		}
		if got != *want {
			t.Errorf("got = %v; want = %v", got, testVersion)
		}
	}
}

func TestGetHealth(t *testing.T) {
	testCases := []struct {
		name               string
		instanceStatus     status.InstanceStatus
		responseStatusCode int
		responseBody       api.Status
	}{
		{
			name:               "up",
			instanceStatus:     status.InstanceStatus{Health: status.Up},
			responseStatusCode: http.StatusOK,
			responseBody:       api.Status{status.Up.String()},
		},
		{
			name:               "down",
			instanceStatus:     status.InstanceStatus{Health: status.Down},
			responseStatusCode: http.StatusServiceUnavailable,
			responseBody:       api.Status{status.Down.String()},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			server := newServer()
			server.InstanceStatus = &tC.instanceStatus
			handler := api.Handler(server)
			request := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
			response := sendRequestToHandler(request, handler)
			defer response.Body.Close()
			if got, want := response.StatusCode, tC.responseStatusCode; got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
			responseBody := api.Status{}
			if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
				t.Error(err)
			}
			if got, want := responseBody, tC.responseBody; got != want {
				t.Errorf("got = %v; want = %v", got, want)
			}
		})
	}
}

func TestPutHealth(t *testing.T) {
	testInstanceStatus := status.InstanceStatus{}
	testCases := []struct {
		statusCode   int
		requestBody  io.Reader
		responseBody api.Status
	}{
		{
			statusCode:   http.StatusOK,
			requestBody:  strings.NewReader("status=up"),
			responseBody: api.Status{status.Up.String()},
		},
		{
			statusCode:   http.StatusOK,
			requestBody:  strings.NewReader("status=down"),
			responseBody: api.Status{status.Down.String()},
		},
	}
	for _, tC := range testCases {
		server := newServer()
		server.InstanceStatus = &testInstanceStatus
		handler := api.Handler(server)
		request := httptest.NewRequest(http.MethodPut, "/health", tC.requestBody)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := sendRequestToHandler(request, handler)
		defer response.Body.Close()
		if got, want := response.StatusCode, tC.statusCode; got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
		body := api.Status{}
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if got, want := body, tC.responseBody; got != want {
			t.Errorf("got = %v; want = %v", got, want)
		}
	}
}

func TestMemory(t *testing.T) {
	testCases := []struct {
		name     string
		provider memory.VirtualMemoryProvider
		code     int
		body     *memory.VirtualMemoryStat
		err      error
	}{
		{
			name:     "Success",
			provider: mockVirtualMemoryProvider(OutcomeSuccess),
			code:     http.StatusOK,
			body:     &testVirtualMemory,
		},
		{
			name:     "Failure",
			provider: mockVirtualMemoryProvider(OutcomeFailure),
			code:     http.StatusInternalServerError,
		},
		{
			name:     "Canceled",
			provider: mockVirtualMemoryProvider(OutcomeCanceled),
			err:      http.ErrAbortHandler,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			server := newServer()
			server.VirtualMemoryProvider = tC.provider
			handler := api.Handler(server)
			request := httptest.NewRequest(http.MethodGet, "/memory", http.NoBody)
			defer passIfPanicWithError(tC.err, t)
			response := sendRequestToHandler(request, handler)
			defer response.Body.Close()
			if got, want := response.StatusCode, tC.code; got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
			if response.StatusCode == http.StatusOK {
				body := memory.VirtualMemoryStat{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if got, want := body, *tC.body; got != want {
					t.Errorf("got = %v; want = %v", got, want)
				}
			}
		})
	}
}

func interfaceSlicesEqual(a, b []net.Interface) bool {
	if len(a) != len(b) {
		return false
	}
	for i, netInterface := range a {
		if !netInterface.Equal(&b[i]) {
			return false
		}
	}
	return true
}

func TestInterfaces(t *testing.T) {
	testCases := []struct {
		name     string
		provider net.InterfacesProvider
		code     int
		body     []net.Interface
	}{
		{
			name:     "Success",
			provider: mockInterfacesProvider(OutcomeSuccess),
			code:     http.StatusOK,
			body:     testInterfaces,
		},
		{
			name:     "Failure",
			provider: mockInterfacesProvider(OutcomeFailure),
			code:     http.StatusInternalServerError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			server := newServer()
			server.InterfacesProvider = tC.provider
			handler := api.Handler(server)
			request := httptest.NewRequest(http.MethodGet, "/interfaces", http.NoBody)
			response := sendRequestToHandler(request, handler)
			defer response.Body.Close()
			if got, want := response.StatusCode, tC.code; got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
			if response.StatusCode == http.StatusOK {
				var got []net.Interface
				want := tC.body
				if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
					t.Error(err)
				} else {
					if !interfaceSlicesEqual(got, want) {
						t.Errorf("got = %v; want = %v", got, want)
					}
				}
			}
		})
	}
}
func TestHostname(t *testing.T) {
	testCases := []struct {
		name     string
		provider os.HostnameProvider
		code     int
		body     os.HostnameInfo
	}{
		{
			name:     "Success",
			provider: mockHostnameProvider(OutcomeSuccess),
			code:     http.StatusOK,
			body:     testHostname,
		},
		{
			name:     "Failure",
			provider: mockHostnameProvider(OutcomeFailure),
			code:     http.StatusInternalServerError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			server := newServer()
			server.HostnameProvider = tC.provider
			handler := api.Handler(server)
			request := httptest.NewRequest(http.MethodGet, "/hostname", http.NoBody)
			response := sendRequestToHandler(request, handler)
			defer response.Body.Close()
			if got, want := response.StatusCode, tC.code; got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
			if response.StatusCode == http.StatusOK {
				var got os.HostnameInfo
				want := tC.body
				if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
					t.Error(err)
				} else {
					if got != want {
						t.Errorf("got = %v; want = %v", got, want)
					}
				}
			}
		})
	}
}

func TestMemoryStress(t *testing.T) {
	type methodStatusPair struct {
		method     string
		statusCode int
	}
	testCases := []struct {
		name           string
		methodSequence []methodStatusPair
	}{
		{
			name: "PostGetDelete",
			methodSequence: []methodStatusPair{
				{http.MethodPost, http.StatusOK},
				{http.MethodGet, http.StatusOK},
				{http.MethodDelete, http.StatusOK},
			},
		},
		{
			name: "PostDeleteGet",
			methodSequence: []methodStatusPair{
				{http.MethodPost, http.StatusOK},
				{http.MethodDelete, http.StatusOK},
				{http.MethodGet, http.StatusBadRequest},
			},
		},
		{
			name: "GetPostDelete",
			methodSequence: []methodStatusPair{
				{http.MethodGet, http.StatusBadRequest},
				{http.MethodPost, http.StatusOK},
				{http.MethodDelete, http.StatusOK},
			},
		},
		{
			name: "GetDeletePostDelete",
			methodSequence: []methodStatusPair{
				{http.MethodGet, http.StatusBadRequest},
				{http.MethodDelete, http.StatusBadRequest},
				{http.MethodPost, http.StatusOK},
				{http.MethodDelete, http.StatusOK},
			},
		},
		{
			name: "DeletePostGet",
			methodSequence: []methodStatusPair{
				{http.MethodDelete, http.StatusBadRequest},
				{http.MethodPost, http.StatusOK},
				{http.MethodGet, http.StatusOK},
				{http.MethodDelete, http.StatusOK},
			},
		},
		{
			name: "DeleteGetPostDelete",
			methodSequence: []methodStatusPair{
				{http.MethodDelete, http.StatusBadRequest},
				{http.MethodGet, http.StatusBadRequest},
				{http.MethodPost, http.StatusOK},
				{http.MethodDelete, http.StatusOK},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			handler := api.Handler(newServer())
			for _, pair := range tC.methodSequence {
				request := httptest.NewRequest(pair.method, "/memory/stresssession", http.NoBody)
				response := sendRequestToHandler(request, handler)
				defer response.Body.Close()
				if got, want := response.StatusCode, pair.statusCode; got != want {
					t.Errorf("got = %d; want = %d", got, want)
				}
			}
		})
	}
}
