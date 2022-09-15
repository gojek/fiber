package grpc

import (
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
)

// CodecName is the name registered for the proto compressor.
const CodecName = "FiberCodec"

type FiberCodec struct{}

func (FiberCodec) Marshal(v interface{}) ([]byte, error) {
	b, ok := v.([]byte)
	if ok {
		return b, nil
	}
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal, message is %T, want proto.Message or []byte", v)
	}
	return proto.Marshal(vv)
}

func (FiberCodec) Unmarshal(data []byte, v interface{}) error {

	myType, ok := v.(io.Writer)
	if ok {
		myType.Write(data)
		return nil
	}
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to marshal, message is %T, want proto.Message or io.writer", v)
	}
	return proto.Unmarshal(data, vv)
}

func (FiberCodec) Name() string {
	return CodecName
}
