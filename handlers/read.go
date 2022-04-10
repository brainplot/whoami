package handlers

import (
	"net/http"

	ghandlers "github.com/gorilla/handlers"
)

func ReadHandler(handler http.Handler) http.Handler {
	result := ghandlers.MethodHandler{}
	result[http.MethodGet] = handler
	result[http.MethodHead] = handler
	return result
}
