package client

import (
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/broker"
	"github.com/lolizeppelin/micro/selector"
	"github.com/lolizeppelin/micro/transport"
	"time"
)

var (
	// DefaultBackoff is the default backoff function for retries.
	DefaultBackoff = exponentialBackoff
	// DefaultRetry is the default check-for-retry function for retries.
	DefaultRetry = RetryOnError
)

const (
	// DefaultRetries is the default number of times a request is tried.
	DefaultRetries = 5
	// DefaultRequestTimeout is the default request timeout.
	DefaultRequestTimeout = time.Second * 30
	// DefaultPoolSize sets the connection pool size.
	DefaultPoolSize = 100
	// DefaultPoolTTL sets the connection pool ttl.
	DefaultPoolTTL = time.Minute * 1
)

// Options are the Client options.
type Options struct {
	Registry  micro.Registry
	Selector  selector.Selector
	Transport transport.Transport

	// Plugged interfaces
	Broker broker.Broker

	// Middleware for client
	Wrappers []Wrapper

	// Connection Pool
	PoolSize int
	PoolTTL  time.Duration

	// Default Call Options
	CallOptions CallOptions
}

// NewOptions creates new Client options.
func NewOptions(options ...Option) Options {
	opts := Options{
		CallOptions: CallOptions{
			Backoff:           DefaultBackoff,
			Retry:             DefaultRetry,
			Retries:           DefaultRetries,
			RequestTimeout:    DefaultRequestTimeout,
			ConnectionTimeout: transport.DefaultDialTimeout,
			DialTimeout:       transport.DefaultDialTimeout,
		},
		PoolSize: DefaultPoolSize,
		PoolTTL:  DefaultPoolTTL,
	}

	for _, o := range options {
		o(&opts)
	}

	return opts
}

func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

func Registry(r micro.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// PoolSize sets the connection pool size.
func PoolSize(d int) Option {
	return func(o *Options) {
		o.PoolSize = d
	}
}

// PoolTTL sets the connection pool ttl.
func PoolTTL(d time.Duration) Option {
	return func(o *Options) {
		o.PoolTTL = d
	}
}

// Adds a Wrapper to a list of options passed into the client.
func Wrap(w Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w)
	}
}

// Adds a Wrapper to the list of CallFunc wrappers.
func WrapCall(cw ...CallWrapper) Option {
	return func(o *Options) {
		o.CallOptions.CallWrappers = append(o.CallOptions.CallWrappers, cw...)
	}
}

// Backoff is used to set the backoff function used
// when retrying Calls.
func Backoff(fn BackoffFunc) Option {
	return func(o *Options) {
		o.CallOptions.Backoff = fn
	}
}

// Retries set the number of retries when making the request.
func Retries(i int) Option {
	return func(o *Options) {
		o.CallOptions.Retries = i
	}
}

// Retry sets the retry function to be used when re-trying.
func Retry(fn RetryFunc) Option {
	return func(o *Options) {
		o.CallOptions.Retry = fn
	}
}

// RequestTimeout set the request timeout.
func RequestTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.RequestTimeout = d
	}
}

// StreamTimeout sets the stream timeout.
func StreamTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.StreamTimeout = d
	}
}

// DialTimeout sets the transport dial timeout.
func DialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.DialTimeout = d
	}
}

// CallOptions are options used to make calls to a server.
type CallOptions struct {

	// Node id of remote hosts
	Node string
	// node version
	Internal bool

	// Request/Response timeout of entire srv.Call, for single request timeout set ConnectionTimeout.
	RequestTimeout time.Duration
	// Stream timeout for the stream
	StreamTimeout time.Duration

	// Backoff func
	Backoff BackoffFunc
	// Check if retriable func
	Retry RetryFunc
	// node filters
	Filters []selector.Filter
	// Middleware for low level call func
	CallWrappers []CallWrapper
	// ConnectionTimeout of one request to the server.
	// Set this lower than the RequestTimeout to enable retries on connection timeout.
	ConnectionTimeout time.Duration
	// Duration to cache the response for
	CacheExpiry time.Duration
	// Transport Dial Timeout. Used for initial dial to establish a connection.
	DialTimeout time.Duration
	// Number of Call attempts
	Retries int
	// Use the services own auth token
	ServiceToken bool
	// ConnClose sets the Connection: close header.
	ConnClose bool
}

func WithSelectFilters(filters ...selector.Filter) CallOption {
	return func(o *CallOptions) {
		o.Filters = append(o.Filters, filters...)
	}
}

// WithCallWrapper is a CallOption which adds to the existing CallFunc wrappers.
func WithCallWrapper(cw ...CallWrapper) CallOption {
	return func(o *CallOptions) {
		o.CallWrappers = append(o.CallWrappers, cw...)
	}
}

// WithBackoff is a CallOption which overrides that which
// set in Options.CallOptions.
func WithBackoff(fn BackoffFunc) CallOption {
	return func(o *CallOptions) {
		o.Backoff = fn
	}
}

// WithRetry is a CallOption which overrides that which
// set in Options.CallOptions.
func WithRetry(fn RetryFunc) CallOption {
	return func(o *CallOptions) {
		o.Retry = fn
	}
}

// WithRetries sets the number of tries for a call.
// This CallOption overrides Options.CallOptions.
func WithRetries(i int) CallOption {
	return func(o *CallOptions) {
		o.Retries = i
	}
}

// WithRequestTimeout is a CallOption which overrides that which
// set in Options.CallOptions.
func WithRequestTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.RequestTimeout = d
	}
}

// WithConnClose sets the Connection header to close.
func WithConnClose() CallOption {
	return func(o *CallOptions) {
		o.ConnClose = true
	}
}

// WithStreamTimeout sets the stream timeout.
func WithStreamTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.StreamTimeout = d
	}
}

// WithDialTimeout is a CallOption which overrides that which
// set in Options.CallOptions.
func WithDialTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.DialTimeout = d
	}
}

// WithServiceToken is a CallOption which overrides the
// authorization header with the services own auth token.
func WithServiceToken() CallOption {
	return func(o *CallOptions) {
		o.ServiceToken = true
	}
}

func WithNode(node string) CallOption {
	return func(o *CallOptions) {
		o.Node = node
	}
}
