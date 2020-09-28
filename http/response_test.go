package http_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	fiberHTTP "github.com/gojek/fiber/http"
	"github.com/gojek/fiber/internal/testutils"
)

var responsePayload, _ = testutils.ReadFile("../internal/testdata/response_payload.json")

type responseTestCase struct {
	name     string
	response *http.Response
	expected struct {
		payload []byte
		status  int
	}
}

func TestNewHTTPResponse(t *testing.T) {
	suite := []responseTestCase{
		{
			name: "ok: response",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       makeBody(responsePayload),
			},
			expected: struct {
				payload []byte
				status  int
			}{payload: responsePayload, status: http.StatusOK},
		},
		{
			name:     "error: nil http response",
			response: nil,
			expected: struct {
				payload []byte
				status  int
			}{status: http.StatusInternalServerError, payload: []byte(`{
  "code": 500,
  "error": "fiber: empty response received"
}`),
			},
		},
		{
			name: "error: failure response",
			response: &http.Response{
				StatusCode: http.StatusForbidden,
				Body:       makeBody([]byte("access denied")),
			},
			expected: struct {
				payload []byte
				status  int
			}{status: http.StatusForbidden, payload: []byte(`{
  "code": 403,
  "error": "access denied"
}`),
			},
		},
		{
			name: "error: closed body",
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(&errorBody{}),
			},
			expected: struct {
				payload []byte
				status  int
			}{status: http.StatusInternalServerError, payload: []byte(`{
  "code": 500,
  "error": "exception"
}`),
			},
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {

			resp := fiberHTTP.NewHTTPResponse(tt.response)

			require.NotNil(t, resp)
			require.Equal(t, string(tt.expected.payload), string(resp.Payload()))
			require.Equal(t, tt.expected.status, resp.StatusCode())
			require.Equal(t, tt.expected.status/100 == 2, resp.IsSuccess())
		})
	}
}
