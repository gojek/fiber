package extras

import (
	"context"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/errors"
)

// FastestResponseFanIn is a FanIn that selects a first (fastest) OK response from the channel of responseQueue
type FastestResponseFanIn struct {
	fiber.BaseFanIn
}

// Aggregate returns the first (fastest) response from the result channel
func (r *FastestResponseFanIn) Aggregate(
	_ context.Context,
	_ fiber.Request,
	queue fiber.ResponseQueue,
) fiber.Response {
	out := make(chan fiber.Response, 1)
	go func() {
		defer close(out)

		for resp := range queue.Iter() {
			if resp.IsSuccess() {
				out <- resp
				return
			}
		}
		out <- fiber.NewErrorResponse(errors.ErrServiceUnavailable)
	}()
	return <-out
}
