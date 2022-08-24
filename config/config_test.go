package config_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/config"
	fibergrpc "github.com/gojek/fiber/grpc"
	fiberhttp "github.com/gojek/fiber/http"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type durCfgTestSuite struct {
	input    string
	duration time.Duration
	success  bool
}

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

	grpcDispatcher, _ := fibergrpc.NewDispatcher(
		fibergrpc.DispatcherConfig{
			ServiceMethod: "myService",
			Endpoint:      "localhost:1234",
			Timeout:       timeout,
		})
	grpcCaller, _ := fiber.NewCaller("proxy_name", grpcDispatcher)
	grpcProxy := fiber.NewProxy(backend, grpcCaller)

	tests := []struct {
		name       string
		configPath string
		want       fiber.Component
	}{
		{
			name:       "http proxy",
			configPath: "../internal/testdata/config/http_proxy.yaml",
			want:       httpProxy,
		},
		{
			name:       "grpc proxy",
			configPath: "../internal/testdata/config/grpc_proxy.yaml",
			want:       grpcProxy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.FromConfig(tt.configPath)
			assert.NoError(t, err)
			assert.True(t,
				cmp.Equal(tt.want, got,
					cmpopts.IgnoreUnexported(grpc.ClientConn{}),
					cmp.AllowUnexported(
						fiber.BaseComponent{},
						fiber.Proxy{},
						fiber.Caller{},
						fibergrpc.Dispatcher{},
						fiberhttp.Dispatcher{}),
				),
				"config not equal to expected")
		})
	}
}
