package grpc

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
	tp "github.com/lolizeppelin/micro/transport/grpc/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type grpcTransport struct {
}

func (t *grpcTransport) Dial(addr string, timeout time.Duration, stream bool) (transport.Client, error) {
	if timeout <= 0 {
		timeout = transport.DefaultDialTimeout
	}

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	}

	conn, err := grpc.NewClient(addr, options...)
	if err != nil {
		return nil, err
	}

	c := &grpcTransportClient{
		conn:   conn,
		local:  "localhost",
		remote: addr,
	}

	// create stream
	if stream {
		var s tp.Transport_StreamClient
		s, err = tp.NewTransportClient(conn).Stream(context.Background())
		if err != nil {
			c.Close()
			return nil, err
		}
		c.stream = s
	}

	// return a client
	return c, nil
}

func (t *grpcTransport) String() string {
	return "grpc"
}

func NewTransport() transport.Transport {
	return &grpcTransport{}
}
