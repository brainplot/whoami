package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type errorResponse struct {
	Message string `json:"message"`
}

func renderError(w http.ResponseWriter, r *http.Request, message string, code int) {
	render.Status(r, code)
	render.JSON(w, r, errorResponse{message})
}
