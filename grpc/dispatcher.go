package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/gojek/fiber"
	fiberError "github.com/gojek/fiber/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	TimeoutDefault = time.Second
)

type Dispatcher struct {
	timeout time.Duration
	// ServiceMethod is the service and method of server point in the format "{grpc_service_name}/{method_name}"
	serviceMethod string
	// Endpoint is the host+port of the grpc server, eg "127.0.0.1:50050"
	endpoint string

	conn *grpc.ClientConn
}

type DispatcherConfig struct {
	ServiceMethod string
	Endpoint      string
	DialOption    grpc.DialOption
	Timeout       time.Duration
	ResponseProto proto.Message
}

func (d *Dispatcher) Do(request fiber.Request) fiber.Response {

	grpcRequest, ok := request.(*Request)
	if !ok {
		return fiber.NewErrorResponse(
			fiberError.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "fiber: grpc.Dispatcher supports only grpc.Request type of requests",
			})
	}

	err := d.isValid()
	if err != nil {
		return fiber.NewErrorResponse(err)
	}

	//TODO add timeout to dial option
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, grpcRequest.Metadata)

	responseProto := proto.Clone(grpcRequest.ResponseProto)
	var responseHeader metadata.MD
	err = d.conn.Invoke(ctx, d.serviceMethod, request.Payload(), responseProto, grpc.Header(&responseHeader))
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
		Metadata:        responseHeader,
		ResponsePayload: responseProto,
		Status:          *status.New(codes.OK, "Success"),
	}
}

func (d *Dispatcher) isValid() error {

	if d.endpoint == "" || d.serviceMethod == "" {
		return fiberError.ErrInvalidInput(
			fiber.GRPC.String(),
			errors.New("missing endpoint/serviceMethod"))
	}

	if d.timeout <= 0 {
		return fiberError.ErrInvalidInput(
			fiber.GRPC.String(),
			errors.New("invalid or no timeout configured"))
	}

	if d.conn == nil {
		return fiberError.NewFiberError(
			fiber.GRPC.String(),
			errors.New("connection not created, use dispatcher constructor"),
		)
	}

	//if d.ResponseProto == nil {
	//	return fiberError.NewFiberError(
	//		fiber.GRPC.String(),
	//		errors.New("response proto not specified in dispatcher"),
	//	)
	//}
	return nil
}

// NewDispatcher is the constructor to create a dispatcher. It will create the clientconn and set defaults.
// Endpoint, serviceMethod and response proto are required minimally to work.
func NewDispatcher(config DispatcherConfig) (*Dispatcher, error) {

	configuredTimeout := TimeoutDefault
	if config.Timeout != 0 {
		configuredTimeout = config.Timeout
	}

	if config.Endpoint == "" || config.ServiceMethod == "" {
		return nil,
			fiberError.ErrInvalidInput(
				fiber.GRPC.String(),
				errors.New("missing endpoint/serviceMethod"))
	}

	//TODO pass in dialoption
	conn, err := grpc.Dial(config.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// if ok is false, unknown codes.Unknown and Status msg is returned in Status
		responseStatus, _ := status.FromError(err)
		return nil,
			fiberError.NewFiberError(
				fiber.GRPC.String(),
				errors.New(responseStatus.String()),
			)
	}

	return &Dispatcher{
		timeout:       configuredTimeout,
		serviceMethod: config.ServiceMethod,
		endpoint:      config.Endpoint,
		conn:          conn,
	}, nil
}
