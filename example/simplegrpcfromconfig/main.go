package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gojek/fiber/config"
	"github.com/gojek/fiber/grpc"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
)

const (
	port1 = 50555
	port2 = 50556
)

func main() {

	testutils.RunTestUPIServer(testutils.GrpcTestServer{
		Port: port1,
	})
	testutils.RunTestUPIServer(testutils.GrpcTestServer{
		Port: port2,
	})

	// initialize root-level fiber component from the config
	component, err := config.InitComponentFromConfig("./example/simplegrpcfromconfig/fiber.yaml")
	if err != nil {
		log.Fatalf("\nerror: %v\n", err)
	}

	var req = &grpc.Request{
		Message: &testproto.PredictValuesRequest{
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
			log.Print(resp.Payload())

			//values can be retrieved using protoReflect or marshalled into proto
			payload, ok := resp.Payload().(proto.Message)
			if !ok {
				log.Fatalf("fail to convert response to proto")
			}
			payloadByte, err := proto.Marshal(payload)
			if err != nil {
				log.Fatalf("fail to marshal to proto")
			}
			responseProto := &testproto.PredictValuesResponse{}
			err = proto.Unmarshal(payloadByte, responseProto)
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
