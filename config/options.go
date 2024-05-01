package config

import (
	"context"
	"github.com/lolizeppelin/micro/config/loader"
	"github.com/lolizeppelin/micro/config/reader"
	"github.com/lolizeppelin/micro/config/source"
)

// WithLoader sets the loader for manager config.
func WithLoader(l loader.Loader) Option {
	return func(o *Options) {
		o.Loader = l
	}
}

// WithSource appends a source to list of sources.
func WithSource(s source.Source) Option {
	return func(o *Options) {
		o.Source = append(o.Source, s)
	}
}

// WithReader sets the config reader.
func WithReader(r reader.Reader) Option {
	return func(o *Options) {
		o.Reader = r
	}
}

func WithWatcherDisabled() Option {
	return func(o *Options) {
		o.WithWatcherDisabled = true
	}
}

type Options struct {
	Loader loader.Loader
	Reader reader.Reader

	// for alternative data
	Context context.Context

	Source []source.Source

	WithWatcherDisabled bool
}

type Option func(o *Options)
