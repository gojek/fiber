## Simple Example 

Basic example shows how to define a fiber component programmatically 
using fiber API and start serving traffic with `net/http` handler.

### Define fiber component

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

### Start serving traffic

```go
package main

import (
    "net/http"
    "time"
   
    fiberhttp "github.com/gojek/fiber/http"
)

func main() {
    
    component := // defined above 

    fiberHandler := fiberhttp.NewHandler(component, fiberhttp.Options{
        Timeout: 20 * time.Second,
    })

    http.ListenAndServe(":8080", fiberHandler)
}
```

NOTE: This is just a minimal snippet to start serving requests. 
Error handling is omitted for simplicity.