package testutils

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"

	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
)

type grpcServer struct {
	Port int
}

func (s *grpcServer) PredictValues(ctx context.Context, request *testproto.PredictValuesRequest) (*testproto.PredictValuesResponse, error) {
	return &testproto.PredictValuesResponse{
		Metadata: &testproto.ResponseMetadata{
			PredictionId: "123",
			ExperimentId: strconv.Itoa(s.Port),
		},
	}, nil
}

func RunTestUPIServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	testproto.RegisterUniversalPredictionServiceServer(s, &grpcServer{port})
	log.Printf("Running Test Server at %v", port)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
