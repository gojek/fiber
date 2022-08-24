package main

import (
	"context"
	"fmt"
	"github.com/gojek/fiber/config"
	"github.com/gojek/fiber/grpc"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"log"
)

const (
	port1 = 50555
	port2 = 50556
)

func main() {

	testutils.RunTestUPIServer(port1)
	testutils.RunTestUPIServer(port2)

	// initialize root-level fiber component from the config
	component, err := config.InitComponentFromConfig("./example/simplegrpcfromconfig/fiber.yaml")
	if err != nil {
		log.Fatalf("\nerror: %v\n", err)
	}

	var req = &grpc.Request{
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
		if resp.StatusCode() == 0 {
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
