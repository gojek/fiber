package fiber_grpc

import (
	"github.com/gojek/fiber"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type Request struct {
	Metadata       metadata.MD
	RequestPayload proto.Message

	hostport string
}

func (r *Request) Payload() interface{} {
	return r.RequestPayload
}

func (r *Request) Header() map[string][]string {
	return r.Metadata
}

func (r *Request) Clone() (fiber.Request, error) {
	return r, nil
}

func (r *Request) OperationName() string {
	return "grpc"
}

func (r *Request) Transform(backend fiber.Backend) (fiber.Request, error) {
	r.hostport = backend.URL("")
	return r, nil
}
