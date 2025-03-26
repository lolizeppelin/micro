package broker

import (
	"context"
	"crypto/tls"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/transport"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/vmihailenco/msgpack/v5"
)

type Options struct {
	// Registry s
	Registry micro.Registry
	// Other options for implementations of the interface
	// can be stored in a context

	// Handler executed when error happens in broker message
	// processing
	ErrorHandler ErrorHandler

	TLSConfig *tls.Config
	Address   []string
	Secure    bool
}

type ErrorHandler func(context.Context, string, *kgo.Record, error)

type Option func(*Options)

func _unmarshal(b []byte) (*transport.Message, error) {
	msg := new(transport.Message)
	if err := msgpack.Unmarshal(b, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func _noneHandler(context.Context, string, *kgo.Record, error) {
	return
}

func NewOptions(opts ...Option) *Options {
	options := Options{
		ErrorHandler: _noneHandler,
	}
	for _, o := range opts {
		o(&options)
	}

	return &options
}

func WithRegistry(registry micro.Registry) Option {
	return func(o *Options) {
		o.Registry = registry
	}
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
func WithErrorHandler(fallback ErrorHandler) Option {
	return func(o *Options) {
		o.ErrorHandler = fallback
	}
}

// Secure communication with the broker.
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

/*------------- Subscribe -------------*/

type SubscribeUnmarshal func([]byte) (*transport.Message, error)

type SubscribeOptions struct {

	// Subscribers with the same queue name
	// will create a shared subscription where each
	// receives a subset of messages.
	Queue string

	// AutoAck defaults to true. When a handler returns
	// with a nil error the message is acked.
	AutoAck bool

	// 解析
	Unmarshal SubscribeUnmarshal
}

type SubscribeOption func(*SubscribeOptions)

func NewSubscribeOptions(opts ...SubscribeOption) SubscribeOptions {
	opt := SubscribeOptions{
		AutoAck:   true,
		Unmarshal: _unmarshal,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
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

func UnmarshalHandler(unmarshal SubscribeUnmarshal) SubscribeOption {
	return func(o *SubscribeOptions) {
		o.Unmarshal = unmarshal
	}
}
