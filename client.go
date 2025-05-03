package micro

import (
	"context"
	"net/url"
)

// Request is the interface for a synchronous request used by Call or Stream.
type Request interface {
	// primary key
	PrimaryKey() string
	// The host to call
	Host() string
	// The service to call
	Service() string
	// The action to take
	Method() string
	// The endpoint to invoke
	Endpoint() string
	// The content type
	Protocols() *Protocols
	// Query string
	Query() url.Values
	// The unencoded request body
	Body() interface{}
	// service version fileter
	Version() *Version
}

type Response struct {
	Headers map[string]string
	Body    interface{}
}

// Stream is the inteface for a bidirectional synchronous stream.
type Stream interface {
	CloseSend() error
	// Send will encode and send a request
	Send(body []byte) error
	// Recv will decode and read a response
	Recv(string, *Response) error

	// Close closes the stream
	Close(ctx context.Context) error
}

type Target struct {
	ID        string     // primary key
	Method    string     // http method
	Host      string     // http host
	Service   string     // service
	Endpoint  string     // service endpoint
	Version   *Version   // service version
	Protocols *Protocols // request protocols
	Query     url.Values // request query parameters
}
