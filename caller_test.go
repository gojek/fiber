package fiber_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/gojek/fiber"
	testutils "github.com/gojek/fiber/internal/testutils/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDispatcher struct {
	mock.Mock
}

func (h *MockDispatcher) Do(req fiber.Request) fiber.Response {
	args := h.Called(req)
	if resp := args.Get(0); resp != nil {
		return resp.(fiber.Response)
	}
	return nil
}

func TestCaller_Dispatch(t *testing.T) {
	payload := "**BODY**"
	expectedResponse := testutils.MockResp(http.StatusOK, payload, nil, nil)

	dispatcher := new(MockDispatcher)
	dispatcher.On("Do", mock.Anything).Return(expectedResponse)

	caller, _ := fiber.NewCaller("", dispatcher)

	req := testutils.MockReq("GET", "http://:8080/test", "")
	resp := <-caller.Dispatch(context.Background(), req).Iter()

	assert.NotNil(t, resp)
	assert.Empty(t, resp.BackendName())
	assert.True(t, resp.IsSuccess())
	assert.Equal(t, []byte(payload), resp.Payload())

	dispatcher.AssertExpectations(t)
}
