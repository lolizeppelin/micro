// Package transport is an interface for synchronous connection based communication
package transport

import (
	"context"
	"net/url"
	"time"
)

const (
	DefaultDialTimeout = time.Second * 3
	DefaultRPCTimeout  = time.Second * 30
)

// Transport is an interface which is used for communication between
// services. It uses connection based socket send/recv semantics and
// has various implementations; http, grpc, quic.
type Transport interface {
	Dial(addr string, timeout time.Duration, stream bool) (Client, error)
	String() string
}

// Message is a broker message.
type Message struct {
	Header map[string]string `msgpack:"headers,omitempty"`
	Query  url.Values        `msgpack:"query,omitempty"`
	Body   []byte            `msgpack:"body,omitempty"`
}

type Socket interface {
	Recv(*Message) error
	Send(*Message) error
	Call(context.Context, *Message) (*Message, error)
	Close() error
	Local() string
	Remote() string
}

type Client interface {
	Socket
}
