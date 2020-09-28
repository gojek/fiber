package fiber

import (
	"context"

	"github.com/gojek/fiber/util"
)

// Combiner is a network component, that uses BaseFanOut to dispatch incoming request
// by all of its sub-routes and then merge all the responseQueue from them into a single
// response, using provided FanIn
type Combiner struct {
	BaseComponent
	FanOut

	fanIn FanIn
}

// NewCombiner is a factory for the Combiner type.
func NewCombiner(id string) *Combiner {
	if id == "" {
		id = "combiner_" + util.UID()
	}

	return &Combiner{
		BaseComponent: BaseComponent{id: id, kind: CombinerKind},
		FanOut:        NewFanOut("fan_out"),
	}
}

// ID is the getter for the combiner's ID
func (c *Combiner) ID() string {
	return c.BaseComponent.ID()
}

// Kind is the getter for the combiner's type
func (c *Combiner) Kind() ComponentKind {
	return c.BaseComponent.kind
}

// WithFanIn is a Setter for the FanIn (aggregation strategy) on the given Combiner
func (c *Combiner) WithFanIn(fanIn FanIn) *Combiner {
	c.fanIn = fanIn
	return c
}

// Dispatch method on the Combiner will ask its embedded dispatcher to simultaneously
// dispatch the incoming request by all of its nested components. After that, Combiner's FanIn
// listens to responseQueue and aggregate them into a single response, that is being sent to output
func (c *Combiner) Dispatch(ctx context.Context, req Request) ResponseQueue {
	ctx = c.beforeDispatch(ctx, req)
	out := make(chan Response, 1)

	queue := NewResponseQueue(out, 1)
	defer c.afterDispatch(ctx, req, queue)

	go func() {
		defer c.afterCompletion(ctx, req, queue)

		out <- c.fanIn.Aggregate(ctx, req, c.FanOut.Dispatch(ctx, req))
		close(out)
	}()

	return queue
}

// AddInterceptor can be used to add the given interceptor to the Combiner and optionally,
// to all its nested components.
func (c *Combiner) AddInterceptor(recursive bool, interceptor ...Interceptor) {
	if recursive {
		c.FanOut.AddInterceptor(recursive, interceptor...)
	}
	c.BaseComponent.AddInterceptor(recursive, interceptor...)
}
