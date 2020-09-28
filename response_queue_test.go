package fiber_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/gojek/fiber"
	"github.com/gojek/fiber/internal/testutils"
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
		testutils.MockResp(200, "foo", nil, nil),
		testutils.MockResp(200, "bar", nil, nil),
	}

	q := fiber.NewResponseQueue(makeChan(responses...), 0)

	assert.Equal(t, responses, chanToArray(q.Iter()))
	assert.Equal(t, responses, chanToArray(q.Iter()))
}

func TestNewResponseQueueFromResponses(t *testing.T) {
	responses := []fiber.Response{
		testutils.MockResp(200, "foo", nil, nil),
		testutils.MockResp(200, "bar", nil, nil),
	}

	q := fiber.NewResponseQueueFromResponses(responses...)

	assert.Equal(t, responses, chanToArray(q.Iter()))
	assert.Equal(t, responses, chanToArray(q.Iter()))

	q = fiber.NewResponseQueueFromResponses()
	assert.Empty(t, chanToArray(q.Iter()))
}

func TestResponseQueue_Iter(t *testing.T) {
	responses := []fiber.Response{
		testutils.MockResp(200, "fist", nil, nil),
		testutils.MockResp(200, "second", nil, nil),
		testutils.MockResp(200, "third", nil, nil),
		testutils.MockResp(200, "fourth", nil, nil),
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
