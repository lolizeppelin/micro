package client

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/transport"
	"github.com/lolizeppelin/micro/utils"
	"sync/atomic"
	"time"
)

func (r *rpcClient) call(ctx context.Context, node *micro.Node, req micro.Request, resp interface{}, opts CallOptions) error {
	address := node.Address

	headers := transport.CopyFromContext(ctx)
	protocol := req.ContentType()
	accept := req.Accept()
	headers[micro.ContentType] = protocol
	headers[transport.Service] = req.Service()
	headers[transport.Method] = req.Method()
	headers[transport.Endpoint] = req.Endpoint()

	// Set connection timeout for single requests to the server. Should be > 0
	// as otherwise requests can't be made.
	cTimeout := opts.ConnectionTimeout / time.Second
	if cTimeout == 0 {
		log.Debugf("connection timeout was set to 0, overwrite to default connection timeout")
		cTimeout = DefaultConnectionTimeout / time.Second
	}
	// set timeout in nanoseconds
	headers["Timeout"] = utils.UnsafeToString(cTimeout)

	c, err := r.pool.Get(address, opts.DialTimeout)
	if err != nil {
		return exc.InternalServerError("go.micro.client", "connection error: %v", err)
	}

	seq := atomic.AddUint64(&r.seq, 1) - 1
	codec := newRPCCodec(headers, c, protocol, accept, "")

	rsp := &rpcResponse{
		socket: c,
		codec:  codec,
	}

	releaseFunc := func(err error) {
		if err = r.pool.Release(c, err); err != nil {
			log.Errorf("failed to release pool: %s", err.Error())
		}
	}

	stream := &rpcStream{
		id:       fmt.Sprintf("%d", seq),
		context:  ctx,
		request:  req,
		response: rsp,
		codec:    codec,
		closed:   make(chan bool),
		close:    opts.ConnClose,
		release:  releaseFunc,
		sendEOS:  false,
	}

	// close the stream on exiting this function
	defer func() {
		if err := stream.Close(); err != nil {
			log.Errorf("failed to close stream %s", err.Error())
		}
	}()

	// wait for error response
	ch := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- exc.InternalServerError("go.micro.client", "panic recovered: %v", r)
			}
		}()

		// send request
		if err := stream.Send(req.Body()); err != nil {
			ch <- err
			return
		}

		// recv request
		if err := stream.Recv(resp); err != nil {
			ch <- err
			return
		}

		// success
		ch <- nil
	}()

	var grr error

	select {
	case err := <-ch:
		return exc.ClientError(err)
	case <-time.After(cTimeout):
		grr = exc.Timeout("go.micro.client", fmt.Sprintf("%v", ctx.Err()))
	}

	// set the stream error
	if grr != nil {
		stream.Lock()
		stream.err = grr
		stream.Unlock()

		return grr
	}

	return nil
}
