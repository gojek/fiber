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
	Label(key string) []string
	WithLabel(key string, values ...string) Response
	WithLabels(Labels) Response
}

type ErrorResponse struct {
	*CachedPayload
	labels  Labels
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

func (resp *ErrorResponse) Label(key string) []string {
	return resp.labels.Label(key)
}

func (resp *ErrorResponse) WithLabel(key string, values ...string) Response {
	resp.labels = resp.labels.WithLabel(key, values...)
	return resp
}

func (resp *ErrorResponse) WithLabels(labels Labels) Response {
	for _, key := range labels.Keys() {
		resp.labels = resp.labels.WithLabel(key, labels.Label(key)...)
	}
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
		labels:        NewLabelsMap(),
	}
}
