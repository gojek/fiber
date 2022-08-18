package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gojek/fiber/config"
	fiberhttp "github.com/gojek/fiber/http"

	"github.com/gojek/fiber/example/helpers"
)

func main() {
	// initialize root-level fiber component from the config
	component, err := config.FromConfig("./example/simplefromconfig/fiber.yaml")
	if err != nil {
		log.Fatalf("\nerror: %v\n", err)
	}

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
