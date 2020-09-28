package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPError is used to capture the error resulting from a HTTP request
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

// Error is a getter for the error message in a HTTPError object
func (err HTTPError) Error() string {
	return err.Message
}

// ToJSON returns the HTTPError object as a Json encoded byte array
func (err *HTTPError) ToJSON() ([]byte, error) {
	return json.MarshalIndent(err, "", "  ")
}

// NewHTTPError returns an error of type HTTPError from the input error object.
// If the input error is already of the required type, it is returned as is.
// If not, a generic request failed error is created from the given error.
func NewHTTPError(err error) *HTTPError {
	if httpErr, ok := err.(HTTPError); ok {
		return &httpErr
	}
	return ErrRequestFailed(err)
}

var (
	// ErrRouterStrategyTimeoutExceeded is a HTTPError that's returned when
	// the routing strategy fails to respond within given timeout
	ErrRouterStrategyTimeoutExceeded = &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: "fiber: routing strategy failed to respond within given timeout",
	}

	// ErrRouterStrategyReturnedEmptyRoutes is a HTTPError that's returned when
	// the routing strategy routing strategy returns an empty routes list
	ErrRouterStrategyReturnedEmptyRoutes = &HTTPError{
		Code:    http.StatusNotImplemented,
		Message: "fiber: routing strategy returned empty routes list",
	}

	// ErrServiceUnavailable is a HTTPError that's returned when
	// none of the routes in the routing strategy return a valid response
	ErrServiceUnavailable = &HTTPError{
		Code:    http.StatusServiceUnavailable,
		Message: "fiber: no responses received",
	}

	// ErrRequestTimeout is a HTTPError that's returned when
	// no response if received for a given HTTP request within the configured timeout
	ErrRequestTimeout = &HTTPError{
		Code:    http.StatusRequestTimeout,
		Message: "fiber: failed to receive a response within configured timeout",
	}

	// ErrReadRequestFailed is a HTTPError that's returned when a http request cannot
	// be read successfully
	ErrReadRequestFailed = func(err error) *HTTPError {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("fiber: failed to read incoming request: %s", err.Error()),
		}
	}

	// ErrRequestFailed is a generic error that is created when problems are encountered fulfilling
	// a request
	ErrRequestFailed = func(err error) *HTTPError {
		return &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("fiber: request cannot be completed: %s", err.Error()),
		}
	}
)
