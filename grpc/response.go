package grpc

import (
	"strings"

	"github.com/gojek/fiber"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Response struct {
	Metadata metadata.MD
	Message  []byte
	Status   status.Status
}

func (r *Response) IsSuccess() bool {
	return r.StatusCode() == int(codes.OK)
}

func (r *Response) Payload() []byte {
	return r.Message
}

func (r *Response) StatusCode() int {
	return int(r.Status.Code())
}

// Attribute returns all the values associated with the given key, in the response metadata.
// If the key does not exist, an empty slice will be returned.
func (r *Response) Attribute(key string) []string {
	return r.Metadata.Get(key)
}

// WithAttribute appends the given value(s) to the key, in the response metadata.
// If the key does not already exist, a new key will be created.
// The modified response is returned.
func (r *Response) WithAttribute(key string, values ...string) fiber.Response {
	r.Metadata.Append(key, values...)
	return r
}

func (r *Response) BackendName() string {
	return strings.Join(r.Attribute("backend"), ",")
}

// WithBackendName sets the given backend name in the response metadata.
// The modified response is returned.
func (r *Response) WithBackendName(backendName string) fiber.Response {
	r.Metadata.Set("backend", backendName)
	return r
}
