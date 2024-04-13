package client

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/transport/headers"
	"github.com/lolizeppelin/micro/transport/metadata"
)

// CallFunc represents the individual call func.
type CallFunc func(ctx context.Context, node *micro.Node, req micro.Request, rsp interface{}, opts CallOptions) error

// CallWrapper is a low level wrapper for the CallFunc.
type CallWrapper func(CallFunc) CallFunc

// Wrapper wraps a client and returns a client.
type Wrapper func(Client) Client

// StreamWrapper wraps a Stream and returns the equivalent.
type StreamWrapper func(micro.Stream) micro.Stream

type fromServiceWrapper struct {
	Client
	// headers to inject
	headers metadata.Metadata
}

func (f *fromServiceWrapper) setHeaders(ctx context.Context) context.Context {
	// don't overwrite keys
	return metadata.MergeContext(ctx, f.headers, false)
}

func (f *fromServiceWrapper) Call(ctx context.Context, req micro.Request, rsp interface{}, opts ...CallOption) error {
	ctx = f.setHeaders(ctx)
	return f.Client.Call(ctx, req, rsp, opts...)
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
		metadata.Metadata{
			headers.Prefix + "From-Service": name,
		},
	}
}
