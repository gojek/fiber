package config_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/gojek/fiber/config"
)

type durCfgTestSuite struct {
	input    string
	duration time.Duration
	success  bool
}

func TestDurationUnmarshalJSON(t *testing.T) {
	tests := map[string]durCfgTestSuite{
		"valid_seconds": {
			input:    "2s",
			duration: time.Second * 2,
			success:  true,
		},
		"valid_minute": {
			input:    "1m",
			duration: time.Minute,
			success:  true,
		},
		"valid_quoted_time": {
			input:    "\"2s\"",
			duration: time.Second * 2,
			success:  true,
		},
		"invalid_input": {
			input:    "xyz",
			duration: 0,
			success:  false,
		},
	}

	// Run the tests
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			var d config.Duration
			// Unmarshal
			err := d.UnmarshalJSON([]byte(data.input))
			// Verify
			assert.Equal(t, data.duration, time.Duration(d))
			assert.Equal(t, data.success, err == nil)
		})
	}
}

func TestDurationMarshalJSON(t *testing.T) {
	duration := config.Duration(time.Second * 2)
	data, err := json.Marshal(duration)
	assert.Equal(t, `"2s"`, string(data))
	assert.NoError(t, err)
}
