package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrUnknownRequestFormat = errors.New("unknown request format")
)

type errorResponse struct {
	Message string `json:"message"`
}

func renderError(w http.ResponseWriter, r *http.Request, message string, code int) {
	render.Status(r, code)
	render.JSON(w, r, errorResponse{message})
}

func DecodeRequestBody(r *http.Request, v interface{}) error {
	switch render.GetRequestContentType(r) {
	case render.ContentTypeJSON:
		return render.DecodeJSON(r.Body, v)
	case render.ContentTypeForm:
		return render.DecodeForm(r.Body, v)
	default:
		return ErrUnknownRequestFormat
	}
}
