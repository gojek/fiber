package testutils

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcTestServer struct {
	Port         int
	MockResponse *testproto.PredictValuesResponse
}

func (s *GrpcTestServer) PredictValues(_ context.Context, _ *testproto.PredictValuesRequest) (*testproto.PredictValuesResponse, error) {
	if s.MockResponse != nil {
		log.Println("response = " + s.MockResponse.String())
		return s.MockResponse, nil
	}
	return &testproto.PredictValuesResponse{
		Metadata: &testproto.ResponseMetadata{
			PredictionId: "123",
			ExperimentId: strconv.Itoa(s.Port),
		},
	}, nil
}

func RunTestUPIServer(srv GrpcTestServer) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.Port))
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	testproto.RegisterUniversalPredictionServiceServer(s, &srv)
	reflection.Register(s)
	log.Printf("Running Test Server at %v", srv.Port)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
