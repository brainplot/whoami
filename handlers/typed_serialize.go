// This file will house all of the "constructors" for known SerializeHandler structs.
//go:build !test

package handlers

import (
	"net/http"

	"github.com/desotech-it/whoami/serialize"
)

func JSONSerializeHandler(v any) http.Handler {
	return SerializeHandler{
		Payload:     v,
		Serializer:  serialize.SerializerJSON,
		ContentType: "application/json; charset=utf-8",
	}
}
