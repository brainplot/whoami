// This file will house all of the "constructors" for known SerializeHandler structs.
//go:!test

package handlers

import (
	"net/http"

	"github.com/desotech-it/whoami/serialize"
)

func JSONSerializerHandler(statusCode int, v any) http.Handler {
	return SerializerHandler{
		StatusCode:  statusCode,
		Payload:     v,
		Serializer:  serialize.SerializerJSON,
		ContentType: "application/json; charset=utf-8",
	}
}
