package client

import (
	"net/http"
	"net/url"
)

type Request struct {
	// Method specifies the HTTP method (GET, POST, PUT, etc.).
	Method string `json:"method"`

	// Path specifies the request's path
	Path string `json:"path"`

	// The protocol version for incoming requests.
	Proto      string `json:"proto"`      // "HTTP/1.0"
	ProtoMajor int    `json:"protoMajor"` // 1
	ProtoMinor int    `json:"protoMinor"` // 0

	// Header contains the request header fields either received
	// by the server or to be sent by the client.
	//
	// If a server received a request with header lines,
	//
	//	Host: example.com
	//	accept-encoding: gzip, deflate
	//	Accept-Language: en-us
	//	fOO: Bar
	//	foo: two
	//
	// then
	//
	//	Header = map[string][]string{
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Foo": {"Bar", "two"},
	//	}
	//
	// For incoming requests, the Host header is promoted to the
	// Request.Host field and removed from the Header map.
	//
	// HTTP defines that header names are case-insensitive. The
	// request parser implements this by using CanonicalHeaderKey,
	// making the first character and any characters following a
	// hyphen uppercase and the rest lowercase.
	Headers http.Header `json:"headers"`

	// Form contains the parsed form data, including both the URL
	// field's query parameters and the PATCH, POST, or PUT form data.
	Form url.Values `json:"form,omitempty"`

	// Trailer specifies additional headers that are sent after the request
	// body.
	//
	// For server requests, the Trailer map initially contains only the
	// trailer keys, with nil values. (The client declares which trailers it
	// will later send.)  While the handler is reading from Body, it must
	// not reference Trailer. After reading from Body returns EOF, Trailer
	// can be read again and will contain non-nil values, if they were sent
	// by the client.
	TrailerHeaders http.Header `json:"trailerHeaders,omitempty"`

	// For server requests, Host specifies the host on which the
	// URL is sought. For HTTP/1 (per RFC 7230, section 5.4), this
	// is either the value of the "Host" header or the host name
	// given in the URL itself. For HTTP/2, it is the value of the
	// ":authority" pseudo-header field.
	// It may be of the form "host:port". For international domain
	// names, Host may be in Punycode or Unicode form. Use
	// golang.org/x/net/idna to convert it to either format if
	// needed.
	// To prevent DNS rebinding attacks, server Handlers should
	// validate that the Host header has a value for which the
	// Handler considers itself authoritative. The included
	// ServeMux supports patterns registered to particular host
	// names and thus protects its registered Handlers.
	Host string `json:"host"`

	// RemoteAddr allows HTTP servers and other software to record
	// the network address that sent the request, usually for
	// logging.
	RemoteAddr string `json:"remoteAddress"`
}

func RequestFromStdRequest(r *http.Request) *Request {
	return &Request{
		Method:         r.Method,
		Path:           r.URL.Path,
		Proto:          r.Proto,
		ProtoMajor:     r.ProtoMajor,
		ProtoMinor:     r.ProtoMinor,
		Headers:        r.Header,
		Form:           r.Form,
		TrailerHeaders: r.Trailer,
		Host:           r.Host,
		RemoteAddr:     r.RemoteAddr,
	}
}
