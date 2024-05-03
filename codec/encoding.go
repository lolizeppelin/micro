package codec

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/encoding"
)

func init() {
	encoding.RegisterCodec(jsonCodec{})
	encoding.RegisterCodec(protoCodec{})
	encoding.RegisterCodec(bytesCodec{})
}

type protoCodec struct{}

func (protoCodec) Marshal(v interface{}) ([]byte, error) {
	m, ok := v.(proto.Message)
	if !ok {
		return nil, ErrInvalidMessage
	}
	return proto.Marshal(m)
}

func (protoCodec) Unmarshal(data []byte, v interface{}) error {
	m, ok := v.(proto.Message)
	if !ok {
		return ErrInvalidMessage
	}
	return proto.Unmarshal(data, m)
}

func (protoCodec) Name() string {
	return "application/grpc+proto"
}

type jsonCodec struct{}

func (jsonCodec) Marshal(v interface{}) ([]byte, error) {

	return json.Marshal(v)
}

func (jsonCodec) Unmarshal(data []byte, v interface{}) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

func (jsonCodec) Name() string {
	return "application/grpc+json"
}

type bytesCodec struct{}

func (bytesCodec) Marshal(v interface{}) ([]byte, error) {
	b, ok := v.([]byte)
	if !ok {
		return nil, ErrInvalidMessage
	}
	return b, nil
}

func (bytesCodec) Unmarshal(data []byte, v interface{}) error {
	if data == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		return ErrInvalidMessage
	}
	copy(b, data)
	return nil
}

func (bytesCodec) Name() string {
	//return "bytes"
	return "application/grpc+bytes"
}
