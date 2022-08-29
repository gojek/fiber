package grpc

import (
	"github.com/gojek/fiber"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type Request struct {
	// Metadata will hold the grpc headers for request
	Metadata metadata.MD
	// RequestPayload is the grpc request
	RequestPayload proto.Message
	// RequestPayload is the grpc expected response type
}

func (r *Request) Protocol() string {
	return fiber.GRPC
}

func (r *Request) Payload() interface{} {
	return r.RequestPayload
}

func (r *Request) Header() map[string][]string {
	return r.Metadata
}

func (r *Request) Clone() (fiber.Request, error) {
	return &Request{
		Metadata:       r.Metadata,
		RequestPayload: r.RequestPayload,
	}, nil
}

// OperationName is naming used in tracing interceptors
func (r *Request) OperationName() string {
	// For grpc implementation, serviceMethod and endpoint is init with dispatcher
	return "grpc"
}

// Transform is use by backend component within a Proxy to abstract endpoint from dispatcher
func (r *Request) Transform(_ fiber.Backend) (fiber.Request, error) {
	// For grpc implementation, endpoint is init with dispatcher
	return r, nil
}
