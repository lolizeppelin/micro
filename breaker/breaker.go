package breaker

import (
	"github.com/lolizeppelin/micro"
	_ "github.com/sony/gobreaker"
)

import (
	"context"
	"sync"

	"github.com/lolizeppelin/micro/client"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/sony/gobreaker"
)

type Scope int

const (
	BreakService Scope = iota
	BreakServiceEndpoint
)

type clientWrapper struct {
	bs    gobreaker.Settings
	scope Scope
	cbs   map[string]*gobreaker.TwoStepCircuitBreaker
	mu    sync.Mutex
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req micro.Request, rsp interface{}, opts ...client.CallOption) error {
	var svc string

	switch c.scope {
	case BreakService:
		svc = req.Service()
	case BreakServiceEndpoint:
		svc = req.Service() + "." + req.Endpoint()
	}

	c.mu.Lock()
	cb, ok := c.cbs[svc]
	if !ok {
		cb = gobreaker.NewTwoStepCircuitBreaker(c.bs)
		c.cbs[svc] = cb
	}
	c.mu.Unlock()

	cbAllow, err := cb.Allow()
	if err != nil {
		return exc.New(req.Service(), err.Error(), 502)
	}

	if err = c.Client.Call(ctx, req, rsp, opts...); err == nil {
		cbAllow(true)
		return nil
	}

	ex := exc.Parse(err.Error())
	switch {
	case ex.Code == 0:
		ex.Code = 503
	case len(ex.Id) == 0:
		ex.Id = req.Service()
	}

	if ex.Code >= 500 {
		cbAllow(false)
	} else {
		cbAllow(true)
	}

	return ex
}

// NewClientWrapper returns a client Wrapper.
func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		w := &clientWrapper{}
		w.bs = gobreaker.Settings{}
		w.cbs = make(map[string]*gobreaker.TwoStepCircuitBreaker)
		w.Client = c
		return w
	}
}

// NewCustomClientWrapper takes a gobreaker.Settings and BreakerMethod. Returns a client Wrapper.
func NewCustomClientWrapper(bs gobreaker.Settings, scope Scope) client.Wrapper {
	return func(c client.Client) client.Client {
		w := &clientWrapper{}
		w.scope = scope
		w.bs = bs
		w.cbs = make(map[string]*gobreaker.TwoStepCircuitBreaker)
		w.Client = c
		return w
	}
}
