package grpc

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
	tp "github.com/lolizeppelin/micro/transport/grpc/proto"
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
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// dial the server
	conn, err := grpc.DialContext(ctx, addr, options...)
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
