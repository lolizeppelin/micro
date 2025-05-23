package codec

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
	"io"
)

const (
	MaxMessageSize = 1024 * 1024 * 4 // 4Mb
	maxInt         = int(^uint(0) >> 1)
)

// Unmarshal 返回值反序列化
func Unmarshal(protocol string, buff []byte, payload *micro.Response) error {
	switch protocol {
	case "application/grpc+bytes":
		payload.Body = buff
	case "application/msgpack", "application/grpc+msgpack":
		return msgpack.Unmarshal(buff, payload.Body)
	case "application/grpc+json", "application/json":
		return json.Unmarshal(buff, payload.Body)
	case "application/grpc+proto", "application/grpc":
		pb, ok := payload.Body.(proto.Message)
		if !ok {
			return fmt.Errorf("proto.Message requried for codec.Unmarshal")
		}
		return proto.Unmarshal(buff, pb)
	}
	return fmt.Errorf("protocol '%s' not support for codec.Unmarshal", protocol)
}

// Marshal 发送值序列化
func Marshal(protocol string, b interface{}) ([]byte, error) {
	if v, ok := b.([]byte); ok { // 已经序列化
		return v, nil
	}
	switch protocol {
	case "application/msgpack", "application/grpc+msgpack":
		return msgpack.Marshal(b)
	case "application/grpc+json", "application/json":
		return json.Marshal(b)
	case "application/grpc+proto", "application/grpc":
		pb, ok := b.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("proto.Message requried for codec.Marshal")

		}
		return proto.Marshal(pb)
	}
	return nil, fmt.Errorf("protocol '%s' not support for codec.Marshal", protocol)

}

// StreamRead 流式读
func StreamRead(r io.Reader) (uint8, []byte, error) {
	header := make([]byte, 5)

	// read the header
	if _, err := r.Read(header); err != nil {
		return uint8(0), nil, err
	}

	// get encoding format e.g compressed
	cf := uint8(header[0])

	// get message length
	length := binary.BigEndian.Uint32(header[1:])

	// no encoding format
	if length == 0 {
		return cf, nil, nil
	}

	//
	if int64(length) > int64(maxInt) {
		return cf, nil, fmt.Errorf("grpc: received message larger than max "+
			"length allowed on current machine (%d vs. %d)", length, maxInt)
	}
	if int(length) > MaxMessageSize {
		return cf, nil, fmt.Errorf("grpc: received message larger than max (%d vs. %d)", length, MaxMessageSize)
	}

	msg := make([]byte, int(length))

	if _, err := r.Read(msg); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return cf, nil, err
	}

	return cf, msg, nil
}

// StreamWrite 流式写
func StreamWrite(cf uint8, buf []byte, w io.Writer) error {
	header := make([]byte, 5)

	// set compression
	header[0] = byte(cf)

	// write length as header
	binary.BigEndian.PutUint32(header[1:], uint32(len(buf)))

	// read the header
	if _, err := w.Write(header); err != nil {
		return err
	}

	// write the buffer
	_, err := w.Write(buf)
	return err
}
