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
	) (route Component, fallbacks []Component, attr Attributes, err error)
}

type baseRoutingStrategy struct {
	RoutingStrategy
	BaseFiberType
}

func (s *baseRoutingStrategy) getRoutesOrder(
	ctx context.Context,
	req Request,
	routes map[string]Component,
) (<-chan []Component, <-chan error) {
	out := make(chan []Component)
	errCh := make(chan error, 1)

	go func() {
		route, fallbacks, _, err := s.SelectRoute(ctx, req, routes)

		if err != nil {
			errCh <- err
		} else {
			// Append routes
			routes := fallbacks
			if route != nil {
				routes = append([]Component{route}, routes...)
			}
			out <- routes
		}
		// Close both channels
		close(out)
		close(errCh)
	}()

	return out, errCh
}
