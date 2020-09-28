package interceptor

import (
	"context"
	"fmt"

	"github.com/gojek/fiber"

	"github.com/opentracing/opentracing-go"
)

// NewTracingInterceptor creates a TracingInterceptor
func NewTracingInterceptor(tracer opentracing.Tracer) fiber.Interceptor {
	return &TracingInterceptor{
		tracer: tracer,
	}
}

// TracingInterceptor allows for tracing requests
type TracingInterceptor struct {
	fiber.NoopAfterDispatchInterceptor
	tracer opentracing.Tracer
}

func (i *TracingInterceptor) operationName(ctx context.Context, req fiber.Request) string {
	componentID := ctx.Value(fiber.CtxComponentIDKey)
	return fmt.Sprintf("[%s] %s ", componentID, req.OperationName())
}

// BeforeDispatch starts and returns a span with the given operation name
func (i *TracingInterceptor) BeforeDispatch(ctx context.Context, req fiber.Request) context.Context {
	_, ctx = opentracing.StartSpanFromContextWithTracer(ctx, i.tracer, i.operationName(ctx, req))
	return ctx
}

// AfterCompletion returns the Span previously associated with context
func (i *TracingInterceptor) AfterCompletion(ctx context.Context, req fiber.Request, queue fiber.ResponseQueue) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.Finish()
	}
}
