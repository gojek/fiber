package testutils

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"

	upiv1 "github.com/gojek/fiber/gen/proto/go/upi/v1"
)

const (
	port = 50055
)

type grpcServer struct{}

func (s *grpcServer) PredictValues(ctx context.Context, request *upiv1.PredictValuesRequest) (*upiv1.PredictValuesResponse, error) {
	return &upiv1.PredictValuesResponse{
		Metadata: &upiv1.ResponseMetadata{
			PredictionId: "123",
			ExperimentId: strconv.Itoa(port),
		},
	}, nil
}

func RunTestUPIServer() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("%v", err)
	}
	s := grpc.NewServer()
	upiv1.RegisterUniversalPredictionServiceServer(s, &grpcServer{})
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
