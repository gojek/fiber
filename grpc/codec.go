package grpc

import (
	"io"

	"google.golang.org/grpc/encoding"
)

// CodecName is the name registered for the proto compressor.
const CodecName = "FiberCodec"

// FiberCodec is a custom codec to prevent marshaling and unmarshalling
// when unnecessary, base on the inputs
type FiberCodec struct {
	defaultCodec encoding.Codec
}

// Marshal will attempt to pass the request directly if it is a byte slice,
// otherwise unmarshal the request proto using the default implementation
func (fc *FiberCodec) Marshal(v interface{}) ([]byte, error) {
	b, ok := v.([]byte)
	if ok {
		return b, nil
	}
	return fc.getDefaultCodec().Marshal(v)
}

// Unmarshal will attempt to write the request directly if it is a writer,
// otherwise unmarshal the request proto using the default implementation
func (fc *FiberCodec) Unmarshal(data []byte, v interface{}) error {
	writer, ok := v.(io.Writer)
	if ok {
		_, err := writer.Write(data)
		return err
	}
	return fc.getDefaultCodec().Unmarshal(data, v)
}

func (*FiberCodec) Name() string {
	return CodecName
}

func (fc *FiberCodec) getDefaultCodec() encoding.Codec {
	if fc.defaultCodec == nil {
		fc.defaultCodec = encoding.GetCodec("proto")
	}
	return fc.defaultCodec
}
