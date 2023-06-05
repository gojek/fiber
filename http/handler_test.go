package http_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gojek/fiber"
	fiberHTTP "github.com/gojek/fiber/http"
	"github.com/gojek/fiber/internal/testutils"
	testUtilsHttp "github.com/gojek/fiber/internal/testutils/http"
	"github.com/stretchr/testify/assert"
)

type handlerTestCase struct {
	name      string
	request   *http.Request
	responses []testUtilsHttp.DelayedResponse
	expected  *http.Response
	timeout   time.Duration
}

func (tt *handlerTestCase) mockComponent() fiber.Component {
	return testutils.NewMockComponent("component", tt.responses...)
}

func TestNewHandler(t *testing.T) {
	component := testutils.NewMockComponent("test")
	handler := fiberHTTP.NewHandler(component, fiberHTTP.Options{})

	assert.NotNil(t, handler)
}

func makeBody(body []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(body))
}

type errorBody struct{}

func (b *errorBody) Read([]byte) (n int, err error) {
	return 0, errors.New("exception")
}

func readBytes(closer io.ReadCloser) []byte {
	defer closer.Close()
	data, _ := io.ReadAll(closer)
	return data
}

func TestHandler_ServeHTTP(t *testing.T) {
	suite := []handlerTestCase{
		{
			name: "ok scenario",
			request: newHTTPRequest(
				"POST",
				"localhost:8080/handler",
				io.NopCloser(bytes.NewBuffer([]byte("request body")))),
			responses: []testUtilsHttp.DelayedResponse{
				{
					Response: testUtilsHttp.MockResp(
						200,
						string(responsePayload),
						http.Header{
							"Request-Id": {
								fmt.Sprintf("fiber-%d", time.Now().Unix()),
							},
						},
						nil),
					Latency: 20 * time.Millisecond,
				},
			},

			expected: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Request-Id": {
						fmt.Sprintf("fiber-%d", time.Now().Unix()),
					},
				},
				Body: makeBody(responsePayload),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:      "error: no responses",
			request:   newHTTPRequest("POST", "localhost:8080/handler", http.NoBody),
			responses: []testUtilsHttp.DelayedResponse{},
			expected: &http.Response{
				StatusCode: http.StatusBadGateway,
				Header:     http.Header{},
				Body: makeBody([]byte(
					`{
  "code": 502,
  "error": "fiber: no valid responses received from routes"
}`)),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:    "error: timeout exceeded",
			request: newHTTPRequest("POST", "localhost:8080/handler", http.NoBody),
			responses: []testUtilsHttp.DelayedResponse{
				{
					Response: testUtilsHttp.MockResp(
						200,
						string(responsePayload),
						http.Header{
							"Request-Id": {
								fmt.Sprintf("fiber-%d", time.Now().Unix()),
							},
						},
						nil),
					Latency: 100 * time.Millisecond,
				},
			},
			expected: &http.Response{
				StatusCode: http.StatusRequestTimeout,
				Header:     http.Header{},
				Body: makeBody([]byte(
					`{
  "code": 408,
  "error": "fiber: failed to receive a response within configured timeout"
}`)),
			},
			timeout: 20 * time.Millisecond,
		},
		{
			name:    "error: fail to read request",
			request: newHTTPRequest("POST", "localhost:8080/handler", &errorBody{}),
			expected: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Header:     http.Header{},
				Body: makeBody([]byte(
					`{
  "code": 500,
  "error": "fiber: failed to read incoming request: exception"
}`)),
			},
			timeout: 20 * time.Millisecond,
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {
			component := tt.mockComponent()

			handler := fiberHTTP.NewHandler(component, fiberHTTP.Options{Timeout: tt.timeout})

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, tt.request)

			assert.Equal(t, tt.expected.StatusCode, recorder.Code)
			assert.Equal(t, string(readBytes(tt.expected.Body)), recorder.Body.String())
			assert.Equal(t, tt.expected.Header, recorder.Header())
		})
	}
}
