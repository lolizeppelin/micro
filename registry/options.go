package registry

import (
	"crypto/tls"
	"time"
)

type Options struct {
	TLSConfig *tls.Config
	Address   []string
	Timeout   time.Duration
	TTL       time.Duration
}

type Option func(*Options)

func WithAddress(address []string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

// WithTimeout timeout with seconds
func WithTimeout(timeout int32) Option {
	return func(o *Options) {
		o.Timeout = time.Second * time.Duration(timeout)
	}
}

// WithTTL timeout with seconds
func WithTTL(timeout int32) Option {
	return func(o *Options) {
		o.TTL = time.Second * time.Duration(timeout)
	}
}

func WithTLSConfig(config *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = config
	}
}
