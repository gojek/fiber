package grpc

import (
	"errors"

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
	ResponseProto proto.Message
	// ServiceMethod is the service and method of server point in the format "{grpc_service_name}/{method_name}"
	ServiceMethod string

	// Endpoint is the host+port of the grpc server, eg "127.0.0.1:50050"
	endpoint string
}

func (r *Request) Protocol() fiber.Protocol {
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
		ResponseProto:  r.ResponseProto,
		ServiceMethod:  r.ServiceMethod,
		endpoint:       r.endpoint,
	}, nil
}

// OperationName is naming used in tracing interceptors
func (r *Request) OperationName() string {
	return r.ServiceMethod
}

// Transform is use by backend component within a Proxy to abstract endpoint from dispatcher
func (r *Request) Transform(backend fiber.Backend) (fiber.Request, error) {
	if backend == nil {
		return nil, errors.New("backend cannot be nil")
	}
	r.endpoint = backend.URL("")
	return r, nil
}
