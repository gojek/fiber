package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gojek/fiber"
	fiberErrors "github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/protocol"
)

// Options captures a set of options that can be used as configurations for
// the Request handler
type Options struct {
	Timeout time.Duration
}

// Handler is a structure used to capture a fiber component and a set of
// options for making requests
type Handler struct {
	fiber.Component

	options Options
}

// NewHandler is a creator factory for the Handler
func NewHandler(c fiber.Component, options Options) *Handler {
	return &Handler{
		Component: c,
		options:   options,
	}
}

// ServeHTTP takes an incoming request, dipatches it on the fiber component
// and writes the response using the given ResponseWriter
func (h *Handler) ServeHTTP(writer http.ResponseWriter, httpReq *http.Request) {
	resp, err := h.DoRequest(httpReq)
	if err != nil {
		// Create error response
		resp = fiber.NewErrorResponse(err)
	}
	if err := h.write(resp, writer); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

// DoRequest executes the given http request and returns the response / error
func (h *Handler) DoRequest(httpReq *http.Request) (fiber.Response, *fiberErrors.FiberError) {
	if req, err := NewHTTPRequest(httpReq); err == nil {
		ctx, cancel := context.WithTimeout(req.Context(), h.options.Timeout)
		defer cancel()

		select {
		case resp, ok := <-h.Dispatch(ctx, req).Iter():
			if ok {
				return resp, nil
			}
			return nil, fiberErrors.ErrServiceUnavailable(protocol.HTTP)
		case <-time.After(h.options.Timeout):
			return nil, fiberErrors.ErrRequestTimeout(protocol.HTTP)
		}
	} else {
		return nil, fiberErrors.ErrReadRequestFailed(protocol.HTTP, err)
	}
}

// write takes a response and writes its contents to the given writer
func (h *Handler) write(resp fiber.Response, writer http.ResponseWriter) (err error) {
	if httpResp, ok := resp.(*Response); ok {
		for key, values := range httpResp.Header() {
			for i := range values {
				writer.Header().Add(key, values[i])
			}
		}
	}

	writer.WriteHeader(resp.StatusCode())
	bytePayLoad, ok := resp.Payload().([]byte)
	if !ok {
		return fiberErrors.NewFiberError(protocol.HTTP, errors.New("unable to parse payload"))
	}
	_, err = writer.Write(bytePayLoad)
	return err
}
