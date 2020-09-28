package fiber

import "context"

// FanIn is the base interface for structural FanIn components,
// that is supposed to listen for zero or more incoming responseQueue from `in`
// and aggregates them into a single Response
type FanIn interface {
	Type
	Aggregate(ctx context.Context, req Request, queue ResponseQueue) Response
}

// BaseFanIn provides the default implementation for the Type interface
type BaseFanIn struct {
	BaseFiberType
}
