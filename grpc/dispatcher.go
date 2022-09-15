package grpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gojek/fiber"
	fiberError "github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func init() {
	encoding.RegisterCodec(FiberCodec{})
}

const (
	TimeoutDefault = time.Second
)

type Dispatcher struct {
	timeout time.Duration
	// serviceMethod is the service and method of server point in the format "{grpc_service_name}/{method_name}"
	serviceMethod string
	// endpoint is the host+port of the grpc server, eg "127.0.0.1:50050"
	endpoint string
	// conn is the grpc connection dialed upon creation of dispatcher
	conn *grpc.ClientConn
}

type DispatcherConfig struct {
	Service  string
	Method   string
	Endpoint string
	Timeout  time.Duration
}

func (d *Dispatcher) Do(request fiber.Request) fiber.Response {
	grpcRequest, ok := request.(*Request)
	if !ok {
		return fiber.NewErrorResponse(
			fiberError.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "fiber: grpc dispatcher: only grpc.Request type of requests are supported",
			})
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, grpcRequest.Metadata)

	response := new(bytes.Buffer)
	var responseHeader metadata.MD

	// Dispatcher will send both request and payload as bytes, with the use of codec
	// to prevent marshaling. The codec content type will be sent with request and
	// the server will attempt to unmarshal with the codec.
	err := d.conn.Invoke(
		ctx,
		d.serviceMethod,
		grpcRequest.Payload(),
		response,
		grpc.Header(&responseHeader),
		grpc.CallContentSubtype(CodecName),
	)
	if err != nil {
		// if ok is false, unknown codes.Unknown and Status msg is returned in Status
		responseStatus, _ := status.FromError(err)
		return fiber.NewErrorResponse(
			fiberError.FiberError{
				Code:    int(responseStatus.Code()),
				Message: responseStatus.String(),
			})
	}

	return &Response{
		Metadata: responseHeader,
		Message:  response.Bytes(),
		Status:   *status.New(codes.OK, "Success"),
	}
}

// NewDispatcher is the constructor to create a dispatcher. It will create the clientconn and set defaults.
// Endpoint, serviceMethod and response proto are required minimally to work.
func NewDispatcher(config DispatcherConfig) (*Dispatcher, error) {
	configuredTimeout := TimeoutDefault
	if config.Timeout != 0 {
		configuredTimeout = config.Timeout
	}

	if config.Endpoint == "" || config.Service == "" || config.Method == "" {
		return nil, fiberError.ErrInvalidInput(
			protocol.GRPC,
			errors.New("grpc dispatcher: missing config (endpoint/service/method)"))
	}

	conn, err := grpc.DialContext(context.Background(), config.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// if ok is false, unknown codes.Unknown and Status msg is returned in Status
		responseStatus, _ := status.FromError(err)
		return nil, fiberError.ErrRequestFailed(
			protocol.GRPC,
			errors.New("grpc dispatcher: "+responseStatus.String()))
	}

	dispatcher := &Dispatcher{
		timeout:       configuredTimeout,
		serviceMethod: fmt.Sprintf("%s/%s", config.Service, config.Method),
		endpoint:      config.Endpoint,
		conn:          conn,
	}
	return dispatcher, nil
}
