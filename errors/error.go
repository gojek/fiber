package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
)

// FiberError is used to capture the error resulting from a HTTP request
type FiberError struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

// Error is a getter for the error message in a FiberError object
func (err FiberError) Error() string {
	return err.Message
}

// ToJSON returns the FiberError object as a Json encoded byte array
func (err *FiberError) ToJSON() ([]byte, error) {
	return json.MarshalIndent(err, "", "  ")
}

// NewFiberError returns an error of type FiberError from the input error object.
// If the input error is already of the required type, it is returned as is.
// If not, a generic request failed error is created from the given error.
func NewFiberError(protocol string, err error) *FiberError {
	if fiberError, ok := err.(FiberError); ok {
		return &fiberError
	}
	return ErrRequestFailed(protocol, err)
}

var (
	// ErrRouterStrategyTimeoutExceeded is a FiberError that's returned when
	// the routing strategy fails to respond within given timeout
	ErrRouterStrategyTimeoutExceeded = func(protocol string) *FiberError {
		statusCode := http.StatusInternalServerError
		if protocol == "GRPC" {
			statusCode = int(codes.Internal)
		}
		return &FiberError{
			Code:    statusCode,
			Message: "fiber: routing strategy failed to respond within given timeout",
		}
	}

	// ErrRouterStrategyReturnedEmptyRoutes is a FiberError that's returned when
	// the routing strategy routing strategy returns an empty routes list
	ErrRouterStrategyReturnedEmptyRoutes = func(protocol string) *FiberError {
		statusCode := http.StatusNotImplemented
		if protocol == "GRPC" {
			statusCode = int(codes.Unimplemented)
		}
		return &FiberError{
			Code:    statusCode,
			Message: "fiber: routing strategy returned empty routes list",
		}
	}

	// ErrServiceUnavailable is a FiberError that's returned when
	// none of the routes in the routing strategy return a valid response
	ErrServiceUnavailable = func(protocol string) *FiberError {
		statusCode := http.StatusServiceUnavailable
		if protocol == "GRPC" {
			statusCode = int(codes.Unavailable)
		}
		return &FiberError{
			Code:    statusCode,
			Message: "fiber: no responses received",
		}
	}

	// ErrRequestTimeout is a FiberError that's returned when
	// no response if received for a given HTTP request within the configured timeout
	ErrRequestTimeout = func(protocol string) *FiberError {
		statusCode := http.StatusRequestTimeout
		if protocol == "GRPC" {
			statusCode = int(codes.DeadlineExceeded)
		}
		return &FiberError{
			Code:    statusCode,
			Message: "fiber: failed to receive a response within configured timeout",
		}
	}

	// ErrReadRequestFailed is a FiberError that's returned when a http request cannot
	// be read successfully
	ErrReadRequestFailed = func(protocol string, err error) *FiberError {
		statusCode := http.StatusInternalServerError
		if protocol == "GRPC" {
			statusCode = int(codes.Internal)
		}
		return &FiberError{
			Code:    statusCode,
			Message: fmt.Sprintf("fiber: failed to read incoming request: %s", err.Error()),
		}
	}

	// ErrRequestFailed is a generic error that is created when problems are encountered fulfilling
	// a request
	ErrRequestFailed = func(protocol string, err error) *FiberError {
		statusCode := http.StatusInternalServerError
		if protocol == "GRPC" {
			statusCode = int(codes.Internal)
		}
		return &FiberError{
			Code:    statusCode,
			Message: fmt.Sprintf("fiber: request cannot be completed: %s", err.Error()),
		}
	}

	ErrInvalidInput = func(protocol string, err error) *FiberError {
		statusCode := http.StatusBadRequest
		if protocol == "GRPC" {
			statusCode = int(codes.InvalidArgument)
		}
		return &FiberError{
			Code:    statusCode,
			Message: fmt.Sprintf("fiber: %s", err.Error()),
		}
	}
)
