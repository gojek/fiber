package grpc

import (
	"log"
	"testing"

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
		res         Response
		want        Response
		backendName string
	}{
		{
			name: "ok",
			res: Response{
				Metadata: map[string][]string{},
			},
			want: Response{
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
		res             Response
		expectedCode    int
		expectedSuccess bool
	}{
		{
			name: "ok",
			res: Response{
				Status: *status.New(codes.OK, ""),
			},
			expectedCode:    0,
			expectedSuccess: true,
		},
		{
			name: "ok",
			res: Response{
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
	tests := []struct {
		name     string
		req      Response
		expected interface{}
	}{
		{
			name: "",
			req: Response{
				ResponsePayload: response,
			},
			expected: proto.Clone(response),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t,
				proto.Equal(
					tt.expected.(proto.Message),
					tt.req.Payload().(proto.Message),
				), "payload not equal to expected")
		})
	}
}
