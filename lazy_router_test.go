package fiber_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gojek/fiber"
	fiberErrors "github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/internal/testutils"
	testUtilsHttp "github.com/gojek/fiber/internal/testutils/http"
	"github.com/gojek/fiber/protocol"
	"github.com/stretchr/testify/assert"
)

type lazyRouterTestCase struct {
	name              string
	routes            map[string]fiber.Component
	strategy          []string
	strategyLatency   time.Duration
	strategyException error
	expected          []fiber.Response
	timeout           time.Duration
}

func TestLazyRouter_Dispatch(t *testing.T) {
	suite := []lazyRouterTestCase{
		{
			name: "ok: successful response",
			routes: map[string]fiber.Component{
				"route-a": testutils.NewMockComponent(
					"route-a",
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "A-OK", nil, nil)}),
				"route-b": testutils.NewMockComponent(
					"route-b",
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "B-OK", nil, nil)}),
			},
			strategy: []string{
				"route-b", "route-a",
			},
			expected: []fiber.Response{
				testUtilsHttp.MockResp(200, "B-OK", nil, nil).WithBackendName("route-b"),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name: "ok: first route failed, fallback response succeeded",
			routes: map[string]fiber.Component{
				"route-a": testutils.NewMockComponent(
					"route-a",
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(500, "A-NOK", nil, fiberErrors.ErrServiceUnavailable(protocol.HTTP))}),
				"route-b": testutils.NewMockComponent(
					"route-b",
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "B-OK", nil, nil)}),
			},
			strategy: []string{
				"route-a", "route-b",
			},
			expected: []fiber.Response{
				testUtilsHttp.MockResp(200, "B-OK", nil, nil).WithBackendName("route-b"),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name: "error: no route succeeded",
			routes: map[string]fiber.Component{
				"route-a": testutils.NewMockComponent(
					"route-a",
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(500, "A-NOK", nil, fiberErrors.ErrServiceUnavailable(protocol.HTTP))}),
				"route-b": testutils.NewMockComponent(
					"route-b",
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(500, "B-NOK", nil, nil)}),
			},
			strategy: []string{
				"route-a", "route-b",
			},
			expected: []fiber.Response{
				testUtilsHttp.MockResp(501, "", nil, fiberErrors.ErrRouterStrategyReturnedEmptyRoutes(protocol.HTTP)),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name: "error: routing strategy succeeded, but route timeout exceeded",
			routes: map[string]fiber.Component{
				"route-a": testutils.NewMockComponent(
					"route-a",
					testUtilsHttp.DelayedResponse{
						Latency:  100 * time.Millisecond,
						Response: testUtilsHttp.MockResp(200, "A-OK", nil, nil)}),
				"route-b": testutils.NewMockComponent(
					"route-b",
					testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "B-OK", nil, nil)}),
			},
			strategy: []string{
				"route-a", "route-b",
			},
			strategyLatency: 50 * time.Millisecond,
			expected: []fiber.Response{
				testUtilsHttp.MockResp(408, "", nil, fiberErrors.ErrRequestTimeout(protocol.HTTP)),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:            "error: strategy timeout exceeded",
			strategyLatency: 200 * time.Millisecond,
			expected: []fiber.Response{
				testUtilsHttp.MockResp(500, "", nil, fiberErrors.ErrRouterStrategyTimeoutExceeded(protocol.HTTP)),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:     "error: routing strategy returned empty routes",
			strategy: []string{},
			expected: []fiber.Response{
				testUtilsHttp.MockResp(501, "", nil, fiberErrors.ErrRouterStrategyReturnedEmptyRoutes(protocol.HTTP)),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:              "error: routing strategy responded with exception",
			strategyException: errors.New("unexpected exception happened"),
			expected: []fiber.Response{
				testUtilsHttp.MockResp(500, "", nil, fiberErrors.NewFiberError(protocol.HTTP, errors.New("unexpected exception happened"))),
			},
			timeout: 100 * time.Millisecond,
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {
			router := fiber.NewLazyRouter("lazy-router")
			router.SetRoutes(tt.routes)

			strategy := testutils.NewMockRoutingStrategy(
				tt.routes,
				tt.strategy,
				tt.strategyLatency,
				tt.strategyException)
			router.SetStrategy(strategy)

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)

			received := make([]fiber.Response, 0)
			request := testUtilsHttp.MockReq("POST", "http://localhost:8080/lazy-router", "payload")
			for responsesCh := router.Dispatch(ctx, request).Iter(); ; {
				select {
				case resp, ok := <-responsesCh:
					if ok {
						received = append(received, resp)
						continue
					}
				case <-time.After(tt.timeout + tt.timeout/2):
					assert.Fail(t, fmt.Sprintf("[%s] failed: it didn't terminate after a timeout...", tt.name))
				}

				cancel()
				break
			}

			assert.Equal(t, len(tt.expected), len(received), tt.name)
			for i := 0; i < len(tt.expected); i++ {
				assert.Equal(t, tt.expected[i].Payload(), received[i].Payload(), tt.name)
				assert.Equal(t, tt.expected[i].StatusCode(), received[i].StatusCode(), tt.name)
			}
			strategy.AssertExpectations(t)
		})
	}
}