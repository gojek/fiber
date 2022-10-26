package extras

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/gojek/fiber"
)

// RandomRoutingStrategy is just a reference implementation of a RoutingStrategy.
// It randomly selects a primary route and all other routes are fallbacks (with no specific order)
type RandomRoutingStrategy struct {
	fiber.BaseFiberType
}

// SelectRoute on the RandomRoutingStrategy selects one of the given routes as the primary
// route, at random, and sets the others as fallbacks
func (s *RandomRoutingStrategy) SelectRoute(
	_ context.Context,
	_ fiber.Request,
	routes map[string]fiber.Component,
) (route fiber.Component, fallbacks []fiber.Component, labels fiber.Labels, err error) {
	idx := rand.Intn(len(routes))
	// Add idx to attribute map for logging / debugging upstream
	labels = fiber.NewLabelsMap().WithLabel("idx", strconv.Itoa(idx))

	for _, child := range routes {
		if idx == 0 {
			route = child
		} else {
			fallbacks = append(fallbacks, child)
		}
		idx--
	}
	return route, fallbacks, labels, nil
}
