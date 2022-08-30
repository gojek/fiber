package grpc

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gojek/fiber"
	fiberError "github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/http"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"github.com/gojek/fiber/protocol"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
)

const (
	port    = 50055
	service = "testproto.UniversalPredictionService"
	method  = "PredictValues"
)

var mockResponse *testproto.PredictValuesResponse

func TestMain(m *testing.M) {

	mockResponse = &testproto.PredictValuesResponse{
		Predictions: []*testproto.PredictionResult{
			{
				RowId: "1",
				Value: &testproto.NamedValue{
					Name:        "str",
					Type:        testproto.NamedValue_TYPE_STRING,
					StringValue: "213",
				},
			},
			{
				RowId: "2",
				Value: &testproto.NamedValue{
					Name:        "double",
					Type:        testproto.NamedValue_TYPE_DOUBLE,
					DoubleValue: 123.45,
				},
			},
			{
				RowId: "3",
				Value: &testproto.NamedValue{
					Name:         "int",
					Type:         testproto.NamedValue_TYPE_INTEGER,
					IntegerValue: 2,
				},
			},
		},
		Metadata: &testproto.ResponseMetadata{
			PredictionId: "abc",
			ModelName:    "linear",
			ModelVersion: "1.2",
			ExperimentId: "1",
			TreatmentId:  "2",
		},
	}

	//Test server will run upi server at port 50055
	testutils.RunTestUPIServer(
		testutils.GrpcTestServer{
			Port:         port,
			MockResponse: mockResponse,
		},
	)
	os.Exit(m.Run())
}

func TestNewDispatcher(t *testing.T) {
	tests := []struct {
		name              string
		dispatcherConfig  DispatcherConfig
		expected          *Dispatcher
		expectedProtoName string
		expectedErr       *fiberError.FiberError
	}{
		{
			name: "empty endpoint",
			dispatcherConfig: DispatcherConfig{
				Service: service,
				Method:  method,
			},
			expected: nil,
			expectedErr: fiberError.ErrInvalidInput(
				protocol.GRPC,
				errors.New("grpc dispatcher: missing config (endpoint/service/method/response-proto)")),
		},
		{
			name: "empty service",
			dispatcherConfig: DispatcherConfig{
				Method:   method,
				Endpoint: fmt.Sprintf(":%d", port),
			},
			expected: nil,
			expectedErr: fiberError.ErrInvalidInput(
				protocol.GRPC,
				errors.New("grpc dispatcher: missing config (endpoint/service/method/response-proto)")),
		},
		{
			name: "empty method",
			dispatcherConfig: DispatcherConfig{
				Service:  service,
				Endpoint: fmt.Sprintf(":%d", port),
			},
			expected: nil,
			expectedErr: fiberError.ErrInvalidInput(
				protocol.GRPC,
				errors.New("grpc dispatcher: missing config (endpoint/service/method/response-proto)")),
		},
		{
			name: "invalid endpoint",
			dispatcherConfig: DispatcherConfig{
				Service:  service,
				Method:   method,
				Endpoint: ":1",
			},
			expected: nil,
			expectedErr: fiberError.NewFiberError(
				protocol.GRPC,
				errors.New("grpc dispatcher: unable to get reflection information, ensure server reflection is enable and config are correct")),
		},
		{
			name: "invalid response",
			dispatcherConfig: DispatcherConfig{
				Service:  service,
				Method:   "fake method",
				Endpoint: fmt.Sprintf(":%d", port),
			},
			expected: nil,
			expectedErr: fiberError.NewFiberError(
				protocol.GRPC,
				errors.New("grpc dispatcher: unable to fetch file descriptors, ensure config are correct")),
		},
		{
			name: "ok response",
			dispatcherConfig: DispatcherConfig{
				Service:  service,
				Method:   method,
				Endpoint: fmt.Sprintf(":%d", port),
				Timeout:  time.Second * 5,
			},
			expected: &Dispatcher{
				timeout:       time.Second * 5,
				serviceMethod: fmt.Sprintf("%s/%s", service, method),
				endpoint:      fmt.Sprintf(":%d", port),
			},
			expectedProtoName: "PredictValuesResponse",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDispatcher(tt.dispatcherConfig)
			if tt.expectedErr != nil {
				fiberErr, ok := err.(*fiberError.FiberError)
				require.True(t, ok, "error not fiber error")
				require.Equal(t, tt.expectedErr, fiberErr)
			} else {
				require.NoError(t, err)
				// responseProto and conn are ignored as they have pointer which value will not be identical
				diff := cmp.Diff(tt.expected, got,
					cmpopts.IgnoreFields(Dispatcher{}, "responseProto", "conn"),
					cmp.AllowUnexported(Dispatcher{}),
				)
				require.Empty(t, diff)
				responseProto, ok := got.responseProto.(*dynamicpb.Message)
				require.True(t, ok, "fail to convert response proto")
				require.Equal(t, tt.expectedProtoName, string(responseProto.Type().Descriptor().Name()))
			}
		})
	}
}

func TestDispatcher_Do(t *testing.T) {
	dispatcherConfig := DispatcherConfig{
		Service:  service,
		Method:   method,
		Endpoint: fmt.Sprintf(":%d", port),
		Timeout:  time.Second * 5,
	}
	dispatcher, err := NewDispatcher(dispatcherConfig)
	require.NoError(t, err, "unable to create dispatcher")

	tests := []struct {
		name             string
		dispatcherConfig *DispatcherConfig
		responseProto    proto.Message
		input            fiber.Request
		expected         fiber.Response
	}{
		{
			name:  "non grpc request",
			input: &http.Request{},
			expected: fiber.NewErrorResponse(fiberError.FiberError{
				Code:    int(codes.InvalidArgument),
				Message: "fiber: grpc dispatcher: only grpc.Request type of requests are supported",
			}),
		},
		{
			name: "success",
			input: &Request{
				RequestPayload: &testproto.PredictValuesRequest{}},
			expected: &Response{
				Metadata: metadata.New(map[string]string{
					"content-type": "application/grpc",
				}),
				Status: *status.New(codes.OK, "Success"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := dispatcher.Do(tt.input)
			errResponse, ok := response.(*fiber.ErrorResponse)
			if ok {
				assert.EqualValues(t, tt.expected, errResponse)
			} else {
				grpcResponse, ok := response.(*Response)
				if !ok {
					assert.FailNow(t, "Fail to type assert response")
				}
				require.EqualValues(t, tt.expected.StatusCode(), grpcResponse.StatusCode())
				require.EqualValues(t, tt.expected.BackendName(), grpcResponse.BackendName())
				require.EqualValues(t, tt.expected.IsSuccess(), grpcResponse.IsSuccess())
				payload, ok := grpcResponse.Payload().([]byte)
				if !ok {
					assert.FailNow(t, "Fail to type assert response payload")
				}
				responseProto := &testproto.PredictValuesResponse{}
				err = proto.Unmarshal(payload, responseProto)
				require.NoError(t, err)
				assert.True(t, proto.Equal(mockResponse, responseProto), "actual proto response don't match expected")
			}
		})
	}
}
