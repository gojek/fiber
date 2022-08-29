package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gojek/fiber"
	fiberError "github.com/gojek/fiber/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

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
	// ResponseProto is the proto return type of the service.
	responseProto proto.Message
}

type DispatcherConfig struct {
	Service           string
	Method            string
	Endpoint          string
	Timeout           time.Duration
	ResponseProtoName string
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

	responseProto := proto.Clone(d.responseProto)
	var responseHeader metadata.MD
	err := d.conn.Invoke(ctx, d.serviceMethod, grpcRequest.Payload(), responseProto, grpc.Header(&responseHeader))
	if err != nil {
		// if ok is false, unknown codes.Unknown and Status msg is returned in Status
		responseStatus, _ := status.FromError(err)
		return fiber.NewErrorResponse(
			fiberError.FiberError{
				Code:    int(responseStatus.Code()),
				Message: responseStatus.String(),
			})
	}

	responseByte, err := proto.Marshal(responseProto)
	if err != nil {
		return fiber.NewErrorResponse(err)
	}

	return &Response{
		Metadata:      responseHeader,
		Status:        *status.New(codes.OK, "Success"),
		CachedPayload: fiber.NewCachedPayload(responseByte),
	}
}

// NewDispatcher is the constructor to create a dispatcher. It will create the clientconn and set defaults.
// Endpoint, serviceMethod and response proto are required minimally to work.
func NewDispatcher(config DispatcherConfig) (*Dispatcher, error) {
	configuredTimeout := TimeoutDefault
	if config.Timeout != 0 {
		configuredTimeout = config.Timeout
	}

	if config.Endpoint == "" || config.Service == "" || config.Method == "" || config.ResponseProtoName == "" {
		return nil, fiberError.ErrInvalidInput(
			fiber.GRPC,
			errors.New("grpc dispatcher: missing config (endpoint/service/method/response-proto)"))
	}

	conn, err := grpc.DialContext(context.Background(), config.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// if ok is false, unknown codes.Unknown and Status msg is returned in Status
		responseStatus, _ := status.FromError(err)
		return nil, fiberError.ErrRequestFailed(
			fiber.GRPC,
			errors.New("grpc dispatcher: "+responseStatus.String()))
	}

	// create a reflection client and get FileDescriptorProtos
	reflectionClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	req := &grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: config.Service,
		},
	}
	reflectionInfoClient, err := reflectionClient.ServerReflectionInfo(context.Background())
	if err != nil {
		return nil, fiberError.NewFiberError(
			fiber.GRPC,
			errors.New("grpc dispatcher: unable to get reflection information, ensure server reflection is enable and config are correct"))
	}
	if err = reflectionInfoClient.Send(req); err != nil {
		return nil, fiberError.NewFiberError(fiber.GRPC, err)
	}
	reflectionResponse, err := reflectionInfoClient.Recv()
	if err != nil {
		return nil, fiberError.NewFiberError(fiber.GRPC, err)
	}

	// From the FileDescriptorProtos, find the proto matching config proto name
	var fileDescriptorProtos []*descriptorpb.FileDescriptorProto
	for _, fdBytes := range reflectionResponse.GetFileDescriptorResponse().FileDescriptorProto {

		fd := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(fdBytes, fd); err != nil {
			return nil, fiberError.NewFiberError(fiber.GRPC, err)
		}
		fileDescriptorProtos = append(fileDescriptorProtos, fd)
	}
	var fdProto *descriptorpb.FileDescriptorProto
	for _, fd := range fileDescriptorProtos {
		for _, md := range fd.MessageType {
			if strings.Contains(*md.Name, config.ResponseProtoName) {
				fdProto = fd
				break
			}
		}
	}

	// Create a FileDescriptor from FileDescriptorProto, and get MessageDescriptor to create a dynamic message
	// Note: It might be required to register new proto using protoregistry.Files.RegisterFile() at runtime
	fileDescriptor, err := protodesc.NewFile(fdProto, protoregistry.GlobalFiles)
	if err != nil {
		return nil, fiberError.NewFiberError(
			fiber.GRPC,
			errors.New("grpc dispatcher: unable to find proto in registry"))
	}
	messageDescriptor := fileDescriptor.Messages().ByName(protoreflect.Name(config.ResponseProtoName))

	dispatcher := &Dispatcher{
		timeout:       configuredTimeout,
		serviceMethod: fmt.Sprintf("%s/%s", config.Service, config.Method),
		endpoint:      config.Endpoint,
		conn:          conn,
		responseProto: dynamicpb.NewMessage(messageDescriptor),
	}
	return dispatcher, nil
}
