package handlers

import (
	"bytes"
	"net/http"
	"time"

	"github.com/desotech-it/whoami/serialize"
)

type SerializeHandler struct {
	Payload     any
	Serializer  serialize.Serializer
	ContentType string
}

func (h SerializeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if data, err := h.Serializer.Serialize(h.Payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		header := w.Header()
		header.Set("Content-Type", h.ContentType)
		header.Set("X-Content-Type-Options", "nosniff")
		http.ServeContent(w, r, "", time.Time{}, bytes.NewReader(data))
	}
}
