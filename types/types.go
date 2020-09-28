package types

import (
	"fmt"
	"reflect"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/extras"
)

// Category is an alias for a string
type Category string

const (
	// RoutingStrategy is the type name use to define a RoutingStrategy
	RoutingStrategy = "ROUTING_STRATEGY"
	// FanIn is the type name use to define a FAN_IN
	FanIn = "FAN_IN"
)

var categories = map[Category]reflect.Type{
	RoutingStrategy: reflect.TypeOf((*fiber.RoutingStrategy)(nil)).Elem(),
	FanIn:           reflect.TypeOf((*fiber.FanIn)(nil)).Elem(),
}

var types = map[Category]map[string]reflect.Type{
	RoutingStrategy: {
		"fiber.RandomRoutingStrategy": reflect.TypeOf(&extras.RandomRoutingStrategy{}).Elem(),
	},
	FanIn: {
		"fiber.FastestResponseFanIn": reflect.TypeOf(&extras.FastestResponseFanIn{}).Elem(),
	},
}

func typeByName(category Category, typez string) (interface{}, error) {
	if strategyType, exist := types[category][typez]; exist {
		return reflect.New(strategyType).Interface(), nil
	}
	return nil, fmt.Errorf("unknown %s type: %s", category, typez)
}

// InstallType updates the type map with new sub-types
func InstallType(key string, objectOfType interface{}) error {
	newType := reflect.TypeOf(objectOfType)
	i := 0
	for k, category := range categories {
		if newType.Implements(category) {
			i++
			types[k][key] = newType.Elem()
		}
	}
	if i == 0 {
		return fmt.Errorf("type %s is not compatible with any category", newType.Name())
	}
	return nil
}

// StrategyByName identifies a routing strategy type that matches the type name specified
// and returns an instance of that type
func StrategyByName(name string) (fiber.RoutingStrategy, error) {
	if strategy, err := typeByName(RoutingStrategy, name); err != nil {
		return nil, err
	} else if typed, ok := strategy.(fiber.RoutingStrategy); ok {
		return typed, nil
	}
	return nil, fmt.Errorf("incompatible strategy type: %s", name)
}

// FanInByName identifies a fan in type that matches the type name specified
// and returns an instance of that type
func FanInByName(name string) (fiber.FanIn, error) {
	fanIn, err := typeByName(FanIn, name)
	if err != nil {
		return nil, err
	}
	if typed, ok := fanIn.(fiber.FanIn); ok {
		return typed, nil
	}
	return nil, fmt.Errorf("incompatible fan-in type: %s", name)
}
