package fiber

// Router is a network component, that uses provided RoutingStrategy to
// select a route (child component), that should dispatch an incoming request
type Router interface {
	MultiRouteComponent

	// Sets routing strategy for this router
	SetStrategy(strategy RoutingStrategy)
}
