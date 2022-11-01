package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/errors"
)

// HeaderBackendName is the default backend name
var headerBackendName = "X-Fiber-Route-ID"

type Response struct {
	*fiber.CachedPayload
	response *http.Response
}

// IsSuccess returns the success state of the request, which is true if the status
func (r *Response) IsSuccess() bool {
	return isSuccessStatus(r.StatusCode())
}

// Label returns all the values associated with the given key, in the response header.
// If the key does not exist, an empty slice will be returned.
func (r *Response) Label(key string) []string {
	return r.Header().Values(key)
}

// WithLabel appends the given value(s) to the key, in the response header.
// If the key does not already exist, a new key will be created.
// The modified response is returned.
func (r *Response) WithLabel(key string, values ...string) fiber.Response {
	for _, value := range values {
		r.Header().Add(key, value)
	}
	return r
}

// WithLabels does the same thing as WithLabel but over a collection of key-values.
func (r *Response) WithLabels(labels fiber.Labels) fiber.Response {
	for _, key := range labels.Keys() {
		values := labels.Label(key)
		for _, value := range values {
			r.Header().Add(key, value)
		}
	}
	return r
}

// BackendName returns the backend used to make the request
func (r *Response) BackendName() string {
	return strings.Join(r.Label(headerBackendName), ",")
}

// WithBackendName sets the given backend name in the response header.
// The modified response is returned.
func (r *Response) WithBackendName(backEnd string) fiber.Response {
	r.Header().Set(headerBackendName, backEnd)
	return r
}

// StatusCode returns the response status code
func (r *Response) StatusCode() int {
	return r.response.StatusCode
}

// Header returns the response header
func (r *Response) Header() http.Header {
	if r.response.Header == nil {
		r.response.Header = make(http.Header)
	}
	return r.response.Header
}

// FromHTTP constructs a fiber http or error response from http response / error object
func NewHTTPResponse(httpResponse *http.Response) fiber.Response {
	if httpResponse == nil {
		return fiber.NewErrorResponse(fmt.Errorf("empty response received"))
	}
	// Read the response body
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return fiber.NewErrorResponse(fmt.Errorf("unable to read response body: %s", err.Error()))
	}
	// If StatusCode is not OK, make error response
	if !isSuccessStatus(httpResponse.StatusCode) {
		// Wrap into a Fiber HTTP Error
		err = &errors.FiberError{
			Code:    httpResponse.StatusCode,
			Message: string(body),
		}
		return fiber.NewErrorResponse(err)
	}
	// Return the success response
	return &Response{
		response:      httpResponse,
		CachedPayload: fiber.NewCachedPayload(body),
	}
}

func isSuccessStatus(code int) bool {
	return code/100 == 2
}
