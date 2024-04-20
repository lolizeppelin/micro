package grpc

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
	pb "github.com/lolizeppelin/micro/transport/grpc/proto"
	"google.golang.org/grpc"
)

type grpcTransportClient struct {
	conn   *grpc.ClientConn
	stream pb.Transport_StreamClient

	local  string
	remote string
}

func (g *grpcTransportClient) Local() string {
	return g.local
}

func (g *grpcTransportClient) Remote() string {
	return g.remote
}

func (g *grpcTransportClient) Recv(m *transport.Message) error {
	if m == nil {
		return nil
	}

	msg, err := g.stream.Recv()
	if err != nil {
		return err
	}

	m.Header = msg.Header
	m.Body = msg.Body
	return nil
}

func (g *grpcTransportClient) Send(m *transport.Message) error {
	if m == nil {
		return nil
	}

	return g.stream.Send(&pb.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (g *grpcTransportClient) Close() error {
	return g.conn.Close()
}

func (g *grpcTransportClient) CloseSend() error {
	err := g.stream.CloseSend()
	if err != nil {
		return nil
	}
	stream, err := pb.NewTransportClient(g.conn).Stream(context.Background())
	if err != nil {
		return err
	}
	g.stream = stream
	return nil
}
