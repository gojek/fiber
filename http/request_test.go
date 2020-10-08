package http_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gojek/fiber"
	fiberHTTP "github.com/gojek/fiber/http"
	"github.com/gojek/fiber/internal/testutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var requestPayload, _ = testutils.ReadFile("../internal/testdata/request_payload.json")

type requestTestCase struct {
	name          string
	request       *http.Request
	payload       []byte
	expectedError error
}

func newHTTPRequest(method, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	return req
}

func TestNewHTTPRequest(t *testing.T) {
	suite := []requestTestCase{
		{
			name: "ok scenario",
			request: newHTTPRequest(
				http.MethodPost,
				"http://localhost:9999/test",
				makeBody(requestPayload),
			),
			payload: requestPayload,
		},
		{
			name: "empty payload",
			request: newHTTPRequest(
				http.MethodGet,
				"http://localhost:9999/empty",
				makeBody([]byte{}),
			),
			payload: []byte{},
		},
		{
			name: "error body",
			request: newHTTPRequest(
				http.MethodGet,
				"http://localhost:9999/error",
				&errorBody{},
			),
			expectedError: errors.New("exception"),
		},
		{
			name: "payload that can only be read once",
			request: newHTTPRequest(
				http.MethodGet,
				"http://localhost:9999/stream",
				strings.NewReader("*** can only be read once ***"),
			),
			payload: []byte("*** can only be read once ***"),
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {
			req, err := fiberHTTP.NewHTTPRequest(tt.request)

			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)

				require.Equal(t, map[string][]string(tt.request.Header), req.Header())
				require.Equal(t, tt.request.URL, req.URL)

				actualBody, _ := req.GetBody()
				require.Equal(t, tt.payload, readBytes(actualBody))

				actualBody, _ = req.GetBody()
				require.Equal(t, tt.payload, readBytes(actualBody),
					"it should be possible to read payload more than once")
			}

		})
	}
}

func TestRequest_Clone(t *testing.T) {
	req, _ := fiberHTTP.NewHTTPRequest(newHTTPRequest(
		http.MethodPost,
		"http://localhost:9999/api/mock",
		strings.NewReader("*** request payload ***"),
	))
	req.Request.Header = http.Header{
		"fiber-key": []string{
			fmt.Sprintf("test-%d", time.Now().Unix()),
		},
	}

	compareRequests := func(t *testing.T, original *fiberHTTP.Request, clone fiber.Request) {
		clonedReq, ok := clone.(*fiberHTTP.Request)
		require.True(t, ok, "clone should have the same type as the original")

		require.NotEqual(t, original, clonedReq)
		require.Equal(t, original.Header(), clonedReq.Header())

		expectedBody, err := original.GetBody()
		require.NoError(t, err)
		require.Equal(t, readBytes(expectedBody), readBytes(clonedReq.Body))
		require.Equal(t, req.URL, clonedReq.URL)
	}

	clone, err := req.Clone()
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		compareRequests(t, req, clone)
	})

	cloneOfClone, err := clone.Clone()
	require.NoError(t, err)

	t.Run("success | clone of cloned", func(t *testing.T) {
		compareRequests(t, req, cloneOfClone)
	})
}

func TestRequest_OperationName(t *testing.T) {
	reqPath := "/internal/api"
	method := http.MethodPost
	req, _ := fiberHTTP.NewHTTPRequest(newHTTPRequest(
		method,
		fmt.Sprintf("http://localhost:9999%s", reqPath),
		makeBody(requestPayload),
	))

	require.Equal(t, fmt.Sprintf("%s %s", method, reqPath), req.OperationName())
}

type mockBackend struct {
	mock.Mock
}

func (b *mockBackend) URL(requestURI string) string {
	args := b.Called(requestURI)
	return args.String(0)
}

type requestTransformTestCase struct {
	proxyURL      string
	requestPath   string
	expectedError error
}

func (tt *requestTransformTestCase) proxyEndpoint() string {
	return fmt.Sprintf("%s%s", tt.proxyURL, tt.requestPath)
}

func (tt *requestTransformTestCase) backend() *mockBackend {
	backend := new(mockBackend)
	backend.On("URL", tt.requestPath).
		Once().
		Return(tt.proxyEndpoint())

	return backend
}

func TestRequest_Transform(t *testing.T) {
	suite := map[string]requestTransformTestCase{
		"ok: transform": {
			proxyURL:    "http://proxy:8080",
			requestPath: "/api/path",
		},
		"error: invalid url": {
			proxyURL:      "%invalid url%",
			requestPath:   "/api/path",
			expectedError: errors.New("parse \"%invalid url%/api/path\": invalid URL escape \"%in\""),
		},
	}

	for name, tt := range suite {
		t.Run(name, func(t *testing.T) {
			req, _ := fiberHTTP.NewHTTPRequest(newHTTPRequest(
				http.MethodPost,
				fmt.Sprintf("http://localhost:9999%s", tt.requestPath),
				makeBody(requestPayload),
			))

			backend := tt.backend()
			transformedReq, err := req.Transform(backend)
			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.proxyEndpoint(), transformedReq.(*fiberHTTP.Request).URL.String())
			}

			backend.AssertExpectations(t)
		})
	}
}
