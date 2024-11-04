package client

import (
	"context"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/transport"
	"github.com/lolizeppelin/micro/utils"
	"runtime/debug"
	"sync/atomic"
	"time"
)

func (r *rpcClient) call(ctx context.Context, node *micro.Node,
	request micro.Request, opts CallOptions) (*transport.Message, error) {

	protocol := request.Protocols()
	body, err := codec.Marshal(protocol.Reqeust, request.Body())
	if err != nil {
		return nil, exc.BadRequest("micro.rpc.call", err.Error())
	}

	headers := transport.CopyFromContext(ctx)
	headers[micro.ContentType] = protocol.Reqeust
	headers[micro.Accept] = protocol.Response
	headers[micro.Host] = request.Host()
	headers[micro.PrimaryKey] = request.PrimaryKey()
	headers[transport.Service] = request.Service()
	headers[transport.Method] = request.Method()
	headers[transport.Endpoint] = request.Endpoint()

	// Set connection timeout for single requests to the server. Should be > 0
	// as otherwise requests can't be made.
	cTimeout := opts.ConnectionTimeout
	if cTimeout <= transport.DefaultDialTimeout {
		log.Debugf("overwrite to default connection timeout")
		cTimeout = transport.DefaultDialTimeout
	}
	// set timeout in nanoseconds
	headers["Timeout"] = utils.UnsafeToString(opts.RequestTimeout / time.Second)

	c, err := r.pool.Get(node.Address, opts.DialTimeout)
	if err != nil {
		return nil, exc.InternalServerError("micro.client.call", "connection error: %v", err)
	}

	seq := atomic.AddUint64(&r.seq, 1) - 1
	headers[transport.ID] = utils.UnsafeToString(seq)

	defer func() {

		if e := recover(); e != nil {
			if err != nil {
				err = exc.InternalServerError("micro.client.call", "rpc call panic")
			}
			log.Error("rpc call panic\n%s", debug.Stack())
		}

		if e := r.pool.Release(c, err); e != nil {
			log.Errorf("failed to close stream %v", e.Error())
		}
	}()

	var msg *transport.Message
	msg, err = c.Call(ctx, &transport.Message{
		Header: headers,
		Query:  request.Query(),
		Body:   body,
	})
	if err != nil {
		return nil, exc.ClientError("micro.client.call", err)
	}
	return msg, nil

}
