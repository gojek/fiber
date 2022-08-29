package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/extras"
	"github.com/gojek/fiber/grpc"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

const (
	port1             = 50555
	port2             = 50556
	endpoint1         = "localhost:50555"
	endpoint2         = "localhost:50556"
	service           = "testproto.UniversalPredictionService"
	method            = "PredictValues"
	responseProtoName = "PredictValuesResponse"
)

func main() {

	testutils.RunTestUPIServer(testutils.GrpcTestServer{
		Port: port1,
	})
	testutils.RunTestUPIServer(testutils.GrpcTestServer{
		Port: port2,
	})

	// initialize root-level component
	component := fiber.NewEagerRouter("eager-router")
	component.SetStrategy(new(extras.RandomRoutingStrategy))

	upiDispatcher1, _ := grpc.NewDispatcher(grpc.DispatcherConfig{
		Endpoint:          endpoint1,
		Service:           service,
		Method:            method,
		ResponseProtoName: responseProtoName,
	})
	upiDispatcher2, _ := grpc.NewDispatcher(grpc.DispatcherConfig{
		Endpoint:          endpoint2,
		Service:           service,
		Method:            method,
		ResponseProtoName: responseProtoName,
	})

	// Caller is required to work with combiner, fanout. Using a dispatcher plainly doesn't work
	caller1, _ := fiber.NewCaller("", upiDispatcher1)
	caller2, _ := fiber.NewCaller("", upiDispatcher2)

	// For grpc proxy, backend is not used to set endpoints unlike the http proxy
	proxy1 := fiber.NewProxy(nil, caller1)
	proxy2 := fiber.NewProxy(nil, caller2)

	// Set both routes to the router component
	component.SetRoutes(map[string]fiber.Component{
		"route-a": proxy1,
		"route-b": proxy2,
	})

	var req = &grpc.Request{
		RequestPayload: &testproto.PredictValuesRequest{
			PredictionRows: []*testproto.PredictionRow{
				{
					RowId: "1",
				},
				{
					RowId: "2",
				},
			},
		},
	}

	resp, ok := <-component.Dispatch(context.Background(), req).Iter()
	if ok {
		if resp.StatusCode() == int(codes.OK) {
			responseProto := &testproto.PredictValuesResponse{}
			err := proto.Unmarshal(resp.Payload().([]byte), responseProto)
			if err != nil {
				log.Fatalf("fail to unmarshal to proto")
			}

			if !ok {
				log.Fatalf("fail to convert response to proto")
			}
			log.Print(responseProto.String())
		} else {
			log.Fatalf(fmt.Sprintf("%s", resp.Payload()))
		}
	} else {
		log.Fatalf("fail to receive response queue")
	}
}
