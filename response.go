package fiber

import (
	"github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/protocol"
)

type Response interface {
	IsSuccess() bool
	Payload() []byte
	StatusCode() int
	BackendName() string
	WithBackendName(string) Response
	Attribute(key string) []string
	WithAttribute(key string, values ...string) Response
}

type ErrorResponse struct {
	*CachedPayload
	attr    Attributes
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

func (resp *ErrorResponse) Attribute(key string) []string {
	return resp.attr.Attribute(key)
}

func (resp *ErrorResponse) WithAttribute(key string, values ...string) Response {
	resp.attr = resp.attr.WithAttribute(key, values...)
	return resp
}

func NewErrorResponse(err error) Response {
	var fiberErr *errors.FiberError
	if castedError, ok := err.(*errors.FiberError); ok {
		fiberErr = castedError
	} else {
		fiberErr = errors.NewFiberError(protocol.HTTP, err)
	}
	payload, _ := fiberErr.ToJSON()
	return &ErrorResponse{
		CachedPayload: NewCachedPayload(payload),
		code:          fiberErr.Code,
		attr:          NewAttributesMap(),
	}
}
