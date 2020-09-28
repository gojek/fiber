package fiber

import (
	"context"
)

// CtxKey is an alias for a string and is used to associate keys to context objects
type CtxKey string

var (
	// CtxComponentIDKey is used to denote the component's id in the request context
	CtxComponentIDKey CtxKey = "CTX_COMPONENT_ID"
	// CtxComponentKindKey is used to denote the component's kind in the request context
	CtxComponentKindKey CtxKey = "CTX_COMPONENT_KIND"
)

// Interceptor is the interface for a structural interceptor
type Interceptor interface {
	BeforeDispatch(ctx context.Context, req Request) context.Context
	AfterDispatch(ctx context.Context, req Request, queue ResponseQueue)
	AfterCompletion(ctx context.Context, req Request, queue ResponseQueue)
}

// NoopBeforeDispatchInterceptor does no operations before dispatch
type NoopBeforeDispatchInterceptor struct{}

// BeforeDispatch is an empty method
func (i *NoopBeforeDispatchInterceptor) BeforeDispatch(ctx context.Context, req Request) context.Context {
	return ctx
}

// NoopAfterDispatchInterceptor does no operations after dispatch
type NoopAfterDispatchInterceptor struct{}

// AfterDispatch is an empty method
func (i *NoopAfterDispatchInterceptor) AfterDispatch(context.Context, Request, ResponseQueue) {
}

// NoopAfterCompletionInterceptor does no operations after request completion
type NoopAfterCompletionInterceptor struct{}

// AfterCompletion is an empty method
func (i *NoopAfterCompletionInterceptor) AfterCompletion(context.Context, Request, ResponseQueue) {
}
