package testutils

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gojek/fiber"
	fiberHTTP "github.com/gojek/fiber/http"
)

func MockResp(code int, body string, header http.Header, err error) fiber.Response {
	if err != nil {
		return fiber.NewErrorResponse(err)
	}

	httpResp := &http.Response{
		StatusCode: code,
		Header:     header,
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
	}
	return fiberHTTP.NewHTTPResponse(httpResp)
}

type DelayedResponse struct {
	fiber.Response
	Latency time.Duration
}
