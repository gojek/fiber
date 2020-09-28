package testutils

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/gojek/fiber"
)

type MockComponent struct {
	*fiber.BaseComponent
	mock.Mock

	responses []DelayedResponse
}

func NewMockComponent(id string, responses ...DelayedResponse) *MockComponent {
	return &MockComponent{
		BaseComponent: fiber.NewBaseComponent(id, ""),
		responses:     responses,
	}
}

func (m *MockComponent) Dispatch(context.Context, fiber.Request) fiber.ResponseQueue {
	out := make(chan fiber.Response, len(m.responses))

	go func() {
		for _, r := range m.responses {
			time.Sleep(r.Latency)
			out <- r.Response
		}
		close(out)
	}()

	return fiber.NewResponseQueue(out, len(m.responses))
}
