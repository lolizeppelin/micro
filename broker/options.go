package broker

import (
	"context"
	"crypto/tls"
	"github.com/lolizeppelin/micro"
)

type Options struct {

	// Registry used for clustering
	Registry micro.Registry
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	// Handler executed when error happens in broker mesage
	// processing
	ErrorHandler Handler

	TLSConfig *tls.Config
	Address   []string
	Secure    bool
}

type Option func(*Options)

func NewOptions(opts ...Option) *Options {
	options := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &options
}

// Address Address address sets the host addresses to be used by the broker.
func Address(address ...string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

// DisableAutoAck will disable auto acking of messages
// after they have been handled.

// ErrorHandler will catch all broker errors that cant be handled
// in normal way, for example Codec errors.
func ErrorHandler(h Handler) Option {
	return func(o *Options) {
		o.ErrorHandler = h
	}
}

func Registry(r micro.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// Secure communication with the broker.
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

/*------------- Subscribe -------------*/

type SubscribeOptions struct {

	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool
}

// WithQueue sets the name of the queue to share messages on.
func WithQueue(name string) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Queue = name
	}
}

func DisableAutoAck() SubscribeOption {
	return func(o *SubscribeOptions) {
		o.AutoAck = false
	}
}

type SubscribeOption func(*SubscribeOptions)

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}
