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

func (r *rpcClient) stream(ctx context.Context, node *micro.Node,
	req micro.Request, opts CallOptions) (micro.Stream, error) {
	address := node.Address

	headers := make(map[string]string)

	md, ok := transport.FromContext(ctx)
	if ok {
		for k, v := range md {
			headers[k] = v
		}
	}

	// set timeout in nanoseconds
	if opts.StreamTimeout > time.Duration(0) {
		headers["Timeout"] = utils.UnsafeToString(opts.StreamTimeout / time.Second)
	}
	protocol := req.Protocols()
	// set the content type for the request
	headers[micro.ContentType] = protocol.Reqeust
	// set the accept header
	headers[micro.Accept] = protocol.Response
	// set old codecs
	c, err := r.opts.Transport.Dial(address, opts.DialTimeout, true)
	if err != nil {
		return nil, exc.InternalServerError("go.micro.client", "connection error: %v", err)
	}

	// increment the sequence number
	seq := atomic.AddUint64(&r.seq, 1) - 1
	// create codec with stream id
	codec := newRPCCodec(headers, c, protocol, false)

	rsp := &rpcResponse{
		socket: c,
		codec:  codec,
	}

	// set request codec
	if r, ok := req.(*rpcRequest); ok {
		r.codec = codec
	}

	releaseFunc := func(_ error) {
		if err = c.Close(); err != nil {
			log.Error(err)
		}
	}

	stream := &rpcStream{
		id:       seq,
		context:  ctx,
		request:  req,
		response: rsp,
		codec:    codec,
		// used to close the stream
		closed: make(chan bool),
		// signal the end of stream,
		sendEOS: true,
		release: releaseFunc,
	}

	// wait for error response
	ch := make(chan error, 1)

	go func() {
		// send the first message
		ch <- stream.Send(req.Body())
	}()

	var grr error

	select {
	case err := <-ch:
		grr = err
	case <-ctx.Done():
		grr = exc.Timeout("go.micro.client", fmt.Sprintf("%v", ctx.Err()))
	}

	if grr != nil {
		// set the error
		stream.Lock()
		stream.err = grr
		stream.Unlock()

		// close the stream
		if err := stream.Close(); err != nil {
			log.Errorf("failed to close stream: %v", err)
		}

		return nil, grr
	}

	return stream, nil
}
