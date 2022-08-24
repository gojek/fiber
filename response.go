package fiber

import (
	"github.com/gojek/fiber/errors"
)

type Response interface {
	IsSuccess() bool
	Payload() interface{}
	StatusCode() int
	BackendName() string
	WithBackendName(string) Response
}

type ErrorResponse struct {
	*CachedPayload
	code    int
	backend string
}

func (resp *ErrorResponse) IsSuccess() bool {
	return false
}

func (resp *ErrorResponse) BackendName() string {
	return resp.backend
}

func (resp *ErrorResponse) WithBackendName(backendName string) Response {
	resp.backend = backendName
	return resp
}

func (resp *ErrorResponse) StatusCode() int {
	return resp.code
}

func NewErrorResponse(err error) Response {
	var fiberErr *errors.FiberError
	if castedError, ok := err.(*errors.FiberError); ok {
		fiberErr = castedError
	} else {
		fiberErr = errors.NewFiberError(HTTP, err)
	}
	payload, _ := fiberErr.ToJSON()
	return &ErrorResponse{
		CachedPayload: NewCachedPayload(payload),
		code:          fiberErr.Code,
	}
}
