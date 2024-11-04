package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/transport"
	"github.com/lolizeppelin/micro/utils"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

const (
	lastStreamResponseError = "EOS"
)

type rpcStream struct {
	sync.RWMutex
	closed bool

	headers map[string]string
	client  transport.Client
	// signal whether we should send EOS
	sendEOS bool
}

func (r *rpcStream) Send(body []byte) error {
	r.Lock()
	defer r.Unlock()

	if r.closed {
		return io.EOF
	}

	if err := r.client.Send(&transport.Message{
		Header: r.headers,
		Body:   body,
	}); err != nil {
		return err
	}

	return nil
}

func (r *rpcStream) Recv(protocol string, res *micro.Response) error {
	if r.Closed() {
		return io.EOF
	}

	msg := new(transport.Message)
	err := r.client.Recv(msg)

	if err != nil {
		// 非主动关闭
		if errors.Is(err, io.EOF) && !r.Closed() {
			return io.ErrUnexpectedEOF
		}
		return err
	}

	sErr := msg.Header[transport.Error]
	if sErr != "" {
		if sErr != lastStreamResponseError {
			return exc.InternalServerError("go.micro.stream", "recv error")
		}
	}
	return codec.Unmarshal(protocol, msg.Body, res)
}

func (r *rpcStream) CloseSend() error {
	return errors.New("streamer not implemented")
}

func (r *rpcStream) Closed() bool {
	r.RLock()
	closed := r.closed
	r.RUnlock()
	return closed
}

func (r *rpcStream) Close() error {
	r.Lock()
	if r.closed {
		r.Unlock()
		return nil
	}
	r.closed = true
	r.Unlock()

	if r.sendEOS {
		err := r.client.Send(&transport.Message{
			Header: map[string]string{
				transport.Error: lastStreamResponseError,
			},
		})
		if err != nil {
			log.Errorf("send close package failed: %s", err.Error())
		}
	}
	return r.client.Close()

}

func (r *rpcClient) stream(ctx context.Context, node *micro.Node,
	request micro.Request, opts CallOptions) (micro.Stream, error) {

	protocol := request.Protocols()
	body, err := codec.Marshal(protocol.Reqeust, request.Body())
	if err != nil {
		return nil, exc.BadRequest("micro.client.stream", err.Error())
	}

	headers := transport.CopyFromContext(ctx)
	headers[micro.ContentType] = protocol.Reqeust
	headers[micro.Accept] = protocol.Response
	headers[micro.Host] = request.Host()
	headers[micro.PrimaryKey] = request.PrimaryKey()
	headers[transport.Service] = request.Service()
	headers[transport.Method] = request.Method()
	headers[transport.Endpoint] = request.Endpoint()

	// set timeout in nanoseconds
	if opts.StreamTimeout > time.Duration(0) {
		headers["Timeout"] = utils.UnsafeToString(opts.StreamTimeout / time.Second)
	}
	// set old codecs
	c, err := r.opts.Transport.Dial(node.Address, opts.DialTimeout, true)
	if err != nil {
		return nil, exc.InternalServerError("micro.client.stream", "connection error: %v", err)
	}
	// increment the sequence number
	seq := atomic.AddUint64(&r.seq, 1) - 1
	headers[transport.ID] = utils.UnsafeToString(seq)

	stream := &rpcStream{
		headers: headers,
		client:  c,
		// signal the end of stream,
		sendEOS: true,
	}

	// wait for error response
	ch := make(chan error, 1)

	go func() {
		// send the first message
		ch <- stream.Send(body)
	}()

	select {
	case err = <-ch:
	case <-ctx.Done():
		err = exc.Timeout("go.micro.stream", fmt.Sprintf("%v", ctx.Err()))
	}

	if err != nil {
		// close the stream
		if err = stream.Close(); err != nil {
			log.Errorf("failed to close stream: %v", err)
		}
		return nil, err
	}

	return stream, nil
}
