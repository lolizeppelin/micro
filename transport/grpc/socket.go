package grpc

import (
	"context"
	"github.com/lolizeppelin/micro/transport"
	tp "github.com/lolizeppelin/micro/transport/grpc/proto"
	"google.golang.org/grpc"
)

type grpcTransportClient struct {
	conn   *grpc.ClientConn
	stream tp.Transport_StreamClient

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

	return g.stream.Send(&tp.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (g *grpcTransportClient) Call(m *transport.Message) (*transport.Message, error) {
	if m == nil {
		return nil, nil
	}
	result, err := tp.NewTransportClient(g.conn).Call(context.Background(), &tp.Message{
		Header: m.Header,
		Body:   m.Body,
	})
	if err != nil {
		return nil, err
	}

	return &transport.Message{
		Header: result.Header,
		Body:   result.Body,
	}, nil

}

func (g *grpcTransportClient) Close() error {
	return g.conn.Close()
}
