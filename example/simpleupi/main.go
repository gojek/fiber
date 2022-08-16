package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/extras"
	upiv1 "github.com/gojek/fiber/gen/proto/go/upi/v1"
	"github.com/gojek/fiber/grpc"
)

const (
	hostport1     = "localhost:9000"
	hostport2     = "localhost:9001"
	serviceMethod = "upi.v1.UniversalPredictionService/PredictValues"
)

func main() {
	var req = &grpc.Request{
		ServiceMethod: serviceMethod,
		ResponseProto: &upiv1.PredictValuesResponse{},
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

	upiDispatcher := &grpc.Dispatcher{}
	caller1, _ := fiber.NewCaller("", upiDispatcher)
	caller2, _ := fiber.NewCaller("", upiDispatcher)

	proxy1 := fiber.NewProxy(
		fiber.NewBackend("route-a", hostport1),
		caller1)
	proxy2 := fiber.NewProxy(
		fiber.NewBackend("route-b", hostport2),
		caller2)

	component.SetRoutes(map[string]fiber.Component{
		"route-a": proxy1,
		"route-b": proxy2,
	})

	// dispatch proxy directly
	//callProxy(proxy1, req)
	//callProxy(proxy2, req)

	select {
	case resp, ok := <-component.Dispatch(context.Background(), req).Iter():
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

func callProxy(proxy *fiber.Proxy, req *grpc.Request) *upiv1.PredictValuesResponse {
	select {
	case resp, ok := <-proxy.Dispatch(context.Background(), req).Iter():
		if ok {
			if resp.StatusCode() == 0 {
				payload, ok := resp.Payload().(*upiv1.PredictValuesResponse)
				if !ok {
					log.Fatalf("fail to convert response to proto")
				}
				log.Print(payload.String())
				return payload
			} else {
				log.Fatalf(fmt.Sprintf("%s", resp.Payload()))
			}
		} else {
			log.Fatalf("fail to receive response queue")
		}
	}

	return nil
}
