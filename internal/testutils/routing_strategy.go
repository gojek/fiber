package testutils

import (
	"context"
	"time"

	"github.com/gojek/fiber"
	"github.com/stretchr/testify/mock"
)

type MockRoutingStrategy struct {
	mock.Mock
	fiber.BaseFiberType
}

func (s *MockRoutingStrategy) SelectRoute(
	_ context.Context,
	req fiber.Request,
	routes map[string]fiber.Component,
) (route fiber.Component, fallbacks []fiber.Component, attr fiber.Attributes, err error) {
	args := s.Called(req, routes)

	if args.Get(0) == nil {
		return (fiber.Component)(nil), args.Get(1).([]fiber.Component), args.Get(2).(fiber.Attributes), args.Error(3)
	}

	return args.Get(0).(fiber.Component), args.Get(1).([]fiber.Component), args.Get(2).(fiber.Attributes), args.Error(3)
}

func NewMockRoutingStrategy(
	routes map[string]fiber.Component,
	order []string,
	latency time.Duration,
	err error,
) *MockRoutingStrategy {
	strategy := new(MockRoutingStrategy)
	strategy.On("SelectRoute", mock.Anything, routes).
		Run(func(args mock.Arguments) {
			time.Sleep(latency)
		}).
		Return(
			func() (fiber.Component, []fiber.Component, fiber.Attributes, error) {
				attr := fiber.NewAttributesMap().WithAttribute("order", order...)
				if len(order) == 0 {
					return (fiber.Component)(nil), []fiber.Component{}, attr, err
				}
				// Else
				fallbacks := make([]fiber.Component, 0)
				for i := 1; i < len(order); i++ {
					fallbacks = append(fallbacks, routes[order[i]])
				}

				return routes[order[0]], fallbacks, attr, err
			}(),
		)
	return strategy
}
