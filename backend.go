package fiber

import (
	"fmt"
)

// Backend is an abstraction to be used for defining different backend endpoints for routers or combiners
type Backend interface {
	URL(requestURI string) string
}

type backend struct {
	Name string

	// Fully qualified endpoint URL with protocol
	Endpoint string
}

// NewBackend creates a new backend with the given name and endpoint
func NewBackend(name string, endpoint string) Backend {
	return &backend{
		Name:     name,
		Endpoint: endpoint,
	}
}

// URL Returns the full url of the current backend with a given path part
func (backend *backend) URL(requestURI string) string {
	return fmt.Sprintf("%s%s", backend.Endpoint, requestURI)
}
