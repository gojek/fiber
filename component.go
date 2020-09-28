package fiber

import "context"

// ComponentKind can be used to define the types of Fiber components
// that support the Component interface
type ComponentKind string

const (
	// CallerKind represents a Fiber component that implements the Caller interface
	CallerKind ComponentKind = "Caller"
	// CombinerKind represents the Combiner type
	CombinerKind ComponentKind = "Combiner"
	// MultiRouteComponentKind represents a Fiber component that implements
	// the MultiRouteComponent interface
	MultiRouteComponentKind ComponentKind = "MultiRouteComponent"
)

// Component is the Base interface, that other network components should implement
type Component interface {
	Type

	// Returns component id
	ID() string

	// Returns the type of the encompassing structure, set at initialization
	Kind() ComponentKind

	// Dispatches the incoming request and returns a ResponseQueue
	// with zero or more responses in it
	Dispatch(ctx context.Context, req Request) ResponseQueue

	AddInterceptor(recursive bool, interceptors ...Interceptor)
}

// BaseComponent implements those contracts on the Component interface associated with
// the ID and Interceptors. Other network components can embed this type to re-use these
// methods.
type BaseComponent struct {
	BaseFiberType

	id string

	kind ComponentKind

	interceptors []Interceptor
}

// ID is the getter for the BaseComponent's unique ID
func (c *BaseComponent) ID() string {
	return c.id
}

// Kind is the getter for the type of the encompassing structure
func (c *BaseComponent) Kind() ComponentKind {
	return c.kind
}

func (c *BaseComponent) beforeDispatch(ctx context.Context, req Request) context.Context {
	// Add component id and type to the context
	ctx = context.WithValue(ctx, CtxComponentIDKey, c.ID())
	ctx = context.WithValue(ctx, CtxComponentKindKey, c.Kind())
	for _, i := range c.interceptors {
		ctx = i.BeforeDispatch(ctx, req)
	}
	return ctx
}

func (c *BaseComponent) afterDispatch(ctx context.Context, req Request, queue ResponseQueue) {
	for _, i := range c.interceptors {
		go i.AfterDispatch(ctx, req, queue)
	}
}

func (c *BaseComponent) afterCompletion(ctx context.Context, req Request, queue ResponseQueue) {
	for _, i := range c.interceptors {
		go i.AfterCompletion(ctx, req, queue)
	}
}

// AddInterceptor can be used to add one or more interceptors to the BaseComponent
func (c *BaseComponent) AddInterceptor(recursive bool, interceptors ...Interceptor) {
	c.interceptors = append(c.interceptors, interceptors...)
}

func NewBaseComponent(id string, kind ComponentKind) *BaseComponent {
	return &BaseComponent{
		id:   id,
		kind: kind,
	}
}
