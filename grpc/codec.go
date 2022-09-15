package grpc

import (
	"google.golang.org/grpc/encoding"
	"io"
)

// CodecName is the name registered for the proto compressor.
const CodecName = "FiberCodec"

// FiberCodec is a custom codec to prevent marshaling and unmarshalling
// when unnecessary, base on the inputs
type FiberCodec struct{}

// Marshal will attempt to pass the request directly if it is a byte slice,
// otherwise unmarshal the request proto using the default implementation
func (FiberCodec) Marshal(v interface{}) ([]byte, error) {
	b, ok := v.([]byte)
	if ok {
		return b, nil
	}
	defaultCodec := encoding.GetCodec("proto")
	return defaultCodec.Marshal(v)
}

// Unmarshal will attempt to write the request directly if it is a writer,
// otherwise unmarshal the request proto using the default implementation
func (FiberCodec) Unmarshal(data []byte, v interface{}) error {
	myType, ok := v.(io.Writer)
	if ok {
		_, err := myType.Write(data)
		return err
	}
	defaultCodec := encoding.GetCodec("proto")
	return defaultCodec.Unmarshal(data, v)
}

func (FiberCodec) Name() string {
	return CodecName
}
