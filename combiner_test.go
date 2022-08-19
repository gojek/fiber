package fiber_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gojek/fiber"
	testUtilsHttp "github.com/gojek/fiber/internal/testutils/http"
	"github.com/gojek/fiber/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockFanOut struct {
	mock.Mock
	*fiber.BaseMultiRouteComponent
}

func (m *mockFanOut) Dispatch(ctx context.Context, req fiber.Request) fiber.ResponseQueue {
	args := m.Called(ctx, req)
	if queue := args.Get(0); queue != nil {
		return queue.(fiber.ResponseQueue)
	}
	return nil
}

type mockFanIn struct {
	mock.Mock
	*fiber.BaseFanIn
}

func (fanIn *mockFanIn) Aggregate(ctx context.Context, req fiber.Request, queue fiber.ResponseQueue) fiber.Response {
	args := fanIn.Called(ctx, req, queue)
	return args.Get(0).(fiber.Response)
}

type combinerTestCase struct {
	name      string
	request   fiber.Request
	responses fiber.ResponseQueue
	expected  fiber.Response
}

func (tt *combinerTestCase) MockFanOut() *mockFanOut {
	fanOut := &mockFanOut{}
	fanOut.On("Dispatch", mock.Anything, tt.request).Once().Return(tt.responses)
	return fanOut
}

func (tt *combinerTestCase) MockFanIn() *mockFanIn {
	fanIn := &mockFanIn{}
	fanIn.On("Aggregate", mock.Anything, tt.request, tt.responses).Once().Return(tt.expected)
	return fanIn
}

func TestCombiner_Dispatch(t *testing.T) {
	timeout := 200 * time.Millisecond
	suite := []combinerTestCase{
		{
			name:     "two routes/two OK responseQueue",
			request:  testUtilsHttp.MockReq("POST", "http:/combiner:8080", ""),
			expected: testUtilsHttp.MockResp(200, "A-OK,B-OK", nil, nil),
		},
	}

	for _, tt := range suite {
		fanOut := tt.MockFanOut()
		fanIn := tt.MockFanIn()

		combiner := fiber.NewCombiner("")
		combiner.FanOut = fanOut
		combiner.WithFanIn(fanIn)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		received := make([]fiber.Response, 0)
		for responsesCh := combiner.Dispatch(ctx, tt.request).Iter(); ; {
			select {
			case resp, ok := <-responsesCh:
				if ok {
					received = append(received, resp)
					continue
				}
			case <-time.After(timeout + 5*time.Millisecond):
				assert.Fail(t, fmt.Sprintf("[%s] failed: it didn't terminate after a timeout...", tt.name))
			}

			cancel()
			break
		}

		assert.Equal(t, 1, len(received))
		assert.Equal(t, tt.expected, received[0])

		fanOut.AssertExpectations(t)
		fanIn.AssertExpectations(t)
	}
}

func TestCombiner_Id(t *testing.T) {
	id := util.UID()

	combiner := fiber.NewCombiner(id)
	assert.Equal(t, id, combiner.ID())
}
