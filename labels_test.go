package fiber_test

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/gojek/fiber"
	"github.com/gojek/fiber/extras"
	"github.com/gojek/fiber/internal/testutils"
	testUtilsHttp "github.com/gojek/fiber/internal/testutils/http"
	"github.com/stretchr/testify/assert"
)

func TestLabelsMapKeys(t *testing.T) {
	tests := map[string]struct {
		data     fiber.LabelsMap
		expected []string
	}{
		"empty map": {
			data:     fiber.LabelsMap{},
			expected: []string{},
		},
		"non-empty map": {
			data:     fiber.LabelsMap{"k1": []string{"v1", "v2"}, "k2": []string{"v3"}},
			expected: []string{"k1", "k2"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			keys := tt.data.Keys()
			// Sort and compare
			sort.Slice(keys, func(p, q int) bool {
				return keys[p] < keys[q]
			})
			assert.Equal(t, tt.expected, keys)
		})
	}
}

func TestLabelsMapWithLabel(t *testing.T) {
	tests := map[string]struct {
		data     fiber.LabelsMap
		key      string
		values   []string
		expected fiber.LabelsMap
	}{
		"set key": {
			data:     fiber.LabelsMap{},
			key:      "k1",
			values:   []string{"v1", "v2"},
			expected: fiber.LabelsMap{"k1": []string{"v1", "v2"}},
		},
		"overwrite key": {
			data:     fiber.LabelsMap{"k1": []string{"v1", "v2"}, "k2": []string{"v3"}},
			key:      "k2",
			values:   []string{"new-val"},
			expected: fiber.LabelsMap{"k1": []string{"v1", "v2"}, "k2": []string{"new-val"}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			newMap := tt.data.WithLabel(tt.key, tt.values...)
			assert.Equal(t, tt.expected, newMap)
		})
	}
}

func TestLabelsMapLabel(t *testing.T) {
	tests := map[string]struct {
		data     fiber.LabelsMap
		key      string
		expected []string
	}{
		"empty map": {
			data:     fiber.LabelsMap{},
			key:      "k",
			expected: []string{},
		},
		"non-empty map": {
			data:     fiber.LabelsMap{"k1": []string{"v1", "v2"}, "k2": []string{"v3"}},
			key:      "k1",
			expected: []string{"v1", "v2"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			values := tt.data.Label(tt.key)
			assert.Equal(t, tt.expected, values)
		})
	}
}

// This test uses RandomRoutingStrategy which will always append context of idx
// Labels should be preserved for fiber.CtxComponentLabelsKey context key
func Test_Router_Dispatch_Labels(t *testing.T) {
	lazyRouter := fiber.NewLazyRouter("lazy-router")
	lazyRouter.SetStrategy(new(extras.RandomRoutingStrategy))

	eagerRouter := fiber.NewEagerRouter("eager-router")
	eagerRouter.SetStrategy(new(extras.RandomRoutingStrategy))

	testRouters := []fiber.MultiRouteComponent{eagerRouter, lazyRouter}

	tests := []struct {
		name               string
		initialLabelKey    any
		initialLabelValue  any
		expectedLabelKey   string
		expectedLabelValue string
		router             []fiber.MultiRouteComponent
	}{
		{
			name:               "new label",
			expectedLabelKey:   "idx",
			expectedLabelValue: "0",
			router:             testRouters,
		},
		{
			name:               "overwritten label",
			initialLabelKey:    "idx",
			initialLabelValue:  "111",
			expectedLabelKey:   "idx",
			expectedLabelValue: "0",
			router:             testRouters,
		},
		{
			name:               "existing label not preserved, wrong key",
			initialLabelKey:    "t",
			initialLabelValue:  "11",
			expectedLabelKey:   "t",
			expectedLabelValue: "",
			router:             testRouters,
		},
		{
			name:               "existing label preserved",
			initialLabelKey:    fiber.CtxComponentLabelsKey,
			initialLabelValue:  fiber.NewLabelsMap().WithLabel("t", "11"),
			expectedLabelKey:   "t",
			expectedLabelValue: "11",
			router:             testRouters,
		},
		{
			name:               "existing label not preserved, unexpected value type",
			initialLabelKey:    fiber.CtxComponentLabelsKey,
			initialLabelValue:  map[string]string{"t": "11"},
			expectedLabelKey:   "t",
			expectedLabelValue: "",
			router:             testRouters,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, router := range tt.router {
				router.SetRoutes(map[string]fiber.Component{
					"route-a": testutils.NewMockComponent(
						"route-a",
						testUtilsHttp.DelayedResponse{Response: testUtilsHttp.MockResp(200, "A-OK", nil, nil)}),
				})
				ctx := context.Background()
				if tt.initialLabelKey != nil {
					ctx = context.WithValue(ctx, tt.initialLabelKey, tt.initialLabelValue)
				}
				request := testUtilsHttp.MockReq("POST", "http://localhost:8080/router", "payload")
				resp, ok := <-router.Dispatch(ctx, request).Iter()
				assert.True(t, ok)
				label := strings.Join(resp.Label(tt.expectedLabelKey), ",")
				assert.Equal(t, tt.expectedLabelValue, label)
			}
		})
	}
}
