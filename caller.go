package fiber

import (
	"context"
	"errors"

	"github.com/gojek/fiber/util"
)

// Caller is the basic network component, that dispatches incoming request
// using configured transport-agnostic Dispatcher and asynchronously
// sends the response into an output channel
type Caller struct {
	BaseComponent
	dispatcher Dispatcher
}

// NewCaller is a factory method that creates a new instance of Caller
// with given id and provided Dispatcher
func NewCaller(id string, dispatcher Dispatcher) (*Caller, error) {
	if id == "" {
		id = "caller_" + util.UID()
	}

	if dispatcher == nil {
		return nil, errors.New("request dispatcher can not be nil")
	}

	return &Caller{
		BaseComponent: BaseComponent{id: id, kind: CallerKind},
		dispatcher:    dispatcher,
	}, nil
}

// Dispatch uses Dispatcher to process incoming request and asynchronously sends
// received response into the output channel. The output channel will be closed
// after Dispatcher has processed request and response was sent back
func (c *Caller) Dispatch(ctx context.Context, req Request) ResponseQueue {
	ctx = c.beforeDispatch(ctx, req)
	out := make(chan Response, 1)
	queue := NewResponseQueue(out, 1)
	defer c.afterDispatch(ctx, req, queue)

	go func() {
		defer c.afterCompletion(ctx, req, queue)
		out <- c.dispatcher.Do(req)
		close(out)
	}()
	return queue
}
