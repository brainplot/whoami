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

func mockVirtualStatMemoryProvider(outcome Outcome) memory.VirtualMemoryProviderFunc {
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
	config := &api.Config{}
	handler := api.Handler(&testVersion, config)
	request := httptest.NewRequest(http.MethodGet, "/version", http.NoBody)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	t.Run("Code=OK", func(t *testing.T) {
		got := response.StatusCode
		want := http.StatusOK
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
	t.Run("Body=TestVersion", func(t *testing.T) {
		got := version.Info{}
		want := &testVersion
		if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
			t.Error(err)
		}
		if got != *want {
			t.Errorf("got = %v; want = %v", got, testVersion)
		}
	})
}

func TestMemory(t *testing.T) {
	t.Run("Outcome=Success", func(t *testing.T) {
		providerFunc := mockVirtualStatMemoryProvider(OutcomeSuccess)
		config := &api.Config{
			VirtualMemoryProvider: memory.VirtualMemoryProviderFunc(providerFunc),
		}
		handler := api.Handler(&testVersion, config)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/memory", http.NoBody)
		handler.ServeHTTP(recorder, request)
		response := recorder.Result()
		defer response.Body.Close()
		t.Run("Code=OK", func(t *testing.T) {
			got := response.StatusCode
			want := http.StatusOK
			if got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
		})
		t.Run("Body=Memory", func(t *testing.T) {
			got := memory.VirtualMemoryStat{}
			want := &testVirtualMemory
			if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
				t.Error(err)
			}
			if got != *want {
				t.Errorf("got = %v; want = %v", got, want)
			}
		})
	})

	t.Run("Outcome=Failure", func(t *testing.T) {
		providerFunc := mockVirtualStatMemoryProvider(OutcomeFailure)
		config := &api.Config{
			VirtualMemoryProvider: memory.VirtualMemoryProviderFunc(providerFunc),
		}
		handler := api.Handler(&testVersion, config)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/memory", http.NoBody)
		handler.ServeHTTP(recorder, request)
		response := recorder.Result()
		defer response.Body.Close()
		t.Run("Code=InternalServerError", func(t *testing.T) {
			got := response.StatusCode
			want := http.StatusInternalServerError
			if got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
		})
	})

	t.Run("Outcome=Canceled", func(t *testing.T) {
		providerFunc := mockVirtualStatMemoryProvider(OutcomeCanceled)
		config := &api.Config{
			VirtualMemoryProvider: memory.VirtualMemoryProviderFunc(providerFunc),
		}
		handler := api.Handler(&testVersion, config)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/memory", http.NoBody)
		defer passIfPanicWithError(http.ErrAbortHandler, t)
		handler.ServeHTTP(recorder, request)
	})
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
	t.Run("Outcome=Success", func(t *testing.T) {
		config := api.Config{
			InterfacesProvider: mockInterfacesProvider(OutcomeSuccess),
		}
		handler := api.Handler(&testVersion, &config)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/interfaces", http.NoBody)
		handler.ServeHTTP(recorder, request)
		response := recorder.Result()
		defer response.Body.Close()
		t.Run("Code=OK", func(t *testing.T) {
			got := response.StatusCode
			want := http.StatusOK
			if got != want {
				t.Errorf("got = %d; want = %d", got, want)
			}
		})
		t.Run("Body=Interfaces", func(t *testing.T) {
			got := []net.Interface{}
			want := testInterfaces
			if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
				t.Error(err)
			}
			if !interfaceSlicesEqual(got, want) {
				t.Errorf("got = %v; want = %v", got, want)
			}
		})
	})

	t.Run("Outcome=Failure", func(t *testing.T) {
		config := api.Config{
			InterfacesProvider: mockInterfacesProvider(OutcomeFailure),
		}
		handler := api.Handler(&testVersion, &config)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/interfaces", http.NoBody)
		handler.ServeHTTP(recorder, request)
		response := recorder.Result()
		defer response.Body.Close()
		got := response.StatusCode
		want := http.StatusInternalServerError
		if got != want {
			t.Errorf("got = %d; want = %d", got, want)
		}
	})
}
