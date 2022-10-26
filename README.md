# fiber 

**fiber** is a Go library for building dynamic proxies, routers and 
traffic mixers from a set of composable abstract network components. 

Core components of fiber are transport agnostic, however, there is 
Go's `net/http`-based implementation provided in [fiber/http](http) package 
and a grpc implementation in [fiber/grpc](grpc).

The grpc implementation will use the byte payload from the request and response using a custom codec to minimize marshaling overhead.
It is expected that the client [marshal](https://pkg.go.dev/github.com/golang/protobuf/proto#Marshal) the message and unmarshall into the intended proto response.

## Usage

```go
import (
    "github.com/gojek/fiber"                  // fiber core
    "github.com/gojek/fiber/config"           // fiber config 
    fiberhttp "github.com/gojek/fiber/http"   // fiber http if required
    fibergrpc "github.com/gojek/fiber/grpc"   // fiber grpc if required
)
```

Define your fiber component in YAML config. For example:

**fiber.yaml for http:**
```yaml
type: EAGER_ROUTER
id: eager_router
strategy:
  type: fiber.RandomRoutingStrategy
routes:
  - id: route_a
    type: PROXY
    timeout: "20s"
    endpoint: "http://localhost:8080/routes/route-a"
  - id: route_b
    type: PROXY
    timeout: "40s"
    endpoint: "http://localhost:8080/routes/route-b"
```

**fiber.yaml for grpc:**
```yaml
type: EAGER_ROUTER
id: eager_router
strategy:
  type: fiber.RandomRoutingStrategy
routes:
  - id: route_a
    type: PROXY
    timeout: "20s"
    endpoint: "localhost:50555" 
    service: "mypackage.Greeter"
    method: "SayHello"
    protocol: "grpc"
  - id: route_b
    type: PROXY
    timeout: "40s"
    endpoint: "localhost:50555"
    service: "mypackage.Greeter"
    method: "SayHello"
    protocol: "grpc"
```

Construct new fiber component from the config:

**main.go:**
```go
import "github.com/gojek/fiber/config"

compomnent, err := config.FromConfig("./fiber.yaml")
```

Start serving http requests:

**main.go:**
```go
import (
    fiberhttp "github.com/gojek/fiber/http"
)

options := fiberhttp.Options{
    Timeout: 20 * time.Second,
}

fiberHandler := fiberhttp.NewHandler(component, options)

http.ListenAndServe(":8080", fiberHandler)
```

It is also possible to define fiber component programmatically, using fiber API.
For example:

```go
import (
    "github.com/gojek/fiber"
    fiberhttp "github.com/gojek/fiber/http"
)

component := fiber.NewEagerRouter("eager-router")
component.SetStrategy(new(extras.RandomRoutingStrategy))

httpDispatcher, _ := fiberhttp.NewDispatcher(http.DefaultClient)
caller, _ := fiber.NewCaller("", httpDispatcher)

component.SetRoutes(map[string]fiber.Component{
    "route-a": fiber.NewProxy(
        fiber.NewBackend("route-a", "http://localhost:8080/routes/route-a"),
        caller),
    "route-b": fiber.NewProxy(
        fiber.NewBackend("route-b", "http://localhost:8080/routes/route-b"),
        caller),
})
``` 

For more sample code snippets and grpc usage, head over to the [example](./example) directory.

## Concepts

There are few general abstractions used in fiber:

- [Component](component.go) - a fiber abstraction, which has a goal to take an incoming request,
dispatch it, and return a queue with zero or more responses in it. There is also a special kind 
of components – [MultiRouteComponent](multi_route_component.go). The family of multi-route components
contains higher-order components that have one or more children components registered as their routes.  
There are few basic components implemented in fiber, such as [Proxy](proxy_caller.go), [FanOut](fan_out.go),
[Combiner](combiner.go), [Router](router.go) ([EagerRouter](eager_router.go) and [LazyRouter](lazy_router.go)).  
Check [Components](#Components) for more information about existing fiber components.

- [Fan In](fan_in.go) – is an aggregator, that takes a response queue and returns a single response.
It can be one of the responses from the queue, a combination of responses or just a 
new arbitrary response. `Fan In` is a necessary part of [Combiner](combiner.go) implementation. 

- [Routing Strategy](routing_strategy.go) – is used together with Routers (either lazy or eager) and 
is responsible for defining the order (priority) of routes registered with given router. So, for every incoming
request, `Routing Strategy` should tell which route should dispatch this given request and what is the 
order of fallback routes to be used in case the primary route has failed to successfully dispatch the request.
Routing strategy is generally expected to be implemented by the client application, because it might be 
domain specific. However, the simplest possible implementation of routing strategy is provided as a reference in 
[RandomRoutingStrategy](extras/random_routing_strategy.go).

- [Interceptor](interceptor.go) – fiber supports pluggable interceptors to examine request and responses.
Interceptors are useful for implementing various req/response loggers, metrics or distributed traces collectors.
fiber comes with few pre-implemented interceptors ([extras/interceptor](extras/interceptor)) to serve
above-mentioned purposes.  

## Components

fiber allows defining request processing rules by composing basic fiber components
in more complicated execution graphs. There are few standard fiber components implemented in this library:

- `PROXY` – component, that dispatches incoming request against configured proxy backend url.                   
Configuration:               
    - `id` – component ID. Example `my_proxy`
    - `endpoint` - proxy endpoint url. Example for http `http://your-proxy:8080/nested/path` or  grpc `127.0.0.1:50050`
    - `timeout` - request timeout for dispatching a request. Example `100ms` 
    - `protocol` - communication protocol. Only "grpc" or "http" supported.
    - `service` - for grpc only, package name and service name. Example `fiber.Greeter` 
    - `method` - for grpc only, method name of the grpc service to invoke. Example `SayHello`
    
- `FAN_OUT` - component, that dispatches incoming request by sending it to each of its registered 
`routes`. Response queue will contain responses of each route in order they have arrived.  
Configuration:
    - `id` - component ID. Example `fan_out`
    - `routes` – list of fiber components definitions. `Fan Out` will send incoming request to each of its
    route components and collect responses into a queue. 
    
- `COMBINER` - dispatches incoming request by sending it to each of its registered `routes` and 
then aggregating received responses into a single response by using provided `fan_in`.  
Configuration:     
    - `id` - component ID
    - `fan_in` - configuration of the [FanIn](fan_in.go), that will be used in this combiner
       - `type` - registered type name of the fan in. Example: `fiber.FastestResponseFanIn`. 
       (See also [Custom Types](#Custom Types))
       - `properties` - arbitrary yaml configuration that would be passed to the FanIn's 
       `Initialize` method during the component initialization
    - `routes` - list of fiber component definitions that would be registered as this combiner's routes.

- `EAGER_ROUTER` - dispatches incoming request by sending it simultaneously to each registered route and
then returning either a response from the primary route (defined by the routing strategy) or switches 
back to one of the fallback routes. Eager routers are useful in situations, when it's crucial to return
fallback response with a minimal delay in case if primary route failed to respond with successful response.
Configuration:   
    - `id` – component ID
    - `strategy` - configuration of the [RoutingStrategy](routing_strategy.go), that would be used 
    with this router
        - `type` - registered type name of the routing strategy. Example: `fiber.RandomRoutingStrategy`. 
        (See also [Custom Types](#Custom Types))
        - `properties` - arbitrary yaml configuration that would be passed to the RoutingStrategy's 
        `Initialize` method during the component initialization
    - `routes` - list of fiber components definitions that would be registered as this router routes.
    
- `LAZY_ROUTER` - dispatches incoming request by retrieving information about the primary and fallback routes
order from its [RoutingStrategy](routing_strategy.go) and then sending the request to the routes in defined order
until one of the routes responds with successful response.
Configuration:   
    - `id` – component ID
    - `strategy` - configuration of the [RoutingStrategy](routing_strategy.go), that would be used 
    with this router
        - `type` - registered type name of the routing strategy. Example: `fiber.RandomRoutingStrategy`
        - `properties` - arbitrary yaml configuration that would be passed to the RoutingStrategy's 
        `Initialize` method during the component initialization
    - `routes` - list of fiber components definitions that would be registered as this router routes.
    
## Interceptors

fiber comes with few pre-defined interceptors, that are serving the most common use-cases:

- [ResponseLoggingInterceptor](extras/interceptor/logging.go) - subscribes to responses from the response queue 
and uses an instance of `zap.SugaredLogger` to log the response's payload.
 
- [MetricsInterceptor](extras/interceptor/metrics.go) - collects the `count` and `time` metrics of component's 
`Dispatch` method and forwards these time-series data using provided `statsd` client. 

- [TracingInterceptor](extras/interceptor/tracing.go) - uses 
[opentracing/opentracing-go](https://github.com/opentracing/opentracing-go) client to create spans of the `Dispatch`
method execution

### Using interceptors

It's also possible to create a custom interceptor by implementing `fiber.Interceptor` interface:

```go
type Interceptor interface {
	BeforeDispatch(ctx context.Context, req Request) context.Context

	AfterDispatch(ctx context.Context, req Request, queue ResponseQueue)

	AfterCompletion(ctx context.Context, req Request, queue ResponseQueue)
}
```

Then, one or more interceptors can be attached to the fiber component by calling `AddInterceptor` method:

```go
import (
    "github.com/gojek/fiber/config"
    "github.com/gojek/fiber/extras/interceptor"
)

compomnent, err := config.FromConfig("./fiber.yaml")

statsdClient := // initialize statsd client
zapLog :=       // initialize logger

component.AddInterceptor(
    true,  // add interceptors recursively to children components
    interceptor.NewMetricsInterceptor(statsdClient),
    interceptor.NewLoggingInterceptor(zapLog),
)
```

## Custom Types

It is also possible to register a custom `RoutingStrategy` or `FanIn` implementation in `fiber`'s type system.
 
First, create your own RoutingStrategy. For example, let's define a routing strategy, that directs requests
to `route-a` in case if session ID (passed via Header) is odd and to `route-b` if it is even:
```go
package mypackage

import (
	"context"
	"strconv"
)

type OddEvenRoutingStrategy struct {}

func (s *OddEvenRoutingStrategy) SelectRoute(
	ctx context.Context,
	req fiber.Request,
	routes map[string]fiber.Component,
) (fiber.Component, []fiber.Component, fiber.Attributes, error) {
	sessionIdStr := ""
	if sessionHeader, ok := req.Header()["X-Session-ID"]; ok {
		sessionIdStr = sessionHeader[0]
	}
    // Metadata that can be propagated upstream for logging / debugging
    attr := fiber.NewAttributesMap()
	
	if sessionID, err := strconv.Atoi(sessionIdStr); err != nil {
		return nil, nil, err
	} else {
		if sessionID % 2 != 0 {
			return routes["route-a"], []fiber.Component{}, attr.WithAttribute("Match-Type", "even"), nil
		} else {
			return routes["route-b"], []fiber.Component{}, attr.WithAttribute("Match-Type", "odd"), nil
		}
	}
}
```

Then, register this routing strategy in fiber's type system:

```go
package main

import (
	"github.com/gojek/fiber/types"
)

func main() {
    if err := types.InstallType(
        "mypackage.OddEvenRoutingStrategy", 
        &mypackage.OddEvenRoutingStrategy{}); err != nil {
        panic(err)
    }
    // ... 
}
```

So now, `mypackage.OddEvenRoutingStrategy` is registered and can be used in fiber component configuration:

```yaml
type: LAZY_ROUTER
id: lazy-router
routes:
  - id: "route-a"
    type: PROXY
    endpoint: "http://www.mocky.io/v2/5e4caccc310000e2cad8c071"
    timeout: 5s
  - id: "route-b"
    type: PROXY
    endpoint: "http://www.mocky.io/v2/5e4cacd4310000e1cad8c073"
    timeout: 5s
strategy:
  type: mypackage.OddEvenRoutingStrategy
```

## Licensing

[Apache 2.0 License](./LICENSE)