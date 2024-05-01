package source

import (
	"time"

	"github.com/lolizeppelin/micro/config/encoder"
)

type AuthCreds struct {
	Username string
	Password string
}

type Options struct {
	// Encoder
	Encoder encoder.Encoder

	Prefix      string
	StripPrefix bool
	Address     []string
	AuthCreds   *AuthCreds
	Timeout     time.Duration
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Encoder: encoder.NewEncoder(),
		Timeout: time.Second * 3,
		Prefix:  DefaultPrefix,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithEncoder sets the source encoder.
func WithEncoder(e encoder.Encoder) Option {
	return func(o *Options) {
		o.Encoder = e
	}
}

func WithAddress(a ...string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// WithPrefix sets the key prefix to use.
func WithPrefix(p string) Option {
	return func(o *Options) {
		o.Prefix = p
	}
}

// StripPrefix indicates whether to remove the prefix from config entries, or leave it in place.
func StripPrefix(strip bool) Option {
	return func(o *Options) {
		o.StripPrefix = strip
	}
}

// Auth allows you to specify username/password.
func WithAuth(username, password string) Option {
	return func(o *Options) {
		o.AuthCreds = &AuthCreds{Username: username, Password: password}
	}
}

// WithDialTimeout set the time out for dialing to etcd.
func WithDialTimeout(secondes time.Duration) Option {
	return func(o *Options) {
		o.Timeout = secondes * time.Second
	}
}
