package micro

import (
	"context"
	"github.com/lolizeppelin/micro/codec"
)

// Request is the interface for a synchronous request used by Call or Stream.
type Request interface {
	// The service to call
	Service() string
	// The action to take
	Method() string
	// The endpoint to invoke
	Endpoint() string
	// The content type
	Protocols() *Protocols
	// The unencoded request body
	Body() interface{}
	// service version fileter
	Version() *Version
}

// Response is the response received from a service.
type Response interface {
	// Read the response
	Codec() codec.Reader
	// read the header
	Header() map[string]string
	// Read the undecoded response
	Read() ([]byte, error)
}

// Stream is the inteface for a bidirectional synchronous stream.
type Stream interface {
	CloseSend() error
	// Context for the stream
	Context() context.Context
	// The request made
	Request() Request
	// The response read
	Response() Response
	// Send will encode and send a request
	Send(interface{}) error
	// Recv will decode and read a response
	Recv(interface{}) error
	// Error returns the stream error
	Error() error
	// Close closes the stream
	Close() error
}

type Target struct {
	Method   string   // http method
	Service  string   // service
	Version  *Version // service version
	Endpoint string   // service endpoint
}
