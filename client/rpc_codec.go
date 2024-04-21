package client

import (
	errs "errors"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
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
	stream   bool             // 是否流式传输
	protocol *micro.Protocols // 请求与返回协议
	client   transport.Client // 发送客户端

	headers  map[string]string
	response *transport.Message
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

func setHeaders(m *codec.Message, stream bool) {
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
	if stream {
		set(transport.Stream, "true")
	}
}

func newRPCCodec(headers map[string]string, client transport.Client, protocol *micro.Protocols,
	stream bool) codec.Codec {

	c := &rpcCodec{
		client:   client,
		protocol: protocol,
		headers:  headers,
		stream:   stream,
	}
	return c
}

func (c *rpcCodec) Write(message *codec.Message, body interface{}) error {
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
			buff, err := codec.Marshal(c.protocol.Reqeust, body)
			// write to codec
			if err != nil {
				log.Errorf("code write failed %s", err.Error())
				return exc.InternalServerError("go.micro.client.codec", err.Error())
			}
			// set body
			message.Body = buff
		}
	}

	// create new transport message
	msg := transport.Message{
		Header: message.Header,
		Body:   message.Body,
	}

	if c.stream {
		if err := c.client.Send(&msg); err != nil {
			return exc.InternalServerError("go.micro.client.codec", err.Error())
		}
	} else {
		res, err := c.client.Call(&msg)
		if err != nil {
			return err
		}
		c.response = res
	}

	// send the request

	return nil
}

func (c *rpcCodec) ReadHeader(msg *codec.Message, r codec.MessageType) error {
	if c.stream {
		response := new(transport.Message)
		// read message from transport
		if err := c.client.Recv(response); err != nil {
			return exc.InternalServerError("go.micro.client", err.Error())
		}
		c.response = response
	}
	response := c.response
	// set headers from transport
	msg.Header = response.Header
	getHeaders(msg)
	return nil
}

func (c *rpcCodec) ReadBody(b interface{}) error {
	// read body
	// read raw data
	if b == nil {
		return nil
	}
	if v, ok := b.(*codec.Frame); ok {
		v.Data = c.response.Body
		return nil
	}
	return codec.Unmarshal(c.protocol.Response, c.response.Body, b)
}

func (c *rpcCodec) Close() error {
	if err := c.client.Close(); err != nil {
		return exc.InternalServerError("go.micro.client.transport", err.Error())
	}

	return nil
}

func (c *rpcCodec) String() string {
	return "rpc"
}
