package grpc

import (
	"github.com/gojek/fiber"
	"github.com/gojek/fiber/protocol"
	"google.golang.org/grpc/metadata"
)

type Request struct {
	// Metadata will hold the grpc headers for request
	Metadata metadata.MD
	Message  []byte
}

func (r *Request) Protocol() protocol.Protocol {
	return protocol.GRPC
}

func (r *Request) Payload() []byte {
	return r.Message
}

func (r *Request) Header() map[string][]string {
	return r.Metadata
}

func (r *Request) Clone() (fiber.Request, error) {
	var copiedMessage []byte
	if len(r.Message) > 0 {
		copiedMessage = make([]byte, len(r.Message))
		copy(copiedMessage, r.Message)
	}
	return &Request{
		Metadata: r.Metadata,
		Message:  copiedMessage,
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
