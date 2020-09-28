package fiber

import (
	"context"
	"sync"

	"github.com/gojek/fiber/util"
)

// FanOut is the base interface for structural FanOut Components, that is
// used to dispatch the incoming request simultaneously and asynchronously
// across all configured routes, writing the responseQueue, as available, to a
// single results channel.
type FanOut interface {
	MultiRouteComponent
}

// BaseFanOut is a component, that dispatches incoming request by each of its nested sub-routes
type BaseFanOut struct {
	*BaseMultiRouteComponent
}

// NewFanOut initializes a new BaseFanOut component and assigns to it a generated unique ID
// and the list of nested children components
func NewFanOut(id string) *BaseFanOut {
	if id == "" {
		id = "fan-out_" + util.UID()
	}
	return &BaseFanOut{
		BaseMultiRouteComponent: NewMultiRouteComponent(id),
	}
}

// Dispatch creates a copy of incoming request (one for each sub-route), asynchronously dispatches
// these request by its children components and then merges response channels into a
// single response channel with zero or more responseQueue in it
func (fanOut *BaseFanOut) Dispatch(ctx context.Context, req Request) ResponseQueue {
	ctx = fanOut.beforeDispatch(ctx, req)
	out := make(chan Response, len(fanOut.routes))

	queue := NewResponseQueue(out, len(fanOut.routes))
	defer fanOut.afterDispatch(ctx, req, queue)

	go func() {
		defer fanOut.afterCompletion(ctx, req, queue)

		var wg sync.WaitGroup
		wg.Add(len(fanOut.routes))

		for _, route := range fanOut.routes {
			go func(route Component) {
				// Make a copy of incoming request for each sub-name
				copyReq, _ := req.Clone()

				in := route.Dispatch(ctx, copyReq).Iter()

				for {
					select {
					case resp, ok := <-in:
						if ok {
							out <- resp.WithBackendName(route.ID())
							continue
						}
					case <-ctx.Done():
					}
					break
				}
				wg.Done()
			}(route)
		}
		wg.Wait()
		close(out)
	}()

	return queue
}
