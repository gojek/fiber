package fiber

import (
	"net/http"

	"github.com/gojek/fiber/errors"
)

type Response interface {
	IsSuccess() bool
	Payload() []byte
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
	var httpErr *errors.HTTPError

	if castedError, ok := err.(*errors.HTTPError); ok {
		httpErr = castedError
	} else {
		httpErr = &errors.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	payload, _ := httpErr.ToJSON()
	return &ErrorResponse{
		CachedPayload: NewCachedPayload(payload),
		code:          httpErr.Code,
	}
}
