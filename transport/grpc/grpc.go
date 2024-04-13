package grpc

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
	pb "github.com/lolizeppelin/micro/transport/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type grpcTransport struct {
}

func (t *grpcTransport) Dial(addr string, timeout time.Duration) (transport.Client, error) {
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

	// create stream
	stream, err := pb.NewTransportClient(conn).Stream(context.Background())
	if err != nil {
		return nil, err
	}

	// return a client
	return &grpcTransportClient{
		conn:   conn,
		stream: stream,
		local:  "localhost",
		remote: addr,
	}, nil
}

func (t *grpcTransport) String() string {
	return "grpc"
}

func NewTransport() transport.Transport {
	return &grpcTransport{}
}
