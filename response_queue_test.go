package fiber_test

import (
	"testing"
	"time"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/internal/testutils/http"
	"github.com/stretchr/testify/assert"
)

func chanToArray(ch <-chan fiber.Response) []fiber.Response {
	var responses []fiber.Response
	for r := range ch {
		responses = append(responses, r)
	}
	return responses
}

func makeChan(r ...fiber.Response) chan fiber.Response {
	out := make(chan fiber.Response)

	go func() {
		for _, resp := range r {
			out <- resp
		}
		close(out)
	}()

	return out
}
func TestNewResponseQueue(t *testing.T) {
	responses := []fiber.Response{
		http.MockResp(200, "foo", nil, nil),
		http.MockResp(200, "bar", nil, nil),
	}

	q := fiber.NewResponseQueue(makeChan(responses...), 0)

	assert.Equal(t, responses, chanToArray(q.Iter()))
	assert.Equal(t, responses, chanToArray(q.Iter()))
}

func TestNewResponseQueueFromResponses(t *testing.T) {
	responses := []fiber.Response{
		http.MockResp(200, "foo", nil, nil),
		http.MockResp(200, "bar", nil, nil),
	}

	q := fiber.NewResponseQueueFromResponses(responses...)

	assert.Equal(t, responses, chanToArray(q.Iter()))
	assert.Equal(t, responses, chanToArray(q.Iter()))

	q = fiber.NewResponseQueueFromResponses()
	assert.Empty(t, chanToArray(q.Iter()))
}

func TestResponseQueue_Iter(t *testing.T) {
	responses := []fiber.Response{
		http.MockResp(200, "fist", nil, nil),
		http.MockResp(200, "second", nil, nil),
		http.MockResp(200, "third", nil, nil),
		http.MockResp(200, "fourth", nil, nil),
	}

	out := make(chan fiber.Response)

	go func() {
		for i, r := range responses {
			time.Sleep(time.Duration(i*50) * time.Millisecond)
			out <- r
		}
		close(out)
	}()

	q := fiber.NewResponseQueue(out, 0)

	assert.Equal(t, responses, chanToArray(q.Iter()))
	time.Sleep(400 * time.Millisecond)
	assert.Equal(t, responses, chanToArray(q.Iter()))
}
