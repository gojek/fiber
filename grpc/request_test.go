package grpc

import (
	"fmt"
	"testing"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/internal/mocks"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func TestRequest_Clone(t *testing.T) {
	tests := []struct {
		name string
		req  *Request
	}{
		{
			name: "",
			req:  &Request{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone, err := tt.req.Clone()
			fmt.Printf("%p\n", tt.req)
			fmt.Printf("%p\n", clone)

			fmt.Printf("%p\n", tt.req.Metadata)
			fmt.Printf("%p\n", clone.(*Request).Metadata)

			assert.NoError(t, err)
			assert.Equal(t, tt.req, clone)
		})
	}
}

func TestRequest_Header(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		want map[string][]string
	}{
		{
			name: "no metadata",
			req:  Request{},
			want: nil,
		},
		{
			name: "ok metadata",
			req: Request{
				Metadata: metadata.New(map[string]string{"test": "123"}),
			},
			want: map[string][]string{"test": {"123"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.req.Header())
		})
	}
}

func TestRequest_OperationName(t *testing.T) {
	tests := []struct {
		name     string
		req      Request
		expected string
	}{
		{
			name:     "empty request",
			req:      Request{},
			expected: "",
		},
		{
			name: "ok ServiceMethod",
			req: Request{
				ServiceMethod: "service/InvokeMethod",
			},
			expected: "service/InvokeMethod",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.req.OperationName())
		})
	}
}

func TestRequest_Payload(t *testing.T) {

	tests := []struct {
		name string
		req  Request
		want interface{}
	}{
		{
			name: "ok payload",
			req: Request{
				RequestPayload: &testproto.PredictValuesResponse{},
			},
			want: &testproto.PredictValuesResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, proto.Equal(tt.want.(proto.Message), tt.req.Payload().(proto.Message)), "payload not equal to expected")
		})
	}
}

func TestRequest_Protocol(t *testing.T) {

	req := Request{}
	assert.Equal(t, fiber.GRPC, req.Protocol())
}

func TestRequest_Transform(t *testing.T) {

	hostport := "1000"
	mockBackend := new(mocks.Backend)
	mockBackend.On("URL", "").Return(hostport)
	tests := []struct {
		name        string
		req         Request
		backend     fiber.Backend
		expected    Request
		expectedErr string
	}{
		{
			name:        "",
			req:         Request{},
			backend:     nil,
			expected:    Request{},
			expectedErr: "backend cannot be nil",
		},
		{
			name:    "",
			req:     Request{},
			backend: mockBackend,
			expected: Request{
				endpoint: hostport,
			},
			expectedErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.req.Transform(tt.backend)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.expected, *got.(*Request), "Transform(%v)", tt.backend)
			} else {
				assert.Equal(t, tt.expectedErr, err.Error())
			}
		})
	}
}
