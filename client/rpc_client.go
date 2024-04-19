package client

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lolizeppelin/micro"
	"sync/atomic"
	"time"

	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/transport"
)

const (
	packageID = "go.micro.client"
)

type rpcClient struct {
	opts Options
	once atomic.Value
	pool *transport.Pool

	seq uint64
}

// next returns an iterator for the next nodes to call.

func (r *rpcClient) Call(ctx context.Context, request micro.Request, response interface{}, opts ...CallOption) error {

	// make a copy of call opts
	callOpts := r.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	next, err := r.next(request, callOpts)
	if err != nil {
		return err
	}

	// check if we already have a deadline
	d, ok := ctx.Deadline()
	if !ok {
		// no deadline so we create a new one
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, callOpts.RequestTimeout)

		defer cancel()
	} else {
		// got a deadline so no need to setup context
		// but we need to set the timeout we pass along
		opt := WithRequestTimeout(time.Until(d))
		opt(&callOpts)
	}

	// should we noop right here?
	select {
	case <-ctx.Done():
		return exc.Timeout("go.micro.client", fmt.Sprintf("%v", ctx.Err()))
	default:
	}

	// make copy of call method
	rcall := r.call

	// wrap the call in reverse
	for i := len(callOpts.CallWrappers); i > 0; i-- {
		rcall = callOpts.CallWrappers[i-1](rcall)
	}

	// return errors.New("go.micro.client", "request timeout", 408)
	call := func(i int) error {
		// call backoff first. Someone may want an initial start delay
		var t time.Duration
		t, err = callOpts.Backoff(ctx, request, i)
		if err != nil {
			return exc.InternalServerError("go.micro.client", "backoff error: %v", err.Error())
		}

		// only sleep if greater than 0
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		// select next node
		var node *micro.Node
		node, err = next()
		service := request.Service()

		if err != nil {
			return err
		}

		log.Debugf("find service %s", node.Id)

		// make the call
		err = rcall(ctx, node, request, response, callOpts)
		r.opts.Selector.Mark(service, node, err)

		return err
	}

	// get the retries
	retries := callOpts.Retries

	// disable retries when using a proxy
	// Note: I don't see why we should disable retries for proxies, so commenting out.
	// if _, _, ok := net.Proxy(request.Service(), callOpts.Address); ok {
	// 	retries = 0
	// }

	ch := make(chan error, retries+1)

	var gerr error

	for i := 0; i <= retries; i++ {
		go func(i int) {
			ch <- call(i)
		}(i)

		select {
		case <-ctx.Done():
			return exc.Timeout("go.micro.client", fmt.Sprintf("call timeout: %v", ctx.Err()))
		case err = <-ch:
			// if the call succeeded lets bail early
			if err == nil {
				return nil
			}

			retry, rerr := callOpts.Retry(ctx, request, i, err)
			if rerr != nil {
				return rerr
			}

			if !retry {
				return err
			}
			if err != nil {
				log.Error(err.Error())
			}
			log.Debugf("Retrying request. Previous attempt failed with: %v", err)

			gerr = err
		}
	}

	return gerr
}

func (r *rpcClient) Publish(ctx context.Context, request micro.Request, opts ...CallOption) error {
	if request.Protocols().Reqeust == "" {
		return fmt.Errorf("content type required on publish")
	}
	if r.opts.Broker != nil {
		return r.publish(ctx, request, opts...)
	}
	return r.Call(ctx, request, new(empty.Empty), opts...)
}

func (r *rpcClient) Stream(ctx context.Context, request micro.Request, opts ...CallOption) (micro.Stream, error) {

	// make a copy of call opts
	callOpts := r.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	next, err := r.next(request, callOpts)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, exc.Timeout("go.micro.client", fmt.Sprintf("%v", ctx.Err()))
	default:
	}

	call := func(i int) (micro.Stream, error) {
		// call backoff first. Someone may want an initial start delay
		var t time.Duration
		t, err = callOpts.Backoff(ctx, request, i)
		if err != nil {
			return nil, exc.InternalServerError("go.micro.client", "backoff error: %v", err.Error())
		}

		// only sleep if greater than 0
		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		var node *micro.Node
		node, err = next()
		service := request.Service()

		if err != nil {
			return nil, err
		}

		var stream micro.Stream
		stream, err = r.stream(ctx, node, request, callOpts)
		r.opts.Selector.Mark(service, node, err)

		return stream, err
	}

	type response struct {
		stream micro.Stream
		err    error
	}

	// get the retries
	retries := callOpts.Retries

	ch := make(chan response, retries+1)

	var grr error

	for i := 0; i <= retries; i++ {
		go func(i int) {
			s, sErr := call(i)
			ch <- response{s, sErr}
		}(i)

		select {
		case <-ctx.Done():
			return nil, exc.Timeout("go.micro.client", fmt.Sprintf("call timeout: %v", ctx.Err()))
		case rsp := <-ch:
			// if the call succeeded lets bail early
			if rsp.err == nil {
				return rsp.stream, nil
			}

			retry, rerr := callOpts.Retry(ctx, request, i, rsp.err)
			if rerr != nil {
				return nil, rerr
			}

			if !retry {
				return nil, rsp.err
			}

			grr = rsp.err
		}
	}

	return nil, grr
}
