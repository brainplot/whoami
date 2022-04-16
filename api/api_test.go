package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/desotech-it/whoami/api"
	"github.com/desotech-it/whoami/api/memory"
	"github.com/desotech-it/whoami/api/net"
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

	errTest           = errors.New("test error")
	errUnknownOutcome = errors.New("unknown outcome")
)

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

func passIfPanicWithError(err error, t *testing.T) {
	got := err
	want := recover()
	if got != want {
		t.Errorf("got = %v; want = %v", got, want)
	}
}

func TestVersion(t *testing.T) {
	var response *http.Response = nil
	{
		server := api.NewServer(&testVersion)
		handler := api.Handler(server)
		request := httptest.NewRequest(http.MethodGet, "/version", http.NoBody)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		response = recorder.Result()
	}
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

func TestMemory(t *testing.T) {
	testCases := []struct {
		name     string
		provider memory.VirtualMemoryProvider
		code     int
		err      error
	}{
		{
			name:     "Success",
			provider: mockVirtualMemoryProvider(OutcomeSuccess),
			code:     http.StatusOK,
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
			var response *http.Response = nil
			{
				server := api.NewServer(&testVersion)
				server.VirtualMemoryProvider = tC.provider
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodGet, "/memory", http.NoBody)
				handler := api.Handler(server)
				defer passIfPanicWithError(tC.err, t)
				handler.ServeHTTP(recorder, request)
				response = recorder.Result()
			}
			defer response.Body.Close()
			got := response.StatusCode
			want := tC.code
			if got != want {
				t.Errorf("got = %d; want = %d", got, want)
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
			var response *http.Response = nil
			{
				server := api.NewServer(&testVersion)
				server.InterfacesProvider = tC.provider
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodGet, "/interfaces", http.NoBody)
				handler := api.Handler(server)
				handler.ServeHTTP(recorder, request)
				response = recorder.Result()
			}
			defer response.Body.Close()
			{
				got := response.StatusCode
				want := tC.code
				if got != want {
					t.Errorf("got = %d; want = %d", got, want)
				}
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
