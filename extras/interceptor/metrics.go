package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/gojek/fiber"
)

// StatsdClient is an interface for a stats listener
type StatsdClient interface {
	Increment(string)
	Unique(string, string)
	Timing(string, interface{})
}

// MetricsKey is an alias for a string
type MetricsKey string

// CtxDispatchStartTimeKey is used to record the start time of a request
var (
	CtxDispatchStartTimeKey MetricsKey = "CTX_DISPATCH_START_TIME"
)

// NewMetricsInterceptor creates a new MetricsInterceptor with the given client
func NewMetricsInterceptor(client StatsdClient) fiber.Interceptor {
	return &MetricsInterceptor{
		statsd: client,
	}
}

// MetricsInterceptor is an interceptor to log metrics
type MetricsInterceptor struct {
	fiber.NoopAfterDispatchInterceptor
	statsd StatsdClient
}

func (i *MetricsInterceptor) operationName(ctx context.Context, req fiber.Request, suffix string) string {
	componentID := ctx.Value(fiber.CtxComponentIDKey)
	return fmt.Sprintf("fiber.%s.%s", componentID, suffix)
}

// BeforeDispatch records count of the requests through a component and the request start time
func (i *MetricsInterceptor) BeforeDispatch(ctx context.Context, req fiber.Request) context.Context {
	i.statsd.Increment(i.operationName(ctx, req, "count"))

	ctx = context.WithValue(ctx, CtxDispatchStartTimeKey, time.Now())
	return ctx
}

// AfterCompletion computes and submits the request completion time
func (i *MetricsInterceptor) AfterCompletion(ctx context.Context, req fiber.Request, queue fiber.ResponseQueue) {
	if startTime, ok := ctx.Value(CtxDispatchStartTimeKey).(time.Time); ok {
		i.statsd.Timing(i.operationName(ctx, req, "timing"), int(time.Since(startTime)/time.Millisecond))
	}
}
