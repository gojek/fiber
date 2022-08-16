package grpc

import (
	"github.com/gojek/fiber"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type Request struct {
	Metadata       metadata.MD
	RequestPayload proto.Message
	ResponseProto  proto.Message
	ServiceMethod  string

	hostport string
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
		hostport:       r.hostport,
	}, nil
}

// OperationName is naming used in tracing interceptors
func (r *Request) OperationName() string {
	return r.ServiceMethod
}

// Transform is use by backend component within a Proxy to abstract endpoint from dispatcher
func (r *Request) Transform(backend fiber.Backend) (fiber.Request, error) {
	r.hostport = backend.URL("")
	return r, nil
}
