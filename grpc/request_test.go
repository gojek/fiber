package grpc

import (
	"testing"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/protocol"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestRequest_Clone(t *testing.T) {
	tests := []struct {
		name string
		req  *Request
	}{
		{
			name: "empty",
			req:  &Request{},
		},
		{
			name: "simple",
			req: &Request{
				Metadata: metadata.New(map[string]string{"test": "1"}),
				Message:  []byte("Testing"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone, err := tt.req.Clone()
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
			expected: "grpc",
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
		name     string
		req      Request
		expected []byte
	}{
		{
			name: "ok payload",
			req: Request{
				Message: []byte("Testing"),
			},
			expected: []byte("Testing"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.req.Payload())
		})
	}
}

func TestRequest_Protocol(t *testing.T) {

	req := Request{}
	assert.Equal(t, protocol.GRPC, req.Protocol())
}

func TestRequest_Transform(t *testing.T) {

	tests := []struct {
		name     string
		req      Request
		backend  fiber.Backend
		expected Request
	}{
		{
			name:     "",
			req:      Request{},
			backend:  nil,
			expected: Request{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.req.Transform(tt.backend)
			assert.NoError(t, err)
			assert.EqualValues(t, tt.expected, *got.(*Request), "Transform(%v)", tt.backend)
		})
	}
}
