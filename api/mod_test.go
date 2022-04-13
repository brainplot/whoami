package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/desotech-it/whoami/api"
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
)

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
