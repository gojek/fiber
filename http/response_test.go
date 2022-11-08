package http_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gojek/fiber"
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
  "error": "fiber: request cannot be completed: empty response received"
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
  "error": "fiber: request cannot be completed: unable to read response body: exception"
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

func TestHTTPResponseLabel(t *testing.T) {
	tests := map[string]struct {
		response fiber.Response
		key      string
		expected []string
	}{
		"empty labels": {
			response: fiberHTTP.NewHTTPResponse(&http.Response{
				Body:       makeBody([]byte("{}")),
				StatusCode: http.StatusOK,
			}),
			key: "dummy-key",
		},
		"case insensitive key": {
			response: fiberHTTP.NewHTTPResponse(&http.Response{
				Header:     http.Header{"Key": []string{"v1", "v2"}},
				Body:       makeBody([]byte("{}")),
				StatusCode: http.StatusOK,
			}),
			key:      "Key",
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

func TestHTTPResponseWithLabel(t *testing.T) {
	tests := map[string]struct {
		response fiber.Response
		key      string
		values   []string
		expected []string
	}{
		"new labels": {
			response: fiberHTTP.NewHTTPResponse(&http.Response{
				Body:       makeBody([]byte("{}")),
				StatusCode: http.StatusOK,
			}),
			key:      "k1",
			values:   []string{"v1", "v2"},
			expected: []string{"v1", "v2"},
		},
		"append labels": {
			response: fiberHTTP.NewHTTPResponse(&http.Response{
				Header:     http.Header{"K1": []string{"v1", "v2"}},
				Body:       makeBody([]byte("{}")),
				StatusCode: http.StatusOK,
			}),
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

func TestHTTPResponseWithLabels(t *testing.T) {
	tests := map[string]struct {
		response fiber.Response
		labels   fiber.Labels
		key      string
		expected []string
	}{
		"new labels": {
			response: fiberHTTP.NewHTTPResponse(&http.Response{
				Body:       makeBody([]byte("{}")),
				StatusCode: http.StatusOK,
			}),
			labels:   fiber.LabelsMap{"k1": []string{"v1", "v2"}},
			key:      "K1",
			expected: []string{"v1", "v2"},
		},
		"append labels": {
			response: fiberHTTP.NewHTTPResponse(&http.Response{
				Header:     http.Header{"K1": []string{"v1", "v2"}},
				Body:       makeBody([]byte("{}")),
				StatusCode: http.StatusOK,
			}),
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
