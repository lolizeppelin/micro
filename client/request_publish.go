package client

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/registry"
	"github.com/lolizeppelin/micro/transport"
)

func (r *rpcClient) publish(ctx context.Context, request micro.Request, opts ...CallOption) error {

	callOpts := r.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}
	node := callOpts.Node
	// 有node id或版本限定,通过过滤器筛选node,设置节点
	if node == "" && request.Version() != nil {
		next, err := r.next(request, callOpts)
		if err != nil {
			return err
		}
		var n *micro.Node
		n, err = next()
		if err != nil {
			return err
		}
		node = n.Id
	}

	topic := registry.Topic(request.Service(), request.Version(), node)
	headers := transport.CopyFromContext(ctx)
	protocol := request.Protocols()
	headers[micro.ContentType] = protocol.Reqeust
	headers[transport.Service] = request.Service()
	headers[transport.Method] = request.Method() // http method
	headers[transport.Endpoint] = request.Endpoint()

	msg := &transport.Message{
		Header: headers,
	}

	body := request.Body()

	b, err := codec.Marshal(protocol.Reqeust, body)
	if err != nil {
		return exc.InternalServerError("micro.rpc.publish", err.Error())
	}
	// set the body
	msg.Body = b

	l, ok := r.once.Load().(bool)
	if !ok {
		return fmt.Errorf("failed to cast to bool")
	}
	if !l {
		if err = r.opts.Broker.Connect(); err != nil {
			return exc.InternalServerError("micro.rpc.publish", err.Error())
		}
		r.once.Store(true)
	}

	return r.opts.Broker.Publish(ctx, topic, msg)

}
