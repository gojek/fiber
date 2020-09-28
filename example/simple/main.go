package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/extras"

	fiberhttp "github.com/gojek/fiber/http"

	"github.com/gojek/fiber/example/helpers"
)

func main() {
	// initialize root-level component
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

	// specify options for component's net/http handler
	options := fiberhttp.Options{
		Timeout: 20 * time.Second,
	}

	// create net/http handler from the component component and handler options
	fiberHandler := fiberhttp.NewHandler(component, options)

	// helper: serve sample responses at /routes/** path
	http.Handle(
		"/routes/",
		http.StripPrefix("/routes", helpers.NewSampleServer("route-a", "route-b")))

	// register component http handler
	http.Handle("/", fiberHandler)

	log.Printf("Listening on port :8080")
	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		log.Printf("\nerror: %v\n", err)
		return
	}
}
