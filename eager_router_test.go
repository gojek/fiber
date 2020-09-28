package fiber_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/gojek/fiber"
	fiberErrors "github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/internal/testutils"
)

type eagerRouterTestCase struct {
	name              string
	responses         map[string][]testutils.DelayedResponse
	routes            map[string]fiber.Component
	strategy          []string
	strategyLatency   time.Duration
	strategyException error
	expected          []fiber.Response
}

func (tt *eagerRouterTestCase) Routes() map[string]fiber.Component {
	if tt.routes == nil {
		tt.routes = make(map[string]fiber.Component)
		for name, resp := range tt.responses {
			tt.routes[name] = testutils.NewMockComponent(name, resp...)
		}
	}
	return tt.routes
}

func TestEagerRouter_Dispatch(t *testing.T) {
	timeout := 100 * time.Millisecond
	suite := []eagerRouterTestCase{
		{
			name: "all OK responseQueue, no delays",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "A-OK", nil, nil)},
				},
				"route-b": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)},
				},
			},
			strategy: []string{
				"route-b", "route-a",
			},
			expected: []fiber.Response{
				testutils.MockResp(200, "B-OK", nil, nil).WithBackendName("route-b"),
			},
		},
		{
			name: "primary route failed, receiving from fallback",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{Response: testutils.MockResp(500, "A-NOK", nil, fiberErrors.ErrServiceUnavailable)},
				},
				"route-b": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)},
				},
			},
			strategy: []string{
				"route-a", "route-b",
			},
			expected: []fiber.Response{
				testutils.MockResp(200, "B-OK", nil, nil).WithBackendName("route-b"),
			},
		},
		{
			name: "expected response comes after the fallback response",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{
						Latency:  75 * time.Millisecond,
						Response: testutils.MockResp(200, "A-OK", nil, nil),
					},
				},
				"route-b": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)},
				},
			},
			strategy: []string{
				"route-a", "route-b",
			},
			expected: []fiber.Response{
				testutils.MockResp(200, "A-OK", nil, nil).WithBackendName("route-a"),
			},
		},
		{
			name: "primary route exceeds timeout, receiving from fallback",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{
						Latency:  timeout + 10*time.Millisecond,
						Response: testutils.MockResp(200, "A-OK", nil, nil),
					},
				},
				"route-b": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)},
				},
			},
			strategy: []string{
				"route-a", "route-b",
			},
			expected: []fiber.Response{
				testutils.MockResp(200, "B-OK", nil, nil).WithBackendName("route-b"),
			},
		},
		{
			name: "primary route exceeds timeout, fallback route failed, receiving from the next fallback",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{Response: testutils.MockResp(500, "A-NOK", nil, fiberErrors.ErrServiceUnavailable)},
				},
				"route-b": {
					testutils.DelayedResponse{
						Response: testutils.MockResp(200, "B-OK", nil, nil),
						Latency:  timeout + 10*time.Millisecond,
					},
				},
				"route-c": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "C-OK", nil, nil)},
				},
			},
			strategy: []string{
				"route-a", "route-b", "route-c",
			},
			expected: []fiber.Response{
				testutils.MockResp(200, "C-OK", nil, nil).WithBackendName("route-c"),
			},
		},
		{
			name: "primary route and all fallbacks failed, receiving error response",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{
						Response: testutils.MockResp(500, "A-NOK", nil, fiberErrors.ErrServiceUnavailable),
					},
				},
				"route-b": {
					testutils.DelayedResponse{
						Response: testutils.MockResp(200, "B-OK", nil, nil),
						Latency:  timeout + 10*time.Millisecond,
					},
				},
				"route-c": {
					testutils.DelayedResponse{
						Response: testutils.MockResp(408, "C-NOK", nil, fiberErrors.ErrRequestTimeout),
					},
				},
			},
			strategy: []string{
				"route-a", "route-b", "route-c",
			},
			expected: []fiber.Response{
				testutils.MockResp(500, "", nil, fiberErrors.ErrServiceUnavailable),
			},
		},
		{
			name: "routing strategy response comes after all routes replied",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{Response: testutils.MockResp(500, "A-NOK", nil, fiberErrors.ErrServiceUnavailable)},
				},
				"route-b": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)},
				},
			},
			strategyLatency: timeout / 2,
			strategy: []string{
				"route-a", "route-b",
			},
			expected: []fiber.Response{
				testutils.MockResp(200, "B-OK", nil, nil).WithBackendName("route-b"),
			},
		},
		{
			name: "routing strategy response exceeds timeout",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "A-OK", nil, nil)},
				},
				"route-b": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)},
				},
			},
			strategyLatency: timeout + timeout,
			strategy: []string{
				"route-a", "route-b",
			},
			expected: []fiber.Response{
				testutils.MockResp(500, "", nil, fiberErrors.ErrRouterStrategyTimeoutExceeded),
			},
		},
		{
			name: "routing strategy returned empty routes",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "A-OK", nil, nil)},
				},
				"route-b": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)},
				},
			},
			strategy: []string{},
			expected: []fiber.Response{
				testutils.MockResp(501, "", nil, fiberErrors.ErrRouterStrategyReturnedEmptyRoutes),
			},
		},
		{
			name: "routing strategy failed with exception",
			responses: map[string][]testutils.DelayedResponse{
				"route-a": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "A-OK", nil, nil)},
				},
				"route-b": {
					testutils.DelayedResponse{Response: testutils.MockResp(200, "B-OK", nil, nil)},
				},
			},
			strategyException: errors.New("routing strategy exception"),
			expected: []fiber.Response{
				testutils.MockResp(500, "", nil, fiberErrors.NewHTTPError(errors.New("routing strategy exception"))),
			},
		},
	}

	for _, tt := range suite {
		router := fiber.NewEagerRouter("eager-router")
		router.SetRoutes(tt.Routes())

		strategy := testutils.NewMockRoutingStrategy(
			tt.Routes(),
			tt.strategy,
			tt.strategyLatency,
			tt.strategyException)
		router.SetStrategy(strategy)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		received := make([]fiber.Response, 0)
		request := testutils.MockReq("GET", "http://test:8080", "")
		for responsesCh := router.Dispatch(ctx, request).Iter(); ; {
			select {
			case resp, ok := <-responsesCh:
				if ok {
					received = append(received, resp)
					continue
				}
			case <-time.After(timeout + timeout/2):
				assert.Fail(t, fmt.Sprintf("[%s] failed: it didn't terminate after a timeout...", tt.name))
			}

			cancel()
			break
		}

		assert.Equal(t, len(tt.expected), len(received), tt.name)
		for i := 0; i < len(tt.expected); i++ {
			assert.Equal(t, string(tt.expected[i].Payload()), string(received[i].Payload()), tt.name)
			assert.Equal(t, tt.expected[i].StatusCode(), received[i].StatusCode(), tt.name)
		}
		strategy.AssertExpectations(t)
	}
}
