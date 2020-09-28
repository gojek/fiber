package fiber

import "sync"

type ResponseQueue interface {
	Iter() <-chan Response
}

type responseQueue struct {
	lock          sync.RWMutex
	items         []Response
	buffer        int
	subscriptions []chan Response

	done chan struct{}
}

func (r *responseQueue) append(resp Response) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.items = append(r.items, resp)
	for _, subscription := range r.subscriptions {
		subscription <- resp
	}
}

func (r *responseQueue) Iter() <-chan Response {
	out := make(chan Response, r.buffer)

	go func() {
		r.lock.Lock()
		r.subscriptions = append(r.subscriptions, out)
		for _, resp := range r.items {
			out <- resp
		}
		r.lock.Unlock()

		<-r.done
		close(out)
	}()
	return out
}

// NewResponseQueue takes an input channel and creates a Queue with all responseQueue from it
func NewResponseQueue(in <-chan Response, bufferSize int) ResponseQueue {
	queue := &responseQueue{
		buffer: bufferSize,
		done:   make(chan struct{}),
	}

	go func(q *responseQueue) {
		defer close(q.done)

		for resp := range in {
			q.append(resp)
		}
	}(queue)
	return queue
}

// NewResponseQueueFromResponses takes list of responses and constructs
// an instance of ResponseQueue from them
func NewResponseQueueFromResponses(responses ...Response) ResponseQueue {
	queue := &responseQueue{
		items:  responses,
		buffer: len(responses),
		done:   make(chan struct{}),
	}

	close(queue.done)
	return queue
}
