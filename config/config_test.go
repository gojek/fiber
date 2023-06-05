package config_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/config"
	fibergrpc "github.com/gojek/fiber/grpc"
	fiberhttp "github.com/gojek/fiber/http"
	testutils "github.com/gojek/fiber/internal/testutils/grpc"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/dynamicpb"
)

type durCfgTestSuite struct {
	input    string
	duration time.Duration
	success  bool
}

const port = 50555

func TestDurationUnmarshalJSON(t *testing.T) {
	tests := map[string]durCfgTestSuite{
		"valid_seconds": {
			input:    "2s",
			duration: time.Second * 2,
			success:  true,
		},
		"valid_minute": {
			input:    "1m",
			duration: time.Minute,
			success:  true,
		},
		"valid_quoted_time": {
			input:    "\"2s\"",
			duration: time.Second * 2,
			success:  true,
		},
		"invalid_input": {
			input:    "xyz",
			duration: 0,
			success:  false,
		},
	}

	// Run the tests
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			var d config.Duration
			// Unmarshal
			err := d.UnmarshalJSON([]byte(data.input))
			// Verify
			assert.Equal(t, data.duration, time.Duration(d))
			assert.Equal(t, data.success, err == nil)
		})
	}
}

func TestDurationMarshalJSON(t *testing.T) {
	duration := config.Duration(time.Second * 2)
	data, err := json.Marshal(duration)
	assert.Equal(t, `"2s"`, string(data))
	assert.NoError(t, err)
}

func TestFromConfig(t *testing.T) {
	timeout := 20 * time.Second
	backend := fiber.NewBackend("proxy_name", "localhost:1234")

	httpDispatcher, _ := fiberhttp.NewDispatcher(&http.Client{Timeout: timeout})
	httpCaller, _ := fiber.NewCaller("proxy_name", httpDispatcher)
	httpProxy := fiber.NewProxy(backend, httpCaller)
	testutils.RunTestUPIServer(testutils.GrpcTestServer{
		Port: port,
	})

	grpcDispatcher, _ := fibergrpc.NewDispatcher(
		fibergrpc.DispatcherConfig{
			ServiceMethod: "testproto.UniversalPredictionService/PredictValues",
			Endpoint:      fmt.Sprintf("localhost:%d", port),
			Timeout:       timeout,
		})
	grpcCaller, _ := fiber.NewCaller("proxy_name", grpcDispatcher)
	grpcProxy := fiber.NewProxy(nil, grpcCaller)

	tests := []struct {
		name              string
		configPath        string
		expectedComponent fiber.Component
		expectedErrMsg    string
	}{
		{
			name:              "http proxy",
			configPath:        "../internal/testdata/config/http_proxy.yaml",
			expectedComponent: httpProxy,
		},
		{
			name:              "grpc proxy",
			configPath:        "../internal/testdata/config/grpc_proxy.yaml",
			expectedComponent: grpcProxy,
		},
		{
			name:           "grpc proxy",
			configPath:     "../internal/testdata/config/invalid_grpc_proxy.yaml",
			expectedErrMsg: "fiber: grpc dispatcher: missing config (endpoint/serviceMethod)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.InitComponentFromConfig(tt.configPath)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err)
				assert.True(t,
					cmp.Equal(tt.expectedComponent, got,
						cmpopts.IgnoreUnexported(grpc.ClientConn{}, dynamicpb.Message{}),
						cmpopts.IgnoreInterfaces(struct{ grpc.ClientConnInterface }{}),
						cmp.AllowUnexported(
							fiber.BaseComponent{},
							fiber.Proxy{},
							fiber.Caller{},
							fibergrpc.Dispatcher{},
							fiberhttp.Dispatcher{}),
					),
					"config not equal to expected")
			} else {
				assert.Equal(t, tt.expectedErrMsg, err.Error())
			}
		})
	}
}
