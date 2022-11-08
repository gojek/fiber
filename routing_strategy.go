package fiber

import "context"

// RoutingStrategy picks up primary route and zero or more fallbacks
// from the map of router routes
type RoutingStrategy interface {
	Type
	// req - Incoming request (so the route can be selected based on the request)
	// routes - map of all possible routes
	SelectRoute(ctx context.Context,
		req Request,
		routes map[string]Component,
	) (route Component, fallbacks []Component, labels Labels, err error)
}

type baseRoutingStrategy struct {
	RoutingStrategy
	BaseFiberType
}

type routesOrderResponse struct {
	Components []Component
	Labels     Labels
	Err        error
}

func (s *baseRoutingStrategy) getRoutesOrder(
	ctx context.Context,
	req Request,
	routes map[string]Component,
) <-chan routesOrderResponse {
	out := make(chan routesOrderResponse)

	go func() {
		route, fallbacks, labels, err := s.SelectRoute(ctx, req, routes)

		// Combine preferred route with the fallbacks
		routes := fallbacks
		if route != nil {
			routes = append([]Component{route}, routes...)
		}

		out <- routesOrderResponse{
			Components: routes,
			Err:        err,
			Labels:     labels,
		}
	}()

	return out
}
