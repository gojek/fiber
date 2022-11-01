package integration_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/config"
	fiberError "github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/grpc"
	fiberhttp "github.com/gojek/fiber/http"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	"github.com/gojek/fiber/internal/testutils"
	testGrpcUtils "github.com/gojek/fiber/internal/testutils/grpc"
	"github.com/gojek/fiber/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	httpResponse1 = []byte(`response 1`)
	httpResponse2 = []byte(`response 2`)
	httpResponse3 = []byte(`response 3`)
	httpAddr1     = ":5000"
	httpAddr2     = ":5001"
	httpAddr3     = ":5002"

	grpcPort1     = 50555
	grpcPort2     = 50556
	grpcPort3     = 50557
	grpcResponse1 = &testproto.PredictValuesResponse{
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
	grpcResponse2 = &testproto.PredictValuesResponse{}
	grpcResponse3 = &testproto.PredictValuesResponse{}
)

func TestMain(m *testing.M) {
	// Set up three grpc and http server with fix response for test.
	// Third routes will be set to timeout intentionally.

	runTestGrpcServer(grpcPort1, grpcResponse1, 0)
	runTestGrpcServer(grpcPort2, grpcResponse2, 0)
	runTestGrpcServer(grpcPort3, grpcResponse3, 10)

	runTestHttpServer(httpAddr1, httpResponse1, 0)
	runTestHttpServer(httpAddr2, httpResponse2, 0)
	runTestHttpServer(httpAddr3, httpResponse3, 10)

	os.Exit(m.Run())
}

func runTestHttpServer(addr string, responseBody []byte, delayDuration int) {
	// Create test server
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(time.Second * time.Duration(delayDuration))
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(responseBody)
		if err != nil {
			log.Fatal("set up: fail to write response body")
		}
	})

	go func() {
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatal("set up: start http server")
		}
	}()
}

func runTestGrpcServer(port int, response *testproto.PredictValuesResponse, delayDuration int) {
	testGrpcUtils.RunTestUPIServer(testGrpcUtils.GrpcTestServer{
		Port:         port,
		MockResponse: response,
		DelayTimer:   time.Second * time.Duration(delayDuration),
	})
}

func TestE2EFromConfig(t *testing.T) {
	bytePayload, _ := proto.Marshal(&testproto.PredictValuesRequest{
		PredictionRows: []*testproto.PredictionRow{
			{
				RowId: "1",
			},
			{
				RowId: "2",
			},
		},
	})
	grpcRequest := &grpc.Request{
		Message: bytePayload,
	}

	httpReq, err := http.NewRequest(
		http.MethodGet, "",
		ioutil.NopCloser(bytes.NewReader([]byte{})))
	require.NoError(t, err)
	httpRequest, err := fiberhttp.NewHTTPRequest(httpReq)
	require.NoError(t, err)

	route1 := "route1"
	route2 := "route2"
	route3 := "route3"

	tests := []struct {
		name                 string
		routesOrder          []string
		request              fiber.Request
		expectedMessageProto *testproto.PredictValuesResponse
		expectedFiberErr     fiber.Response
		expectedResponse     fiber.Response
		configPath           string
	}{
		{
			name:                 "grpc route 1",
			configPath:           "./fibergrpc.yaml",
			routesOrder:          []string{route1, route2, route3},
			request:              grpcRequest,
			expectedMessageProto: grpcResponse1,
			expectedResponse: (&grpc.Response{
				Status:   *status.New(codes.OK, "Success"),
				Metadata: metadata.New(map[string]string{}),
			}).WithLabel("order", []string{route1, route2, route3}...),
		},
		{
			name:        "http route 1",
			configPath:  "./fiberhttp.yaml",
			routesOrder: []string{route1, route2, route3},
			request:     httpRequest,
			expectedResponse: fiberhttp.NewHTTPResponse(
				&http.Response{
					StatusCode: http.StatusOK,
					Body:       makeBody(httpResponse1),
				},
			).WithLabel("order", []string{route1, route2, route3}...),
		},
		{
			name:                 "grpc route 2",
			configPath:           "./fibergrpc.yaml",
			routesOrder:          []string{route2, route1, route3},
			request:              grpcRequest,
			expectedMessageProto: grpcResponse2,
			expectedResponse: (&grpc.Response{
				Status:   *status.New(codes.OK, "Success"),
				Metadata: metadata.New(map[string]string{}),
			}).WithLabel("order", []string{route2, route1, route3}...),
		},
		{
			name:        "http route 2",
			configPath:  "./fiberhttp.yaml",
			routesOrder: []string{route2, route1, route3},
			request:     httpRequest,
			expectedResponse: fiberhttp.NewHTTPResponse(
				&http.Response{
					StatusCode: http.StatusOK,
					Body:       makeBody(httpResponse2),
				},
			).WithLabel("order", []string{route2, route1, route3}...),
		},
		{
			name:                 "grpc route3 timeout, route 1 fallback returned",
			configPath:           "./fibergrpc.yaml",
			routesOrder:          []string{route3, route1, route2},
			request:              grpcRequest,
			expectedMessageProto: grpcResponse1,
			expectedResponse: (&grpc.Response{
				Status:   *status.New(codes.OK, "Success"),
				Metadata: metadata.New(map[string]string{}),
			}).WithLabel("order", []string{route3, route1, route2}...),
		},
		{
			name:        "http route3 timeout, route 1 fallback returned",
			configPath:  "./fiberhttp.yaml",
			routesOrder: []string{route3, route1, route2},
			request:     httpRequest,
			expectedResponse: fiberhttp.NewHTTPResponse(
				&http.Response{
					StatusCode: http.StatusOK,
					Body:       makeBody(httpResponse1),
				},
			).WithLabel("order", []string{route3, route1, route2}...),
		},
		{
			name:                 "grpc route3 timeout, route 2 fallback returned",
			configPath:           "./fibergrpc.yaml",
			routesOrder:          []string{route3, route2, route1},
			request:              grpcRequest,
			expectedMessageProto: grpcResponse2,
			expectedResponse: (&grpc.Response{
				Status:   *status.New(codes.OK, "Success"),
				Metadata: metadata.New(map[string]string{}),
			}).WithLabel("order", []string{route3, route2, route1}...),
		},
		{
			name:        "http route3 timeout, route 2 fallback returned",
			configPath:  "./fiberhttp.yaml",
			routesOrder: []string{route3, route2, route1},
			request:     httpRequest,
			expectedResponse: fiberhttp.NewHTTPResponse(
				&http.Response{
					StatusCode: http.StatusOK,
					Body:       makeBody(httpResponse2),
				},
			).WithLabel("order", []string{route3, route2, route1}...),
		},
		{
			name:        "grpc route3 timeout",
			configPath:  "./fibergrpc.yaml",
			routesOrder: []string{route3},
			request:     grpcRequest,
			expectedResponse: &grpc.Response{
				Status: *status.New(codes.Unavailable, ""),
			},
			expectedFiberErr: fiber.
				NewErrorResponse(fiberError.ErrServiceUnavailable(protocol.GRPC)).
				WithLabel("order", []string{route3}...),
		},
		{
			name:        "http route3 timeout",
			configPath:  "./fiberhttp.yaml",
			routesOrder: []string{route3},
			request:     httpRequest,
			expectedFiberErr: fiber.
				NewErrorResponse(fiberError.ErrServiceUnavailable(protocol.HTTP)).
				WithLabel("order", []string{route3}...),
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable locally for parallelizing the tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			component, err := config.InitComponentFromConfig(tt.configPath)
			require.NoError(t, err)
			router, ok := component.(*fiber.EagerRouter)
			require.True(t, ok)

			// Orchestrate route order with mock strategy to fix the order of routes for testing
			strategy := testutils.NewMockRoutingStrategy(
				router.GetRoutes(),
				tt.routesOrder,
				0,
				nil,
			)
			router.SetStrategy(strategy)

			resp, ok := <-router.Dispatch(context.Background(), tt.request).Iter()
			require.True(t, ok)

			var finalResp fiber.Response
			if tt.expectedFiberErr != nil {
				finalResp = tt.expectedFiberErr
				assert.EqualValues(t, tt.expectedFiberErr, resp)
			} else {
				finalResp = tt.expectedResponse
				require.Equal(t, resp.StatusCode(), tt.expectedResponse.StatusCode())
				if tt.request.Protocol() == protocol.GRPC {
					responseProto := &testproto.PredictValuesResponse{}
					err = proto.Unmarshal(resp.Payload(), responseProto)
					require.NoError(t, err)

					assert.True(t, proto.Equal(tt.expectedMessageProto, responseProto), "actual proto response don't match expected")
				} else {
					assert.Equal(t, tt.expectedResponse.Payload(), resp.Payload())
				}
			}
			// Verify that the labels were correctly set, for both success and failure response
			assert.Equal(t, finalResp.Label("order"), resp.Label("order"))
		})
	}
}

func makeBody(body []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(body))
}
