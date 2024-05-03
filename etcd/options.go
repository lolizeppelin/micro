package etcd

import (
	"crypto/tls"
	"time"
)

type AuthCreds struct {
	Username string
	Password string
}

type Options struct {
	TLSConfig   *tls.Config
	Address     []string
	DialTimeout time.Duration
	AuthCreds   *AuthCreds
}

type Option func(*Options)

func WithAddress(address []string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

// WithTimeout dial timeout with seconds
func WithTimeout(seconds int32) Option {
	return func(o *Options) {
		o.DialTimeout = time.Second * time.Duration(seconds)
	}
}

func WithTLSConfig(config *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = config
	}
}

// Auth allows you to specify username/password.
func WithAuth(username, password string) Option {
	return func(o *Options) {
		o.AuthCreds = &AuthCreds{Username: username, Password: password}
	}
}
