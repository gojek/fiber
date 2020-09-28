package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type payload struct {
	RouteID   string    `json:"route_id"`
	Timestamp time.Time `json:"timestamp"`
}

func sampleHandlerFunc(routeID string) func(writer http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if delay := req.URL.Query().Get("delay"); delay != "" {
			delayDuration, _ := time.ParseDuration(delay)
			time.Sleep(delayDuration)
		}

		rw.WriteHeader(http.StatusOK)

		bodyBytes, _ := json.MarshalIndent(&payload{
			RouteID:   routeID,
			Timestamp: time.Now(),
		}, "", "  ")
		count, _ := rw.Write(bodyBytes)

		log.Printf("[%s]: Response bytes: %d\n", routeID, count)
	}
}

func NewSampleServer(routes ...string) http.Handler {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.HandleFunc(fmt.Sprintf("/%s/", route), sampleHandlerFunc(route))
	}

	return mux
}
