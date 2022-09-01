package grpc

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gojek/fiber"
	fiberError "github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/protocol"
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

	return &Response{
		Metadata: responseHeader,
		Message:  responseProto,
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
			errors.New("grpc dispatcher: missing config (endpoint/service/method/response-proto)"))
	}

	conn, err := grpc.DialContext(context.Background(), config.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// if ok is false, unknown codes.Unknown and Status msg is returned in Status
		responseStatus, _ := status.FromError(err)
		return nil, fiberError.ErrRequestFailed(
			protocol.GRPC,
			errors.New("grpc dispatcher: "+responseStatus.String()))
	}

	// Get reflection response from reflection server, which contain FileDescriptorProtos
	reflectionResponse, err := getReflectionResponse(conn, config.Service)
	if err != nil {
		return nil, err
	}
	fileDescriptorProtoBytes := reflectionResponse.GetFileDescriptorResponse().GetFileDescriptorProto()

	messageDescriptor, err := getMessageDescriptor(fileDescriptorProtoBytes, config.Service, config.Method)
	if err != nil {
		return nil, err
	}

	dispatcher := &Dispatcher{
		timeout:       configuredTimeout,
		serviceMethod: fmt.Sprintf("%s/%s", config.Service, config.Method),
		endpoint:      config.Endpoint,
		conn:          conn,
		responseProto: dynamicpb.NewMessage(messageDescriptor),
	}
	return dispatcher, nil
}

func getReflectionResponse(conn *grpc.ClientConn, serviceName string) (*grpc_reflection_v1alpha.ServerReflectionResponse, error) {
	// create a reflection client and get FileDescriptorProtos
	reflectionClient := grpc_reflection_v1alpha.NewServerReflectionClient(conn)
	req := &grpc_reflection_v1alpha.ServerReflectionRequest{
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: serviceName,
		},
	}
	reflectionInfoClient, err := reflectionClient.ServerReflectionInfo(context.Background())
	if err != nil {
		return nil, fiberError.NewFiberError(
			protocol.GRPC,
			errors.New("grpc dispatcher: unable to get reflection information, ensure server reflection is enable and config are correct"))
	}
	if err = reflectionInfoClient.Send(req); err != nil {
		return nil, fiberError.NewFiberError(protocol.GRPC, err)
	}
	reflectionResponse, err := reflectionInfoClient.Recv()
	if err != nil {
		return nil, fiberError.NewFiberError(protocol.GRPC, err)
	}

	return reflectionResponse, nil
}

func getMessageDescriptor(fileDescriptorProtoBytes [][]byte, serviceName string, methodName string) (protoreflect.MessageDescriptor, error) {
	fileDescriptorProto, outputProtoName, err := getFileDescriptorProto(fileDescriptorProtoBytes, serviceName, methodName)
	if err != nil {
		return nil, err
	}

	messageDescriptor, err := getMessageDescriptorByName(fileDescriptorProto, outputProtoName)
	if err != nil {
		return nil, err
	}
	return messageDescriptor, nil
}

func getFileDescriptorProto(fileDescriptorProtoBytes [][]byte, serviceName string, methodName string) (*descriptorpb.FileDescriptorProto, string, error) {
	var fileDescriptorProto *descriptorpb.FileDescriptorProto
	var outputProtoName string

	for _, fdpByte := range fileDescriptorProtoBytes {
		fdp := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(fdpByte, fdp); err != nil {
			return nil, "", fiberError.NewFiberError(protocol.GRPC, err)
		}

		for _, service := range fdp.Service {
			// find matching service descriptors from file descriptor
			if serviceName == fmt.Sprintf("%s.%s", fdp.GetPackage(), service.GetName()) {
				// find matching method from service descriptor
				for _, method := range service.Method {
					if method.GetName() == methodName {
						outputType := method.GetOutputType()
						//Get the proto name without package
						outputProtoName = outputType[strings.LastIndex(outputType, ".")+1:]
						fileDescriptorProto = fdp
						break
					}
				}
			}
			if fileDescriptorProto != nil {
				break
			}
		}
		if fileDescriptorProto != nil {
			break
		}
	}

	if fileDescriptorProto == nil {
		return nil, "", fiberError.NewFiberError(
			protocol.GRPC,
			errors.New("grpc dispatcher: unable to fetch file descriptors, ensure config are correct"))
	}
	return fileDescriptorProto, outputProtoName, nil
}

func getMessageDescriptorByName(fileDescriptorProto *descriptorpb.FileDescriptorProto, outputProtoName string) (protoreflect.MessageDescriptor, error) {
	// Create a FileDescriptor from FileDescriptorProto, and get MessageDescriptor to create a dynamic message
	// Note: It might be required to register new proto using protoregistry.Files.RegisterFile() at runtime
	fileDescriptor, err := protodesc.NewFile(fileDescriptorProto, protoregistry.GlobalFiles)
	if err != nil {
		return nil, fiberError.NewFiberError(
			protocol.GRPC,
			errors.New("grpc dispatcher: unable to find proto in registry"))
	}
	messageDescriptor := fileDescriptor.Messages().ByName(protoreflect.Name(outputProtoName))
	return messageDescriptor, nil
}
