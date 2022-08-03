package fiber_grpc

import (
	"context"
	"github.com/gojek/fiber"
	"github.com/gojek/fiber/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
)

type Dispatcher struct {
	methodName    string
	responseProto proto.Message
}

func (d *Dispatcher) Do(request fiber.Request) fiber.Response {

	grpcRequest, ok := request.(*Request)
	if !ok {
		log.Fatalf("non grpc request")
	}

	conn, err := grpc.Dial(grpcRequest.hostport, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error creating connection")
	}

	payload, ok := request.Payload().(proto.Message)
	if !ok {
		log.Fatalf("error reading payload")
	}

	//TODO need to pass header later
	responseProto := proto.Clone(d.responseProto)
	err = conn.Invoke(context.Background(), d.methodName, payload, responseProto)
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			log.Fatalf("error fetching status")
		}

		//TODO refactor errors.HTTPError into a generic error
		err = &errors.HTTPError{
			Code:    int(status.Code()),
			Message: status.Message(),
		}
		return fiber.NewErrorResponse(err)
	}

	return &Response{
		Metadata:        nil,
		ResponsePayload: responseProto,
	}

}

func NewDispatcher(methodName string, responseProto proto.Message) (fiber.Dispatcher, error) {

	//TODO add more validation here
	return &Dispatcher{
		methodName:    methodName,
		responseProto: responseProto,
	}, nil
}
