package handlers

import (
	"bytes"
	"net/http"
	"time"

	"github.com/desotech-it/whoami/serialize"
)

type SerializerHandler struct {
	StatusCode  int
	Payload     any
	Serializer  serialize.Serializer
	ContentType string
}

func (h SerializerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if data, err := h.Serializer.Serialize(h.Payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		header := w.Header()
		header.Set("Content-Type", h.ContentType)
		header.Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(h.StatusCode)
		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(data))
	}
}
