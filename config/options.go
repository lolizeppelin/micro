package config

import (
	"time"
)

const (
	DefaultPrefix = "/micro/config/"
)

type AuthCreds struct {
	Username string
	Password string
}

type Options struct {
	Prefix    string
	Address   []string
	AuthCreds *AuthCreds
	Timeout   time.Duration
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Timeout: time.Second * 3,
		Prefix:  DefaultPrefix,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}

// WithEncoder sets the source encoder.

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

// Auth allows you to specify username/password.
func WithAuth(username, password string) Option {
	return func(o *Options) {
		o.AuthCreds = &AuthCreds{Username: username, Password: password}
	}
}

// WithDialTimeout set the time out for dialing to etcd.
func WithDialTimeout(seconds time.Duration) Option {
	return func(o *Options) {
		o.Timeout = seconds * time.Second
	}
}
