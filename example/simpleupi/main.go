package main

import (
	"context"
	"fmt"
	"github.com/gojek/fiber"
	"github.com/gojek/fiber/extras"
	"github.com/gojek/fiber/fiber_grpc"
	upiv1 "github.com/gojek/fiber/gen/proto/go/upi/v1"
	"log"
)

const (
	hostport1     = "localhost:9000"
	hostport2     = "localhost:9001"
	serviceMethod = "upi.v1.UniversalPredictionService/PredictValues"
)

func main() {
	var req = &fiber_grpc.Request{
		RequestPayload: &upiv1.PredictValuesRequest{
			PredictionRows: []*upiv1.PredictionRow{
				{
					RowId: "1",
				},
				{
					RowId: "2",
				},
			},
		},
	}

	// initialize root-level component
	component := fiber.NewEagerRouter("eager-router")
	component.SetStrategy(new(extras.RandomRoutingStrategy))

	// TODO serviceMethod can be moved inside request actually, also possible to use a jhump dynamic message to replace the proto but for future versions
	upiDispatcher, _ := fiber_grpc.NewDispatcher(serviceMethod, &upiv1.PredictValuesResponse{})
	caller, _ := fiber.NewCaller("", upiDispatcher)

	proxy1 := fiber.NewProxy(
		fiber.NewBackend("route-a", hostport1),
		caller)
	proxy2 := fiber.NewProxy(
		fiber.NewBackend("route-b", hostport2),
		caller)

	component.SetRoutes(map[string]fiber.Component{
		"route-a": proxy1,
		"route-b": proxy2,
	})

	// working calls
	//callProxy(proxy1, req)
	//callProxy(proxy2, req)

	// This block returning -> "code": 503, error": "fiber: no responses received", need to debug further
	select {
	case resp, ok := <-component.Dispatch(context.Background(), req).Iter():
		if ok {
			if resp.StatusCode() == 0 {
				payload, ok := resp.Payload().(upiv1.PredictValuesResponse)
				if !ok {
					log.Fatalf("fail to convert response to proto")
				}
				log.Print(payload.String())
			} else {
				log.Fatalf(fmt.Sprintf("%s", resp.Payload()))
			}
		}
		log.Fatalf("fail to receive response queue")
	}

}

func callProxy(proxy *fiber.Proxy, req *fiber_grpc.Request) {
	select {
	case resp, ok := <-proxy.Dispatch(context.Background(), req).Iter():
		if ok {
			if resp.StatusCode() == 0 {
				payload, ok := resp.Payload().(*upiv1.PredictValuesResponse)
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
}
