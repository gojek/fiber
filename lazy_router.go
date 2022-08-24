package fiber

import (
	"context"

	"github.com/gojek/fiber/errors"
	"github.com/gojek/fiber/util"
)

// LazyRouter implements Router interface and performs routing of incoming requests
// based on the routing strategy.
// The reason why it's 'lazy' is because it tries to dispatch an incoming request by
// a primary route first and switches to fallback options (one by one) only if
// received response is not OK
type LazyRouter struct {
	*BaseMultiRouteComponent

	strategy *baseRoutingStrategy
}

// NewLazyRouter initializes new LazyRouter
func NewLazyRouter(id string) *LazyRouter {
	if id == "" {
		id = "lazy-router_" + util.UID()
	}
	return &LazyRouter{
		BaseMultiRouteComponent: NewMultiRouteComponent(id),
	}
}

// SetStrategy sets routing strategy for this router
func (r *LazyRouter) SetStrategy(strategy RoutingStrategy) {
	r.strategy = &baseRoutingStrategy{RoutingStrategy: strategy}
}

// Dispatch makes a synchronous call to a routing strategy to select the primary route and fallbacks.
// After receiving a response it asynchronously asks a primary route to dispatch the request.
// If all responseQueue from a primary route are OK, it sends them back to output
// Otherwise it repeats the same with all fallback options one by one until one of fallbacks
// successfully dispatches a request or all fallbacks tried and failed to dispatch it
func (r *LazyRouter) Dispatch(ctx context.Context, req Request) ResponseQueue {
	ctx = r.beforeDispatch(ctx, req)
	out := make(chan Response, 1)

	queue := NewResponseQueue(out, 1)
	defer r.afterDispatch(ctx, req, queue)

	go func() {
		defer r.afterCompletion(ctx, req, queue)
		defer close(out)

		var routes []Component
		routesOrderCh, errCh := r.strategy.getRoutesOrder(ctx, req, r.routes)
		for routesOrderCh != nil || errCh != nil {
			select {
			case orderedRoutes, ok := <-routesOrderCh:
				if ok {
					routes = orderedRoutes
				} else {
					routesOrderCh = nil
				}
			case err, ok := <-errCh:
				if ok {
					out <- NewErrorResponse(errors.NewFiberError(req.Protocol(), err))
					return
				}
				errCh = nil
			case <-ctx.Done():
				out <- NewErrorResponse(errors.ErrRouterStrategyTimeoutExceeded(req.Protocol()))
				return
			}
		}

		if len(routes) > 0 {
			// iterate over an ordered slice of possible routes
			for _, route := range routes {
				copyReq, _ := req.Clone()
				responses := make([]Response, 0)
				responseCh := route.Dispatch(ctx, copyReq).Iter()
				ok := true
				for ok {
					select {
					case resp, notClosed := <-responseCh:
						if notClosed {
							if ok = resp.IsSuccess(); ok {
								responses = append(responses, resp.WithBackendName(route.ID()))
							}
						} else {
							// all responseQueue from selected route are ok, sending them back to output
							// and breaking a cycle over other routes
							for _, resp := range responses {
								out <- resp
							}
							return
						}
					case <-ctx.Done():
						out <- NewErrorResponse(errors.ErrRequestTimeout(req.Protocol()))
						return
					}
				}
			}
		} else {
			out <- NewErrorResponse(errors.ErrRouterStrategyReturnedEmptyRoutes(req.Protocol()))
		}
	}()

	return queue
}
