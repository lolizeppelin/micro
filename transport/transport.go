// Package transport is an interface for synchronous connection based communication
package transport

import "time"

const (
	DefaultDialTimeout = time.Second * 5
)

// Transport is an interface which is used for communication between
// services. It uses connection based socket send/recv semantics and
// has various implementations; http, grpc, quic.
type Transport interface {
	Dial(addr string, timeout time.Duration) (Client, error)
	String() string
}

// Message is a broker message.
type Message struct {
	Header map[string]string
	Body   []byte
}

type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Close() error
	CloseSend() error
	Local() string
	Remote() string
}

type Client interface {
	Socket
}
