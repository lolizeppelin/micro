package selector

import (
	"github.com/lolizeppelin/micro"
	"time"
)

type Options struct {
	Registry micro.Registry
	Strategy Strategy
	Seconds  time.Duration
}

type Option func(*Options)

// WithRegistry sets the registry used by the selector.
func WithRegistry(r micro.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// WithStrategy sets the default strategy for the selector.
func WithStrategy(fn Strategy) Option {
	return func(o *Options) {
		o.Strategy = fn
	}
}

// WithCacheSeconds sets the seconds of cache ttl
func WithCacheSeconds(seconds time.Duration) Option {
	return func(o *Options) {
		o.Seconds = time.Second * seconds
	}
}
