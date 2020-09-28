package fiber

import "encoding/json"

// Type interface provides a method to do custom initialization of fiber components
type Type interface {
	Initialize(cfgProperties json.RawMessage) error
}

// BaseFiberType implements Type
type BaseFiberType struct{}

// Initialize is a dummy implementation for the Type interface
func (b BaseFiberType) Initialize(_ json.RawMessage) error {
	return nil
}
