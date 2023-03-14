package fiber_test

import (
	"sort"
	"testing"

	"github.com/gojek/fiber"
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
