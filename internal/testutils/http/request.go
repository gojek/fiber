package http

import (
	"bytes"
	"net/http"

	fiberHTTP "github.com/gojek/fiber/http"
)

func MockReq(method, url, body string) *fiberHTTP.Request {
	httpReq, _ := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req, _ := fiberHTTP.NewHTTPRequest(httpReq)
	return req
}
