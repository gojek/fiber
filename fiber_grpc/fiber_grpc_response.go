package fiber_grpc

import (
	"github.com/gojek/fiber"
	"github.com/gojek/fiber/errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type Response struct {
	Metadata        metadata.MD
	ResponsePayload proto.Message
	error           errors.HTTPError //TODO should rename this to generic error
}

func (r *Response) IsSuccess() bool {
	return true
}

func (r *Response) Payload() interface{} {
	return r.ResponsePayload
}

func (r *Response) StatusCode() int {
	return 0
}

func (r *Response) BackendName() string {
	return "grpc"
}

func (r *Response) WithBackendName(string) fiber.Response {
	return r
}
