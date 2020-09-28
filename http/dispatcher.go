package http

import (
	"errors"
	"net/http"

	"github.com/gojek/fiber"
)

// Client is the base interface for an http-client (to be able to mock actual implementation)
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Dispatcher struct {
	httpClient Client
}

func (d *Dispatcher) Do(req fiber.Request) fiber.Response {
	if httpReq, ok := req.(*Request); ok {
		resp, err := d.httpClient.Do(httpReq.Request)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
			return NewHTTPResponse(resp)
		}
		return fiber.NewErrorResponse(err)
	}

	return fiber.NewErrorResponse(errors.New("fiber: http.Dispatcher supports only http.Request type of requests"))
}

func NewDispatcher(client Client) (fiber.Dispatcher, error) {
	if client == nil {
		return nil, errors.New("client can not be nil")
	}
	return &Dispatcher{
		httpClient: client,
	}, nil
}
