package fiber_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/internal/testutils"
	testUtilsHttp "github.com/gojek/fiber/internal/testutils/http"
	"github.com/gojek/fiber/protocol"
	"github.com/stretchr/testify/assert"
)

type fanOutTestCase struct {
	name      string
	responses map[string][]testUtilsHttp.DelayedResponse
}

func (tt *fanOutTestCase) Routes() map[string]fiber.Component {
	routes := make(map[string]fiber.Component)
	for name, resp := range tt.responses {
		routes[name] = testutils.NewMockComponent(name, resp...)
	}
	return routes
}

func (tt *fanOutTestCase) ExpectedResponses(timeout time.Duration) map[string][]fiber.Response {
	expectedResponses := make(map[string][]fiber.Response)
	for route, responses := range tt.responses {
		commutativeLatency := time.Duration(0)
		for _, resp := range responses {
			if commutativeLatency += resp.Latency; commutativeLatency < timeout {
				expectedResponses[route] = append(expectedResponses[route], resp.Response)
			} else {
				break
			}
		}
	}
	return expectedResponses
}

func TestFanOut_Dispatch(t *testing.T) {
	timeout := 100 * time.Millisecond

	suite := []fanOutTestCase{
		{
			name: "two routes/two OK responseQueue",
			responses: map[string][]testUtilsHttp.DelayedResponse{
				"route-a": {
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "A-OK", nil, nil)},
				},
				"route-b": {
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "B-OK", nil, nil)},
				},
			},
		},
		{
			name: "two routes/one OK response",
			responses: map[string][]testUtilsHttp.DelayedResponse{
				"route-a": {
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "A-OK", nil, nil)},
				},
				"route-b": {},
			},
		},
		{
			name: "two routes/two OK responseQueue in each",
			responses: map[string][]testUtilsHttp.DelayedResponse{
				"route-a": {
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "A-OK_1", nil, nil)},
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "A-OK_2", nil, nil)},
				},
				"route-b": {
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "B-OK_1", nil, nil)},
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "B-OK_2", nil, nil)},
				},
			},
		},
		{
			name: "one route/one NOK response",
			responses: map[string][]testUtilsHttp.DelayedResponse{
				"route-a": {
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(502, "", nil, errors.ErrNoValidResponseFromRoutes(protocol.HTTP))},
				},
			},
		},
		{
			name: "one route/multiple responseQueue with delays",
			responses: map[string][]testUtilsHttp.DelayedResponse{
				"route-a": {
					testUtilsHttp.DelayedResponse{
						Latency:  10 * time.Millisecond,
						Response: testUtilsHttp.MockResp(200, "OK", nil, nil),
					},
					testUtilsHttp.DelayedResponse{
						Latency:  timeout / 2,
						Response: testUtilsHttp.MockResp(200, "OK", nil, nil),
					},
					// should never be received
					testUtilsHttp.DelayedResponse{
						Latency:  2 * timeout,
						Response: testUtilsHttp.MockResp(200, "OK", nil, nil),
					},
				},
			},
		},
		{
			name: "two routes/one OK, one timeout",
			responses: map[string][]testUtilsHttp.DelayedResponse{
				"route-a": {
					testUtilsHttp.DelayedResponse{
						Latency:  2 * timeout,
						Response: testUtilsHttp.MockResp(200, "OK", nil, nil),
					},
				},
				"route-b": {
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "OK", nil, nil)},
				},
			},
		},
	}

	for _, tt := range suite {
		fanOut := fiber.NewFanOut("")
		fanOut.SetRoutes(tt.Routes())

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		var (
			receivedResponses = make(map[string][]fiber.Response)
			expectedResponses = tt.ExpectedResponses(timeout)
		)

		for responsesCh := fanOut.Dispatch(ctx, testUtilsHttp.MockReq("GET", "http://test:8080", "")).Iter(); ; {
			select {
			case resp, ok := <-responsesCh:
				if ok {
					receivedResponses[resp.BackendName()] = append(receivedResponses[resp.BackendName()], resp)
					continue
				}
			case <-time.After(timeout + timeout/2):
				assert.Fail(t, fmt.Sprintf("[%s] failed: it didn't terminate after a timeout...", tt.name))
			}

			cancel()
			break
		}

		assert.Equal(t, len(expectedResponses), len(receivedResponses))
		for name, resp := range expectedResponses {
			assert.ElementsMatch(t, resp, receivedResponses[name])
		}
	}
}
