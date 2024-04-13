// Package grpc provides a grpc codec
package grpc

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/lolizeppelin/micro/codec"
	"io"
)

type Codec struct {
	Conn     io.ReadWriteCloser
	protocol string
	accept   string
}

func (c *Codec) ReadHeader(m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *Codec) ReadBody(b interface{}) error {
	// no body
	if b == nil {
		return nil
	}

	_, buf, err := decode(c.Conn)
	if err != nil {
		return err
	}

	switch c.protocol {
	case "":
		return nil
	case "application/grpc+json", "application/json":
		return json.Unmarshal(buf, b)
	case "application/grpc+proto", "application/grpc":
		return proto.Unmarshal(buf, b.(proto.Message))
	}

	return errors.New("unsupported Content-Type")
}

func (c *Codec) Write(m *codec.Message, b interface{}) error {
	var buf []byte
	var err error

	// marshal content
	switch c.protocol {
	case "":
		return nil
	case "application/grpc+json", "application/json":
		buf, err = json.Marshal(b)
	case "application/grpc+proto", "application/grpc":
		pb, ok := b.(proto.Message)
		if ok {
			buf, err = proto.Marshal(pb)
		}
	default:
		err = errors.New("unsupported Content-Type")
	}
	// check error
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return nil
	}

	return encode(0, buf, c.Conn)
}

func (c *Codec) Close() error {
	return c.Conn.Close()
}

func (c *Codec) String() string {
	return "grpc"
}

func NewCodec(c io.ReadWriteCloser, protocol, accept string) codec.Codec {
	return &Codec{
		Conn:     c,
		protocol: protocol,
		accept:   accept,
	}
}
