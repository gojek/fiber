package fiber

// MultiRouteComponent - is a network component with zero or more possible routes,
// such as FanOut, Combiner, Router
type MultiRouteComponent interface {
	Component

	SetRoutes(routes map[string]Component)
	GetRoutes() map[string]Component
}

// NewMultiRouteComponent is a factory function for creating a MultiRouteComponent
func NewMultiRouteComponent(id string) *BaseMultiRouteComponent {
	return &BaseMultiRouteComponent{
		BaseComponent: BaseComponent{id: id, kind: MultiRouteComponentKind},
	}
}

// BaseMultiRouteComponent is a reference implementation of a MultiRouteComponent
type BaseMultiRouteComponent struct {
	BaseComponent
	routes map[string]Component
}

// SetRoutes sets possible routes for this multi-route component
func (multiRoute *BaseMultiRouteComponent) SetRoutes(routes map[string]Component) {
	multiRoute.routes = routes
}

// GetRoutes is a getter for the routes configured on the BaseMultiRouteComponent
func (multiRoute *BaseMultiRouteComponent) GetRoutes() map[string]Component {
	return multiRoute.routes
}

// AddInterceptor can be used to (optionally, recursively) add one or more interceptors to
// the BaseMultiRouteComponent
func (multiRoute *BaseMultiRouteComponent) AddInterceptor(recursive bool, interceptors ...Interceptor) {
	if recursive {
		for _, route := range multiRoute.routes {
			route.AddInterceptor(recursive, interceptors...)
		}
	}
	multiRoute.BaseComponent.AddInterceptor(recursive, interceptors...)
}
