package grpc

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
	tp "github.com/lolizeppelin/micro/transport/grpc/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"time"
)

type grpcTransport struct {
	credentials credentials.TransportCredentials
}

func (t *grpcTransport) Dial(addr string, timeout time.Duration, stream bool) (transport.Client, error) {
	if timeout <= 0 {
		timeout = transport.DefaultDialTimeout
	}

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(t.credentials),       // 证书设置
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // OpenTelemetry数据传递
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) { // 设置链接超时
			return net.DialTimeout("tcp", addr, timeout)
		}),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any,
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error { // 设置默认rpc超时
			if _, ok := ctx.Deadline(); !ok {
				timeoutCtx, cancel := context.WithTimeout(ctx, transport.DefaultRPCTimeout)
				defer cancel()
				return invoker(timeoutCtx, method, req, reply, cc, opts...)
			}
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

func NewTransport(c credentials.TransportCredentials) transport.Transport {
	return &grpcTransport{credentials: c}
}
