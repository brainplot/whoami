package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
