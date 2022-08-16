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

type Dispatcher struct {
}

func (d *Dispatcher) Do(request fiber.Request) fiber.Response {

	grpcRequest, ok := request.(*Request)
	if !ok {
		return fiber.NewErrorResponse(errors.New("fiber: grpc.Dispatcher supports only grpc.Request type of requests"))
	}

	conn, err := grpc.Dial(grpcRequest.hostport, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// if ok is false, error is wrap inside status with unknown code
		responseStatus, _ := status.FromError(err)
		return &Response{status: responseStatus}
	}

	payload, ok := request.Payload().(proto.Message)
	if !ok {
		return fiber.NewErrorResponse(errors.New("unable to convert payload to proto message"))
	}

	//TODO add timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = metadata.NewOutgoingContext(ctx, grpcRequest.Metadata)

	responseProto := proto.Clone(grpcRequest.ResponseProto)
	var responseHeader metadata.MD
	err = conn.Invoke(ctx, grpcRequest.ServiceMethod, payload, responseProto, grpc.Header(&responseHeader))
	if err != nil {
		responseStatus, ok := status.FromError(err)
		if !ok {
			return &Response{status: responseStatus}
		}

		//TODO refactor errors.FiberError into a generic error
		err = &fiberError.FiberError{
			Code:    int(responseStatus.Code()),
			Message: responseStatus.Message(),
		}
		return fiber.NewErrorResponse(err)
	}

	return &Response{
		Metadata:        responseHeader,
		ResponsePayload: responseProto,
		status:          status.New(codes.OK, "Success"),
	}
}
