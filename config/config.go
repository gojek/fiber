package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gojek/fiber"
	"github.com/gojek/fiber/grpc"
	fiberHTTP "github.com/gojek/fiber/http"
	"github.com/gojek/fiber/protocol"
	"github.com/gojek/fiber/types"
)

// DefaultClientTimeout defines the default http client timeout to use,
// if it is not supplied in the config
const DefaultClientTimeout = time.Second

// Config is the base interface to initialise a network from a config file
type Config interface {
	initComponent() (fiber.Component, error)
}

// ComponentConfig is used to parse the base properties for a component
type ComponentConfig struct {
	ID   string `json:"id" required:"true"`
	Type string `json:"type" required:"true"`
}

// Routes represent a collection of configurations.
type Routes []Config

// UnmarshalJSON is used to parse a given input byte array to config objects
func (r Routes) UnmarshalJSON(b []byte) error {
	data := make([]json.RawMessage, 0)
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	for idx, route := range data {
		cfg, err := parseConfig(route)
		if err != nil {
			return err
		}
		r[idx] = cfg
	}
	return nil
}

// Duration is an alias for time.Duration (required since time.Duration Unmarshal is not defined)
type Duration time.Duration

// UnmarshalJSON converts the byte representation of the Duration object into time.Duration
func (d *Duration) UnmarshalJSON(b []byte) error {
	val, err := time.ParseDuration(strings.Trim(string(b), `"`))
	*d = Duration(val)
	return err
}

// MarshalJSON converts the Duration object the byte representation of its human-readable format
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// Routes takes in an object of type Routes and returns a map of each route's ID and the route
func (r Routes) Routes() (map[string]fiber.Component, error) {
	routes := make(map[string]fiber.Component)
	for _, routeConfig := range r {
		route, err := routeConfig.initComponent()
		if err != nil {
			return nil, err
		}
		routes[route.ID()] = route
	}
	return routes, nil
}

// MultiRouteConfig is used to parse the configuration for a MultiRouteComponent
type MultiRouteConfig struct {
	ComponentConfig
	Routes Routes `json:"routes" required:"true"`
}

// RouterConfig is used to parse the configuration for a Router
type RouterConfig struct {
	MultiRouteConfig
	Strategy StrategyConfig `json:"strategy" required:"true"`
}

// StrategyConfig is used to parse the configuration for a RoutingStrategy
type StrategyConfig struct {
	Type       string          `json:"type" required:"true"`
	Properties json.RawMessage `json:"properties" yaml:"properties,omitempty"`
}

// Strategy takes a reference to a StrategyConfig and creates a RoutingStrategy
func (c *StrategyConfig) Strategy() (fiber.RoutingStrategy, error) {
	strategy, err := types.StrategyByName(c.Type)
	return strategy, err
}

func (c *RouterConfig) initComponent() (fiber.Component, error) {
	var router fiber.Router
	switch c.Type {
	case "LAZY_ROUTER":
		router = fiber.NewLazyRouter(c.ID)
	case "EAGER_ROUTER":
		router = fiber.NewEagerRouter(c.ID)
	default:
		return nil, fmt.Errorf("unknown router type: [%s]", c.Type)
	}
	routes, err := c.Routes.Routes()
	if err != nil {
		return nil, err
	}
	router.SetRoutes(routes)

	strategy, err := c.Strategy.Strategy()
	if err != nil {
		return nil, err
	}

	// Initialize strategy
	err = strategy.Initialize(c.Strategy.Properties)
	if err != nil {
		return nil, err
	}
	// Set the strategy on the router
	router.SetStrategy(strategy)
	return router, nil
}

// CombinerConfig is used to parse the configuration for a Combiner
type CombinerConfig struct {
	MultiRouteConfig
	FanIn FanInConfig `json:"fan_in" required:"true"`
}

// FanInConfig is used to parse the configuration for a FanIn
type FanInConfig struct {
	Type       string          `json:"type" required:"true"`
	Properties json.RawMessage `json:"properties" yaml:"properties,omitempty"`
}

// FanIn takes a reference to a FanInConfig and creates a FanIn
func (c *FanInConfig) FanIn() (fiber.FanIn, error) {
	fanIn, err := types.FanInByName(c.Type)
	return fanIn, err
}

func (c *CombinerConfig) initComponent() (fiber.Component, error) {
	combiner := fiber.NewCombiner(c.ID)

	routes, err := c.Routes.Routes()
	if err != nil {
		return nil, err
	}
	combiner.SetRoutes(routes)

	fanIn, err := c.FanIn.FanIn()
	if err != nil {
		return nil, err
	}

	// Initialize fanIn
	err = fanIn.Initialize(c.FanIn.Properties)
	if err != nil {
		return nil, err
	}
	// Set the fanIn on the combiner
	return combiner.WithFanIn(fanIn), nil
}

// ProxyConfig is used to parse the configuration for a Proxy
type ProxyConfig struct {
	ComponentConfig
	Endpoint string            `json:"endpoint" required:"true"`
	Timeout  Duration          `json:"timeout"`
	Protocol protocol.Protocol `json:"protocol"`
	GrpcConfig
}

type GrpcConfig struct {
	ServiceMethod string `json:"service_method,omitempty"`
}

func (c *ProxyConfig) initComponent() (fiber.Component, error) {

	var dispatcher fiber.Dispatcher
	var err error
	var backend fiber.Backend
	if strings.EqualFold(string(c.Protocol), string(protocol.GRPC)) {
		dispatcher, err = grpc.NewDispatcher(grpc.DispatcherConfig{
			ServiceMethod: c.ServiceMethod,
			Endpoint:      c.Endpoint,
			Timeout:       time.Duration(c.Timeout),
		})
	} else {
		httpClient := &http.Client{Timeout: time.Duration(c.Timeout)}
		dispatcher, err = fiberHTTP.NewDispatcher(httpClient)
		backend = fiber.NewBackend(c.ID, c.Endpoint)
	}
	if err != nil {
		return nil, err
	}
	caller, err := fiber.NewCaller(c.ID, dispatcher)
	if err != nil {
		return nil, err
	}

	return fiber.NewProxy(backend, caller), nil
}

// InitComponentFromConfig takes in the path to a config file, parses the contents
// and if successful, constructs a fiber Component
func InitComponentFromConfig(configPath string) (fiber.Component, error) {
	if yamlFile, err := os.ReadFile(configPath); err != nil {
		return nil, err
	} else if cfg, err := parseConfig(yamlFile); err != nil {
		return nil, err
	} else {
		return cfg.initComponent()
	}
}

func parseConfig(data []byte) (Config, error) {
	typez := struct {
		Type   string            `json:"type" required:"true"`
		Routes []json.RawMessage `json:"routes"`
	}{}

	if err := yaml.Unmarshal(data, &typez); err != nil {
		return nil, err
	}

	var dst Config
	switch typez.Type {
	case "PROXY":
		dst = &ProxyConfig{
			// Set the default value here, can't find an easier way to supply defaults
			// Ref: https://github.com/go-yaml/yaml/issues/165
			Timeout: Duration(DefaultClientTimeout),
		}
	case "EAGER_ROUTER", "LAZY_ROUTER":
		dst = &RouterConfig{
			MultiRouteConfig: MultiRouteConfig{Routes: make(Routes, len(typez.Routes))},
		}
	case "COMBINER":
		dst = &CombinerConfig{
			MultiRouteConfig: MultiRouteConfig{Routes: make(Routes, len(typez.Routes))},
		}
	default:
		return nil, fmt.Errorf("unknown component type: %s", typez.Type)
	}

	if err := yaml.Unmarshal(data, dst); err != nil {
		return nil, err
	}

	return dst, nil
}
