package selector

import (
	"errors"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/utils"
	"time"

	"github.com/lolizeppelin/micro/registry/cache"
)

type registrySelector struct {
	so   Options
	rc   cache.Cache
	name string
}

func (c *registrySelector) newCache(ttl time.Duration) cache.Cache {
	return cache.New(c.so.Registry, ttl)
}

func (c *registrySelector) Name() string {
	return c.name
}

func (c *registrySelector) Select(service string, filters ...Filter) (Next, error) {

	// get the service
	// try the cache first
	// if that fails go directly to the registry
	services, err := c.rc.GetService(service)
	if err != nil {
		if errors.Is(err, micro.ErrServiceNotFound) {
			return nil, micro.ErrSelectServiceNotFound
		}
		return nil, err
	}

	// apply the filters
	for _, filter := range filters {
		services, err = filter(services)
		if err != nil {
			return nil, err
		}
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return nil, micro.ErrSelectEndpointNotFound
	}

	return c.so.Strategy(services), nil
}

func (c *registrySelector) Mark(service string, node *micro.Node, err error) {
}

func (c *registrySelector) Reset(service string) {
}

// Close stops the watcher and destroys the cache
func (c *registrySelector) Close() error {
	c.rc.Stop()

	return nil
}

func NewSelector(opts ...Option) (Selector, error) {
	_opts := Options{}

	for _, opt := range opts {
		opt(&_opts)
	}
	if _opts.Registry == nil {
		return nil, fmt.Errorf("no register server found")
	}
	var name string
	if _opts.Strategy == nil {
		name = "score-default"
		_opts.Strategy = NewSharedStrategy(
			[]string{
				utils.RandomHex(16),
				utils.RandomHex(16),
			})
	} else {
		name = "custom"
	}
	s := &registrySelector{
		so:   _opts,
		name: name,
	}
	s.rc = s.newCache(_opts.TTL)
	return s, nil
}
