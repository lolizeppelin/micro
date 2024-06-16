package client

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/transport"
)

// CallFunc represents the individual call func.
type CallFunc func(ctx context.Context, node *micro.Node, req micro.Request, opts CallOptions) (*transport.Message, error)

// CallWrapper is a low level wrapper for the CallFunc.
type CallWrapper func(CallFunc) CallFunc

// Wrapper wraps a client and returns a client.
type Wrapper func(Client) Client

// StreamWrapper wraps a Stream and returns the equivalent.
type StreamWrapper func(micro.Stream) micro.Stream

type fromServiceWrapper struct {
	Client
	// headers to inject
	headers transport.Metadata
}

func (f *fromServiceWrapper) setHeaders(ctx context.Context) context.Context {
	// don't overwrite keys
	return transport.MergeContext(ctx, f.headers, false)
}

func (f *fromServiceWrapper) RPC(ctx context.Context, req micro.Request, res *micro.Response, opts ...CallOption) error {
	ctx = f.setHeaders(ctx)
	return f.Client.RPC(ctx, req, res, opts...)
}

func (f *fromServiceWrapper) Call(ctx context.Context, req micro.Request, opts ...CallOption) (*transport.Message, error) {
	ctx = f.setHeaders(ctx)
	return f.Client.Call(ctx, req, opts...)
}

func (f *fromServiceWrapper) Stream(ctx context.Context, req micro.Request, opts ...CallOption) (micro.Stream, error) {
	ctx = f.setHeaders(ctx)
	return f.Client.Stream(ctx, req, opts...)
}

func (f *fromServiceWrapper) Publish(ctx context.Context, req micro.Request, opts ...CallOption) error {
	ctx = f.setHeaders(ctx)
	return f.Client.Publish(ctx, req, opts...)
}

// FromService wraps a client to inject service and auth metadata.
func FromService(name string, c Client) Client {
	return &fromServiceWrapper{
		c,
		transport.Metadata{
			transport.Prefix + "From-Service": name,
		},
	}
}
