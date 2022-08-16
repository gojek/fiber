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
					testutils.DelayedResponse{Response: testutils.MockResp(200, "A-OK", nil, nil)}),
				"route-b": testutils.NewMockComponent(
					"route-b",
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)}),
			},
			strategy: []string{
				"route-b", "route-a",
			},
			expected: []fiber.Response{
				testutils.MockResp(200, "B-OK", nil, nil).WithBackendName("route-b"),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name: "ok: first route failed, fallback response succeeded",
			routes: map[string]fiber.Component{
				"route-a": testutils.NewMockComponent(
					"route-a",
					testutils.DelayedResponse{Response: testutils.MockResp(500, "A-NOK", nil, fiberErrors.ErrServiceUnavailable(fiber.HTTP.String()))}),
				"route-b": testutils.NewMockComponent(
					"route-b",
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)}),
			},
			strategy: []string{
				"route-a", "route-b",
			},
			expected: []fiber.Response{
				testutils.MockResp(200, "B-OK", nil, nil).WithBackendName("route-b"),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name: "error: routing strategy succeeded, but route timeout exceeded",
			routes: map[string]fiber.Component{
				"route-a": testutils.NewMockComponent(
					"route-a",
					testutils.DelayedResponse{
						Latency:  100 * time.Millisecond,
						Response: testutils.MockResp(200, "A-OK", nil, nil)}),
				"route-b": testutils.NewMockComponent(
					"route-b",
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)}),
			},
			strategy: []string{
				"route-a", "route-b",
			},
			strategyLatency: 50 * time.Millisecond,
			expected: []fiber.Response{
				testutils.MockResp(408, "", nil, fiberErrors.ErrRequestTimeout(fiber.HTTP.String())),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:            "error: strategy timeout exceeded",
			strategyLatency: 200 * time.Millisecond,
			expected: []fiber.Response{
				testutils.MockResp(500, "", nil, fiberErrors.ErrRouterStrategyTimeoutExceeded(fiber.HTTP.String())),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:     "error: routing strategy returned empty routes",
			strategy: []string{},
			expected: []fiber.Response{
				testutils.MockResp(501, "", nil, fiberErrors.ErrRouterStrategyReturnedEmptyRoutes(fiber.HTTP.String())),
			},
			timeout: 100 * time.Millisecond,
		},
		{
			name:              "error: routing strategy responded with exception",
			strategyException: errors.New("unexpected exception happened"),
			expected: []fiber.Response{
				testutils.MockResp(500, "", nil, fiberErrors.NewFiberError(fiber.HTTP.String(), errors.New("unexpected exception happened"))),
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
			request := testutils.MockReq("POST", "http://localhost:8080/lazy-router", "payload")
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
