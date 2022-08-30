package integration_test

import (
	"bytes"
	"context"
	"github.com/gojek/fiber/config"
	"github.com/gojek/fiber/grpc"
	fiberhttp "github.com/gojek/fiber/http"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSimpleHttpFromConfig(t *testing.T) {
	// Create test server
	responseBody := []byte(`test`)
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(responseBody)
		require.NoError(t, err, "fail to write to body")
	})
	go func() {
		err := http.ListenAndServe(":5000", handler)
		require.NoError(t, err, "fail to start http")
	}()

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
	require.Equal(t, respByte, responseBody)
}

func TestSimpleGrpcFromConfig(t *testing.T) {
	//Create test server and response
	mockResponse := &testproto.PredictValuesResponse{
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

	testutils.RunTestUPIServer(testutils.GrpcTestServer{
		Port:         50555,
		MockResponse: mockResponse,
	})

	// initialize root-level fiber component from the config
	component, err := config.InitComponentFromConfig("./fibergrpc.yaml")
	require.NoError(t, err)

	resp, ok := <-component.Dispatch(context.Background(), req).Iter()
	if ok {
		if resp.StatusCode() == int(codes.OK) {
			payload, ok := resp.Payload().(proto.Message)
			require.True(t, ok, "fail to convert response to proto")
			payloadByte, err := proto.Marshal(payload)
			require.NoError(t, err, "unable to unmarshal proto")
			responseProto := &testproto.PredictValuesResponse{}
			err = proto.Unmarshal(payloadByte, responseProto)
			require.NoError(t, err, "unable to marshal proto")
			assert.True(t, proto.Equal(mockResponse, responseProto), "actual proto response don't match expected")
		} else {
			assert.FailNow(t, "unexpected status")
		}
	} else {
		assert.FailNow(t, "fail to receive response queue")
	}
}
