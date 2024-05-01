package loader

import (
	"github.com/lolizeppelin/micro/config/reader"
	"github.com/lolizeppelin/micro/config/source"
)

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
