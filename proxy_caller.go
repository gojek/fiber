package fiber

import "context"

// Proxy can be used to configure an intermediary for requests
type Proxy struct {
	Component
	backend Backend
}

// Dispatch is used to dispatch the incoming request against the proxy backend
func (p *Proxy) Dispatch(ctx context.Context, req Request) ResponseQueue {
	proxyReq, err := req.Transform(p.backend)

	if err != nil {
		return NewResponseQueueFromResponses(NewErrorResponse(err))
	}

	return p.Component.Dispatch(ctx, proxyReq)
}

// NewProxy is a factory function to create a new Proxy structure
func NewProxy(backend Backend, component Component) *Proxy {
	return &Proxy{
		Component: component,
		backend:   backend,
	}
}
