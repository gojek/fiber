package interceptor

import (
	"context"

	"github.com/gojek/fiber"

	"go.uber.org/zap"
)

// NewLoggingInterceptor is a creator factory for a ResponseLoggingInterceptor
func NewLoggingInterceptor(log *zap.SugaredLogger) fiber.Interceptor {
	return &ResponseLoggingInterceptor{
		logger: log,
	}
}

// ResponseLoggingInterceptor is the structural interceptor used for logging responses
type ResponseLoggingInterceptor struct {
	fiber.NoopBeforeDispatchInterceptor
	fiber.NoopAfterCompletionInterceptor
	logger *zap.SugaredLogger
}

// AfterDispatch logs the success or failure information of a request
func (i *ResponseLoggingInterceptor) AfterDispatch(ctx context.Context, req fiber.Request, queue fiber.ResponseQueue) {
	for resp := range queue.Iter() {
		if resp.IsSuccess() {
			i.logger.Infof("%s: %s", resp.BackendName(), resp.Payload())
		} else {
			i.logger.Warnf("%s: %s", resp.BackendName(), resp.Payload())
		}
	}
}
