package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

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

func (r *Response) WithBackendName(backEnd string) fiber.Response {
	r.Header().Set(headerBackendName, backEnd)
	return r
}

// BackendName returns the backend used to make the request
func (r *Response) BackendName() string {
	if r.Header() == nil {
		r.response.Header = make(http.Header)
	}
	return r.Header().Get(headerBackendName)
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
		return fiber.NewErrorResponse(fmt.Errorf("fiber: empty response received"))
	}
	// Read the response body
	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return fiber.NewErrorResponse(err)
	}
	// If StatusCode is not OK, make error response
	if !isSuccessStatus(httpResponse.StatusCode) {
		// Wrap into a Fiber HTTP Error
		err = &errors.HTTPError{
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
