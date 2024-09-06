// Package client is an interface for an RPC client
package client

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/selector"
	"github.com/lolizeppelin/micro/transport"
	"github.com/lolizeppelin/micro/transport/grpc"
	"net/url"
)

// Client is the interface used to make requests to services.
// It supports Request/Response via Transport and Publishing via the Broker.
// It also supports bidirectional streaming of requests.
type Client interface {
	Call(ctx context.Context, req micro.Request, opts ...CallOption) (*transport.Message, error)
	RPC(ctx context.Context, req micro.Request, res *micro.Response, opts ...CallOption) error
	Stream(ctx context.Context, req micro.Request, opts ...CallOption) (micro.Stream, error)
	Publish(ctx context.Context, req micro.Request, opts ...CallOption) error
}

// Closer handle client close.
type Closer interface {
	// CloseSend closes the send direction of the stream.
	CloseSend() error
}

// Option used by the Client.
type Option func(*Options)

// CallOption used by Call or Stream.
type CallOption func(*CallOptions)

func NewClient(opts Options) (Client, error) {
	if opts.Registry == nil {
		return nil, fmt.Errorf("not registry found")
	}
	if opts.Transport == nil {
		opts.Transport = grpc.NewTransport()
	}
	if opts.Selector == nil {
		s, _ := selector.NewSelector(selector.WithRegistry(opts.Registry))
		opts.Selector = s
	}

	p := transport.NewPool(
		opts.PoolSize,
		opts.PoolTTL,
		opts.Transport,
	)

	rc := &rpcClient{
		opts: opts,
		pool: p,
		seq:  0,
	}
	rc.once.Store(false)

	c := Client(rc)

	// wrap in reverse
	for i := len(opts.Wrappers); i > 0; i-- {
		c = opts.Wrappers[i-1](c)
	}

	return c, nil
}

func NewRequest(target micro.Target, request interface{}) micro.Request {

	return &rpcRequest{
		service:   target.Service,
		method:    target.Method,
		endpoint:  target.Endpoint,
		query:     target.Query,
		protocols: target.Protocols,
		version:   target.Version,
		body:      request,
	}
}

type rpcRequest struct {
	query     url.Values
	body      interface{}
	service   string
	method    string
	endpoint  string
	protocols *micro.Protocols
	version   *micro.Version
}

func (r *rpcRequest) Protocols() *micro.Protocols {
	return r.protocols
}

func (r *rpcRequest) Query() url.Values {
	return r.query
}

func (r *rpcRequest) Service() string {
	return r.service
}

func (r *rpcRequest) Method() string {
	return r.method
}

func (r *rpcRequest) Endpoint() string {
	return r.endpoint
}

func (r *rpcRequest) Body() interface{} {
	return r.body
}

func (r *rpcRequest) Version() *micro.Version {
	return r.version
}
