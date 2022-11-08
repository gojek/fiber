package grpc_test

import (
	"log"
	"testing"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/grpc"
	testproto "github.com/gojek/fiber/internal/testdata/gen/testdata/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestResponse_Backend(t *testing.T) {
	tests := []struct {
		name        string
		res         grpc.Response
		want        grpc.Response
		backendName string
	}{
		{
			name: "ok",
			res: grpc.Response{
				Metadata: map[string][]string{},
			},
			want: grpc.Response{
				Metadata: metadata.New(map[string]string{"backend": "testing"}),
			},
			backendName: "testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.res.WithBackendName(tt.backendName)
			log.Print(tt.res)
			assert.Equalf(t, tt.want, tt.res, "BackendName()")
		})
	}
}

func TestResponse_Status(t *testing.T) {
	tests := []struct {
		name            string
		res             grpc.Response
		expectedCode    int
		expectedSuccess bool
	}{
		{
			name: "ok",
			res: grpc.Response{
				Status: *status.New(codes.OK, ""),
			},
			expectedCode:    0,
			expectedSuccess: true,
		},
		{
			name: "ok",
			res: grpc.Response{
				Status: *status.New(codes.InvalidArgument, ""),
			},
			expectedCode:    3,
			expectedSuccess: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedCode, tt.res.StatusCode())
			assert.Equal(t, tt.expectedSuccess, tt.res.IsSuccess())
		})
	}
}

func TestResponse_Payload(t *testing.T) {
	response := &testproto.PredictValuesRequest{
		PredictionRows: []*testproto.PredictionRow{
			{
				RowId: "123",
			},
		},
		Metadata: nil,
	}
	responseByte, _ := proto.Marshal(response)
	tests := []struct {
		name     string
		req      grpc.Response
		expected []byte
	}{
		{
			name: "",
			req: grpc.Response{
				Message: responseByte,
			},
			expected: responseByte,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.req.Payload())
		})
	}
}

func TestResponse_Label(t *testing.T) {
	tests := map[string]struct {
		response fiber.Response
		key      string
		expected []string
	}{
		"empty labels": {
			response: &grpc.Response{
				Metadata: map[string][]string{},
			},
			key: "dummy-key",
		},
		"non-empty labels": {
			response: &grpc.Response{
				Metadata: map[string][]string{"key": []string{"v1", "v2"}},
			},
			key:      "key",
			expected: []string{"v1", "v2"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			values := tt.response.Label(tt.key)
			assert.Equal(t, tt.expected, values)
		})
	}
}

func TestResponse_WithLabel(t *testing.T) {
	tests := map[string]struct {
		response fiber.Response
		key      string
		values   []string
		expected []string
	}{
		"new labels": {
			response: &grpc.Response{
				Metadata: map[string][]string{"key": []string{"v1", "v2"}},
			},
			key:      "k1",
			values:   []string{"v1", "v2"},
			expected: []string{"v1", "v2"},
		},
		"append labels": {
			response: &grpc.Response{
				Metadata: map[string][]string{"k1": []string{"v1", "v2"}},
			},
			key:      "k1",
			values:   []string{"v3"},
			expected: []string{"v1", "v2", "v3"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			newLabels := tt.response.WithLabel(tt.key, tt.values...)
			assert.Equal(t, tt.expected, newLabels.Label(tt.key))
		})
	}
}

func TestHTTPResponse_WithLabels(t *testing.T) {
	tests := map[string]struct {
		response fiber.Response
		labels   fiber.Labels
		key      string
		expected []string
	}{
		"new labels": {
			response: &grpc.Response{
				Metadata: map[string][]string{"key": []string{"v1", "v2"}},
			},
			labels:   fiber.LabelsMap{"k1": []string{"v1", "v2"}},
			key:      "k1",
			expected: []string{"v1", "v2"},
		},
		"append labels": {
			response: &grpc.Response{
				Metadata: map[string][]string{"k1": []string{"v1", "v2"}},
			},
			labels:   fiber.LabelsMap{"k1": []string{"v3"}},
			key:      "k1",
			expected: []string{"v1", "v2", "v3"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			newLabels := tt.response.WithLabels(tt.labels)
			assert.Equal(t, tt.expected, newLabels.Label(tt.key))
		})
	}
}
