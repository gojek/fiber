package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/protocol"
)

// Request wraps a standard http request
type Request struct {
	*fiber.CachedPayload
	*http.Request
}

func (r *Request) Protocol() protocol.Protocol {
	return protocol.HTTP
}

func (r *Request) Header() map[string][]string {
	return r.Request.Header
}

// NewHTTPRequest initialize a new client request from incoming server request
func NewHTTPRequest(req *http.Request) (*Request, error) {
	// RequestURI can't be set in client requests
	req.RequestURI = ""

	var payload *fiber.CachedPayload
	if req.ContentLength == 0 && req.Body == nil {
		payload = new(fiber.CachedPayload)
		req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
	} else {
		defer req.Body.Close()

		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		payload = fiber.NewCachedPayload(data)
		req.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(data)
			return ioutil.NopCloser(r), nil
		}
	}

	return &Request{Request: req, CachedPayload: payload}, nil
}

// Copy creates a deep copy of this request
func (r *Request) Clone() (fiber.Request, error) {
	bytePayLoad, ok := r.Payload().([]byte)
	if !ok {
		return nil, errors.New("unable to parse payload")
	}
	bodyReader := bytes.NewReader(bytePayLoad)

	proxyRequest, err := http.NewRequest(r.Method, r.URL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	proxyRequest.GetBody = func() (io.ReadCloser, error) {
		return ioutil.NopCloser(bodyReader), nil
	}

	proxyRequest.Header = r.Header()

	return &Request{CachedPayload: r.CachedPayload, Request: proxyRequest}, nil
}

func (r *Request) OperationName() string {
	return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
}

func (r *Request) Transform(backend fiber.Backend) (fiber.Request, error) {
	reqPath := ""
	// If the URL is set, extract the URI. If nothing is set, ignore.
	// Calling RequestURI() on the empty url will return a '/' and the
	// backend server may or may not handle it, so skip instead.
	if r.URL.String() != "" {
		reqPath = r.URL.RequestURI()
	}

	updatedURL, err := url.Parse(backend.URL(reqPath))
	if err != nil {
		return nil, err
	}

	r.URL = updatedURL
	return r, nil
}
