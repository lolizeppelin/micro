package registry

import (
	clientV3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Options struct {
	Client  *clientV3.Client
	TTL     time.Duration
	Timeout time.Duration
}

type Option func(*Options)

func WithClient(client *clientV3.Client) Option {
	return func(o *Options) {
		o.Client = client
	}
}

// WithTTL timeout with seconds
func WithTTL(timeout int32) Option {
	return func(o *Options) {
		o.TTL = time.Second * time.Duration(timeout)
	}
}

func WithTimeout(seconds int32) Option {
	return func(o *Options) {
		o.Timeout = time.Second * time.Duration(seconds)
	}
}
