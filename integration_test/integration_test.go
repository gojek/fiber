package integration_test

import (
	"bytes"
	"context"
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
	"google.golang.org/protobuf/proto"
)

var (
	httpResponse1 = []byte(`response 1`)
	httpResponse2 = []byte(`response 2`)
	httpAddr1     = ":5000"
	httpAddr2     = ":5001"

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
	// Set up three http and grpc server with fix response for test
	runTestHttpServer(httpAddr1, httpResponse1)
	runTestHttpServer(httpAddr2, httpResponse2)

	// Third routes will be set to timeout intentionally
	runTestGrpcServer(grpcPort1, grpcResponse1, 0)
	runTestGrpcServer(grpcPort2, grpcResponse2, 0)
	runTestGrpcServer(grpcPort3, grpcResponse3, 10)

	os.Exit(m.Run())
}

func runTestHttpServer(addr string, responseBody []byte) {
	// Create test server
	//responseBody1 := []byte(`response 1`)
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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

func TestSimpleHttpFromConfig(t *testing.T) {

	//Generate http request
	httpReq, err := http.NewRequest(
		http.MethodGet, "",
		ioutil.NopCloser(bytes.NewReader([]byte{})))
	require.NoError(t, err)

	req, err := fiberhttp.NewHTTPRequest(httpReq)
	require.NoError(t, err)

	// initialize root-level fiber component from the config
	component, err := config.InitComponentFromConfig("./fiberhttp.yaml")
	require.NoError(t, err)

	resp, ok := <-component.Dispatch(context.Background(), req).Iter()
	require.True(t, ok)
	require.Equal(t, resp.StatusCode(), http.StatusOK)
	respByte, ok := resp.Payload().([]byte)
	require.True(t, ok)
	require.Equal(t, respByte, httpResponse1)
}

func TestSimpleGrpcFromConfig(t *testing.T) {
	req := &grpc.Request{
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

	//Set up the router. route 1 and 2 are working fine, route 3 will always timeout.
	component, err := config.InitComponentFromConfig("./fibergrpc.yaml")
	require.NoError(t, err)
	router, ok := component.(*fiber.EagerRouter)
	require.True(t, ok)
	route1 := "route_a"
	route2 := "route_b"
	route3 := "route_c"

	tests := []struct {
		name                    string
		routesOrder             []string
		expectedResponseMessage *testproto.PredictValuesResponse
		expectedFiberErr        fiber.Response
		expectedStatus          int
	}{
		{
			name:                    "route 1",
			routesOrder:             []string{route1, route2, route3},
			expectedStatus:          int(codes.OK),
			expectedResponseMessage: grpcResponse1,
		},
		{
			name:                    "route 2",
			routesOrder:             []string{route2, route1, route3},
			expectedStatus:          int(codes.OK),
			expectedResponseMessage: grpcResponse2,
		},
		{
			name:                    "route3 timeout, route 1 fallback returned",
			routesOrder:             []string{route3, route1, route2},
			expectedStatus:          int(codes.OK),
			expectedResponseMessage: grpcResponse1,
		},
		{
			name:                    "route3 timeout, route 2 fallback returned",
			routesOrder:             []string{route3, route2, route1},
			expectedStatus:          int(codes.OK),
			expectedResponseMessage: grpcResponse2,
		},
		{
			name:             "route3timeout",
			routesOrder:      []string{route3},
			expectedStatus:   int(codes.Unavailable),
			expectedFiberErr: fiber.NewErrorResponse(fiberError.ErrServiceUnavailable(protocol.GRPC)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Orchestrate route order with mock strategy to fix the order of routes for testing
			strategy := testutils.NewMockRoutingStrategy(
				router.GetRoutes(),
				tt.routesOrder,
				0,
				nil,
			)
			router.SetStrategy(strategy)

			resp, ok := <-component.Dispatch(context.Background(), req).Iter()
			if ok {
				if resp.StatusCode() == tt.expectedStatus {
					if tt.expectedFiberErr != nil {
						payload, ok := resp.(*fiber.ErrorResponse)
						require.True(t, ok, "fail to convert response to err")
						assert.EqualValues(t, tt.expectedFiberErr, payload)
					} else {
						payload, ok := resp.Payload().(proto.Message)
						require.True(t, ok, "fail to convert response to proto")
						payloadByte, err := proto.Marshal(payload)
						require.NoError(t, err, "unable to marshal proto")
						responseProto := &testproto.PredictValuesResponse{}
						err = proto.Unmarshal(payloadByte, responseProto)
						require.NoError(t, err, "unable to unmarshal proto")
						assert.True(t, proto.Equal(tt.expectedResponseMessage, responseProto), "actual proto response don't match expected")
					}
				} else {
					assert.FailNow(t, "unexpected status")
				}
			} else {
				assert.FailNow(t, "fail to receive response queue")
			}
		})
	}
}
