package grpc

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/http"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	port          = 50055
	serviceMethod = "testproto.UniversalPredictionService/PredictValues"
)

func TestDispatcher_Do(t *testing.T) {
	tests := []struct {
		name     string
		input    fiber.Request
		expected fiber.Response
	}{
		{
			name:  "non grpc request",
			input: &http.Request{},
			expected: fiber.NewErrorResponse(errors.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "fiber: grpc.Dispatcher supports only grpc.Request type of requests",
			}),
		},
		{
			name: "missing hostport",
			input: &Request{
				RequestPayload: &testproto.PredictValuesRequest{},
			},
			expected: fiber.NewErrorResponse(errors.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "missing hostport/servicemethod",
			}),
		},
		{
			name: "missing service method",
			input: &Request{
				RequestPayload: &testproto.PredictValuesRequest{},
				hostport:       fmt.Sprintf(":%d", port),
			},
			expected: fiber.NewErrorResponse(errors.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "missing hostport/servicemethod",
			}),
		},
		{
			name: "empty input",
			input: &Request{
				hostport:      fmt.Sprintf(":%d", port),
				ServiceMethod: serviceMethod,
			},
			expected: fiber.NewErrorResponse(errors.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "unable to convert payload to proto message",
			}),
		},
		{
			name: "invalid server address",
			input: &Request{
				RequestPayload: &testproto.PredictValuesRequest{},
				hostport:       "localhost:50050",
				ServiceMethod:  serviceMethod,
			},
			expected: fiber.NewErrorResponse(errors.FiberError{
				Code: int(codes.Unavailable),
				Message: "rpc error: code = Unavailable desc = connection error: desc = " +
					"\"transport: Error while dialing dial tcp [::1]:50050: " +
					"connect: connection refused\"",
			}),
		},
		{
			name: "success",
			input: &Request{
				RequestPayload: &testproto.PredictValuesRequest{},
				hostport:       fmt.Sprintf(":%d", port),
				ServiceMethod:  serviceMethod,
				ResponseProto:  &testproto.PredictValuesResponse{},
			},
			expected: &Response{
				Metadata: metadata.New(map[string]string{
					"content-type": "application/grpc",
				}),
				ResponsePayload: &testproto.PredictValuesResponse{
					Metadata: &testproto.ResponseMetadata{
						PredictionId: "123",
						ExperimentId: strconv.Itoa(50055),
					},
				},
				Status: *status.New(codes.OK, "Success"),
			},
		},
	}

	//Test server will run upi server at port 50055
	testutils.RunTestUPIServer(port)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{}
			response := d.Do(tt.input)

			errResponse, ok := response.(*fiber.ErrorResponse)
			if ok {
				log.Print(string(errResponse.Payload().([]byte)))
				assert.EqualValues(t, tt.expected, errResponse, string(errResponse.Payload().([]byte)))
			} else {
				grpcResponse, ok := response.(*Response)
				if !ok {
					assert.FailNow(t, "Fail to type assert response")
				}

				assert.EqualValues(t, tt.expected.StatusCode(), grpcResponse.StatusCode())
				assert.EqualValues(t, tt.expected.BackendName(), grpcResponse.BackendName())
				assert.EqualValues(t, tt.expected.IsSuccess(), grpcResponse.IsSuccess())
				expectedPayload, ok := tt.expected.Payload().(proto.Message)
				if !ok {
					assert.FailNow(t, "Fail to type assert response payload")
				}
				actualPayload, ok := grpcResponse.Payload().(proto.Message)
				if !ok {
					assert.FailNow(t, "Fail to type assert response")
				}
				assert.True(t, proto.Equal(expectedPayload, actualPayload), "actual payload doesn't equate expected")
			}
		})
	}
}
