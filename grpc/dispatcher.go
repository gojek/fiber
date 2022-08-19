package grpc

import (
	"context"
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
	Timeout time.Duration
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

	if grpcRequest.endpoint == "" || grpcRequest.ServiceMethod == "" {
		return fiber.NewErrorResponse(
			fiberError.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "missing endpoint/servicemethod",
			})
	}

	conn, err := grpc.Dial(grpcRequest.endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// if ok is false, unknown codes.Unknown and Status msg is returned in Status
		responseStatus, _ := status.FromError(err)
		return &Response{Status: *responseStatus}
	}

	payload, ok := request.Payload().(proto.Message)
	if !ok {
		return fiber.NewErrorResponse(
			fiberError.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "unable to convert payload to proto message",
			})
	}

	// Set default timeout
	if d.Timeout == 0 {
		d.Timeout = TimeoutDefault
	}
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, grpcRequest.Metadata)

	responseProto := proto.Clone(grpcRequest.ResponseProto)
	var responseHeader metadata.MD
	err = conn.Invoke(ctx, grpcRequest.ServiceMethod, payload, responseProto, grpc.Header(&responseHeader))
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
