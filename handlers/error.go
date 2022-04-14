package handlers

import "net/http"

func ErrorHandler(error string, code int) http.Handler {
	fn := func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, error, code)
	}
	return http.HandlerFunc(fn)
}
