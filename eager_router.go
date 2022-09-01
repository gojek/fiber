package fiber

import (
	"context"

	"github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/util"
)

// EagerRouter implements Router interface and performs routing of incoming requests
// based on the routing strategy.
// The reason why it's 'eager' is because it dispatches incoming request by its every
// possible route in parallel and then returns either a response from a primary route
// (defined by the routing strategy) or switches back to one of fallback options.
//
// In a sense, EagerRouter is a Combiner, that aggregates responses from its all routes
// into a single response by selecting this response based on a provided RoutingStrategy
type EagerRouter struct {
	*Combiner
}

// NewEagerRouter initializes new EagerRouter
func NewEagerRouter(id string) *EagerRouter {
	if id == "" {
		id = "eager-router_" + util.UID()
	}
	return &EagerRouter{
		Combiner: NewCombiner(id),
	}
}

// SetStrategy sets routing strategy for this router
func (router *EagerRouter) SetStrategy(strategy RoutingStrategy) {
	router.WithFanIn(&eagerRouterFanIn{
		BaseFanIn{},
		&baseRoutingStrategy{RoutingStrategy: strategy},
		router})
}

// EagerRouter's specific FanIn implementation
// It receives the channel with responses from all possible router routes and asynchronously
// retrieves information about primary route and the order of fallbacks to be used.
//
// This FanIn doesn't wait for responses from all of the routes, but returns a response
// as soon, as the preferred order of routes is known and a successful response from
// the primary route is received.

// In case if the response from the primary route is not successful, then the first successful
// response from fallback routes will be sent back.

// If primary route AND all fallback routes responded with not non-successful responses, the error
// response will be created and sent back.
type eagerRouterFanIn struct {
	BaseFanIn
	strategy *baseRoutingStrategy
	router   *EagerRouter
}

func (fanIn *eagerRouterFanIn) Aggregate(
	ctx context.Context,
	req Request,
	queue ResponseQueue,
) Response {
	// use routing strategy to fetch primary route and fallbacks
	// publish the ordered routes into a channel
	routesOrderCh, errCh := fanIn.strategy.getRoutesOrder(ctx, req, fanIn.router.GetRoutes())

	out := make(chan Response, 1)
	go func() {
		defer close(out)

		var (
			// map to temporary store responseQueue from the routes
			responses = make(map[string]Response)

			// routes, ordered according to their priority
			// would be initialized from a routesOrderCh channel
			routes []Component

			// index of current primary route
			currentRouteIdx int

			responseCh = queue.Iter()

			masterResponse Response
		)

		for masterResponse == nil {
			select {
			case resp, ok := <-responseCh:
				if ok {
					responses[resp.BackendName()] = resp
				} else {
					responseCh = nil
				}
			case orderedRoutes, ok := <-routesOrderCh:
				if ok {
					routes = orderedRoutes
				} else {
					routesOrderCh = nil
				}
			case err, ok := <-errCh:
				if ok {
					masterResponse = NewErrorResponse(errors.NewFiberError(req.Protocol(), err))
				} else {
					errCh = nil
				}
			case <-ctx.Done():
				if routes == nil {
					// timeout exceeded, but no routes received. Sending error response
					masterResponse = NewErrorResponse(errors.ErrRouterStrategyTimeoutExceeded(req.Protocol()))
				} else {
					// timeout exceeded
					responseCh = nil
				}
			}

			if routes != nil {
				for ; currentRouteIdx < len(routes); currentRouteIdx++ {
					if currMasterResponse, exist := responses[routes[currentRouteIdx].ID()]; exist {
						if currMasterResponse.IsSuccess() {
							// preferred response found
							masterResponse = currMasterResponse
							break
						}
					} else if responseCh != nil {
						// response from preferred route is not ready; continue listening for new responseQueue
						break
					}
				}

				// all expected routes tried, no OK response received from either of them
				if currentRouteIdx >= len(routes) {
					if len(routes) == 0 {
						masterResponse = NewErrorResponse(errors.ErrRouterStrategyReturnedEmptyRoutes(req.Protocol()))
					} else {
						masterResponse = NewErrorResponse(errors.ErrServiceUnavailable(req.Protocol()))
					}
				}
			}
		}
		out <- masterResponse
	}()

	return <-out
}
