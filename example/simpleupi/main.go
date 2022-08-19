package main

import (
	"context"
	"fmt"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"google.golang.org/grpc/codes"
	"log"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/extras"
	"github.com/gojek/fiber/grpc"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
)

const (
	port1         = 50555
	port2         = 50556
	endpoint1     = "localhost:50555"
	endpoint2     = "localhost:50556"
	serviceMethod = "testproto.UniversalPredictionService/PredictValues"
)

func main() {

	testutils.RunTestUPIServer(port1)
	testutils.RunTestUPIServer(port2)

	// initialize root-level component
	component := fiber.NewEagerRouter("eager-router")
	component.SetStrategy(new(extras.RandomRoutingStrategy))

	upiDispatcher := &grpc.Dispatcher{}
	caller1, _ := fiber.NewCaller("", upiDispatcher)
	caller2, _ := fiber.NewCaller("", upiDispatcher)

	proxy1 := fiber.NewProxy(
		fiber.NewBackend("route-a", endpoint1),
		caller1)
	proxy2 := fiber.NewProxy(
		fiber.NewBackend("route-b", endpoint2),
		caller2)

	component.SetRoutes(map[string]fiber.Component{
		"route-a": proxy1,
		"route-b": proxy2,
	})

	var req = &grpc.Request{
		ServiceMethod: serviceMethod,
		ResponseProto: &testproto.PredictValuesResponse{},
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
			payload, ok := resp.Payload().(*testproto.PredictValuesResponse)
			if !ok {
				log.Fatalf("fail to convert response to proto")
			}
			log.Print(payload.String())
		} else {
			log.Fatalf(fmt.Sprintf("%s", resp.Payload()))
		}
	} else {
		log.Fatalf("fail to receive response queue")
	}

}
