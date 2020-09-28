package fiber_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/gojek/fiber"
	"github.com/gojek/fiber/internal/testutils"
)

func TestProxyCaller_Dispatch(t *testing.T) {
	host, path := "http://localhost:8080", "/recommendations/search"
	req := testutils.MockReq("POST", fmt.Sprintf("%s%s", host, path), "")

	dispatcher := new(MockDispatcher)
	dispatcher.On("Do", req).Return(nil)

	backendName := "test-backend"
	backendEndpoint := "http://proxy-test:9090/api"
	backend := fiber.NewBackend(backendName, backendEndpoint)
	caller, _ := fiber.NewCaller(backendName, dispatcher)

	proxy := fiber.NewProxy(backend, caller)

	assert.Equal(t, backendName, proxy.ID())

	<-proxy.Dispatch(context.Background(), req).Iter()

	expectedURL, _ := url.Parse(fmt.Sprintf("%s%s", backendEndpoint, path))

	assert.Equal(t, expectedURL, req.URL)
	dispatcher.AssertExpectations(t)
}
