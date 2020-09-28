package http_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/gojek/fiber"
	fiberHTTP "github.com/gojek/fiber/http"
	"github.com/gojek/fiber/internal/testutils"
)

type unsupportedRequest struct {
	*fiber.CachedPayload
}

func (r *unsupportedRequest) Clone() (fiber.Request, error) {
	panic("not implemented")
}

func (r *unsupportedRequest) OperationName() string {
	panic("not implemented")
}

func (r *unsupportedRequest) Transform(_ fiber.Backend) (fiber.Request, error) {
	panic("not implemented")
}

func (r *unsupportedRequest) Header() map[string][]string {
	panic("not implemented")
}

type MockHTTPClient struct {
	mock.Mock
}

func (h *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := h.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

type dispatcherTestCase struct {
	name     string
	request  fiber.Request
	response *http.Response
	error    error
	expected fiber.Response
}

func (tt dispatcherTestCase) mockClient() *MockHTTPClient {
	mockClient := new(MockHTTPClient)
	if httpReq, ok := tt.request.(*fiberHTTP.Request); ok {
		mockClient.On("Do", httpReq.Request).Once().Return(tt.response, tt.error)
	}

	return mockClient
}

func TestNewDispatcher(t *testing.T) {
	mockClient := new(MockHTTPClient)

	dispatcher, err := fiberHTTP.NewDispatcher(mockClient)
	assert.NoError(t, err)
	assert.NotNil(t, dispatcher, "dispatcher should not be null")

	dispatcher, err = fiberHTTP.NewDispatcher(nil)

	assert.Errorf(t, err, "client can not be nil")
	assert.Nil(t, dispatcher)
}

func TestDispatcher_Do(t *testing.T) {
	suite := []dispatcherTestCase{
		{
			name:    "valid response",
			request: testutils.MockReq("POST", "localhost:8080/dispatcher", ""),
			response: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("OK response"))),
			},
			expected: testutils.MockResp(200, "OK response", nil, nil),
		},
		{
			name:     "invalid response",
			request:  testutils.MockReq("POST", "localhost:8080/dispatcher", ""),
			error:    errors.New("http: nil Request.URL"),
			expected: fiber.NewErrorResponse(errors.New("http: nil Request.URL")),
		},
		{
			name:    "unsupported request",
			request: &unsupportedRequest{},
			expected: fiber.NewErrorResponse(
				errors.New("fiber: http.Dispatcher supports only http.Request type of requests")),
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := tt.mockClient()

			dispatcher, _ := fiberHTTP.NewDispatcher(mockClient)

			resp := dispatcher.Do(tt.request)
			assert.Equal(t, tt.expected, resp)
			mockClient.AssertExpectations(t)
		})
	}

}
