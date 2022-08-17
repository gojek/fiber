package grpc

import (
	"github.com/gojek/fiber"
	"github.com/gojek/fiber/errors"
	upiv1 "github.com/gojek/fiber/gen/proto/go/upi/v1"
	"github.com/gojek/fiber/http"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strconv"
	"testing"
)

const (
	hostport      = ":50055"
	serviceMethod = "upi.v1.UniversalPredictionService/PredictValues"
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
				RequestPayload: &upiv1.PredictValuesRequest{},
			},
			expected: fiber.NewErrorResponse(errors.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "missing hostport/servicemethod",
			}),
		},
		{
			name: "missing service method",
			input: &Request{
				RequestPayload: &upiv1.PredictValuesRequest{},
				hostport:       hostport,
			},
			expected: fiber.NewErrorResponse(errors.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "missing hostport/servicemethod",
			}),
		},
		{
			name: "empty input",
			input: &Request{
				hostport:      hostport,
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
				RequestPayload: &upiv1.PredictValuesRequest{},
				hostport:       "localhost:50050",
				ServiceMethod:  serviceMethod,
			},
			expected: fiber.NewErrorResponse(errors.FiberError{
				Code:    int(codes.Unavailable),
				Message: "connection error: desc = \"transport: Error while dialing dial tcp [::1]:50050: connect: connection refused\"",
			}),
		},
		{
			name: "success",
			input: &Request{
				RequestPayload: &upiv1.PredictValuesRequest{},
				hostport:       hostport,
				ServiceMethod:  serviceMethod,
				ResponseProto:  &upiv1.PredictValuesResponse{},
			},
			expected: &Response{
				Metadata: metadata.New(map[string]string{
					"content-type": "application/grpc",
				}),
				ResponsePayload: &upiv1.PredictValuesResponse{
					Metadata: &upiv1.ResponseMetadata{
						PredictionId: "123",
						ExperimentId: strconv.Itoa(50055),
					},
				},
				Status: *status.New(codes.OK, "Success"),
			},
		},
	}

	//Test server will run upi server at port 50055
	testutils.RunTestUPIServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{}
			response := d.Do(tt.input)

			errResponse, ok := response.(*fiber.ErrorResponse)
			if ok {
				assert.EqualValues(t, tt.expected, errResponse)
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
				assert.EqualValues(t, expectedPayload.String(), actualPayload.String())
			}
		})
	}
}
