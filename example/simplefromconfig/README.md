## Simple Example 

Basic example shows how to define a fiber component with YAML config,
initialize this defined component and start serving traffic with `net/http` handler.

### Configuration

```yaml
id: eager_router
type: EAGER_ROUTER    # set the type of the top-level fiber component
strategy:             # set routing strategy, that would be used with given router
  type: fiber.RandomRoutingStrategy 
routes:               # set an arbitrary number of routes for given router. 
                      # each item should be a valid fiber component definition 
  - id: route_a
    type: PROXY
    timeout: "20s"
    endpoint: "http://localhost:8080/routes/route-a"
  - id: route_b
    type: PROXY
    timeout: "40s"
    endpoint: "http://localhost:8080/routes/route-b"
```

### Start serving traffic

```go
package main

import (
    "net/http"
    "time"

    "github.com/gojek/fiber/config"
    fiberhttp "github.com/gojek/fiber/http"
)

func main() {
    component, _ := config.FromConfig("./example/simple/fiber.yaml")
    fiberHandler := fiberhttp.NewHandler(component, fiberhttp.Options{
        Timeout: 20 * time.Second,
    })

    http.ListenAndServe(":8080", fiberHandler)
}
```

NOTE: This is just a minimal snippet to start serving requests. 
Error handling is omitted for simplicity.