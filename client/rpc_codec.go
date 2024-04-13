package client

import (
	"bytes"
	errs "errors"
	"github.com/lolizeppelin/micro/codec"
	"github.com/lolizeppelin/micro/codec/grpc"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/transport"
)

const (
	lastStreamResponseError = "EOS"
)

// serverError represents an error that has been returned from
// the remote side of the RPC connection.
type serverError string

func (e serverError) Error() string {
	return string(e)
}

// errShutdown holds the specific error for closing/closed connections.
var (
	errShutdown = errs.New("connection is shut down")
)

type rpcCodec struct {
	client transport.Client
	codec  codec.Codec

	headers map[string]string
	buf     *readWriteCloser

	// signify if its a stream
	stream string
}

type readWriteCloser struct {
	wbuf *bytes.Buffer
	rbuf *bytes.Buffer
}

func (rwc *readWriteCloser) Read(p []byte) (n int, err error) {
	return rwc.rbuf.Read(p)
}

func (rwc *readWriteCloser) Write(p []byte) (n int, err error) {
	return rwc.wbuf.Write(p)
}

func (rwc *readWriteCloser) Close() error {
	rwc.rbuf.Reset()
	rwc.wbuf.Reset()

	return nil
}

func getHeaders(m *codec.Message) {
	set := func(v, hdr string) string {
		if len(v) > 0 {
			return v
		}

		return m.Header[hdr]
	}
	// check error in header
	m.Error = set(m.Error, transport.Error)
	// check endpoint in header
	m.Endpoint = set(m.Endpoint, transport.Endpoint)
	// check method in header
	m.Method = set(m.Method, transport.Method)
	// set the request id
	m.Id = set(m.Id, transport.ID)
}

func setHeaders(m *codec.Message, stream string) {
	set := func(hdr, v string) {
		if len(v) == 0 {
			return
		}

		m.Header[hdr] = v
	}

	set(transport.ID, m.Id)
	set(transport.Service, m.Service)
	set(transport.Method, m.Method)
	set(transport.Endpoint, m.Endpoint)
	set(transport.Error, m.Error)

	if len(stream) > 0 {
		set(transport.Stream, stream)
	}
}

func newRPCCodec(headers map[string]string, client transport.Client, protocol, accept, stream string) codec.Codec {
	rwc := &readWriteCloser{
		wbuf: bytes.NewBuffer(nil),
		rbuf: bytes.NewBuffer(nil),
	}

	return &rpcCodec{
		buf:     rwc,
		client:  client,
		codec:   grpc.NewCodec(rwc, protocol, accept),
		headers: headers,
		stream:  stream,
	}
}

func (c *rpcCodec) Write(message *codec.Message, body interface{}) error {
	c.buf.wbuf.Reset()

	// create header
	if message.Header == nil {
		message.Header = map[string]string{}
	}

	// copy original header
	for k, v := range c.headers {
		message.Header[k] = v
	}

	// set the mucp headers
	setHeaders(message, c.stream)

	// if body is bytes Frame don't encode
	if body != nil {
		if b, ok := body.(*codec.Frame); ok {
			// set body
			message.Body = b.Data
		} else {
			// write to codec
			if err := c.codec.Write(message, body); err != nil {
				return exc.InternalServerError("go.micro.client.codec", err.Error())
			}
			// set body
			message.Body = c.buf.wbuf.Bytes()
		}
	}

	// create new transport message
	msg := transport.Message{
		Header: message.Header,
		Body:   message.Body,
	}

	// send the request
	if err := c.client.Send(&msg); err != nil {
		return exc.InternalServerError("go.micro.client.transport", err.Error())
	}

	return nil
}

func (c *rpcCodec) ReadHeader(msg *codec.Message, r codec.MessageType) error {
	tm := new(transport.Message)
	// read message from transport
	if err := c.client.Recv(tm); err != nil {
		return exc.InternalServerError("go.micro.client.transport", err.Error())
	}

	c.buf.rbuf.Reset()
	c.buf.rbuf.Write(tm.Body)

	// set headers from transport
	msg.Header = tm.Header

	// read header
	err := c.codec.ReadHeader(msg, r)

	// get headers
	getHeaders(msg)

	// return header error
	if err != nil {
		return exc.InternalServerError("go.micro.client.codec", err.Error())
	}

	return nil
}

func (c *rpcCodec) ReadBody(b interface{}) error {
	// read body
	// read raw data
	if v, ok := b.(*codec.Frame); ok {
		v.Data = c.buf.rbuf.Bytes()
		return nil
	}

	if err := c.codec.ReadBody(b); err != nil {
		return exc.InternalServerError("go.micro.client.codec", err.Error())
	}

	return nil
}

func (c *rpcCodec) Close() error {
	if err := c.buf.Close(); err != nil {
		return err
	}

	if err := c.codec.Close(); err != nil {
		return err
	}

	if err := c.client.Close(); err != nil {
		return exc.InternalServerError("go.micro.client.transport", err.Error())
	}

	return nil
}

func (c *rpcCodec) String() string {
	return "rpc"
}
