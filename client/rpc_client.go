package client

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/tracing"
	"github.com/lolizeppelin/micro/transport"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"time"
)

type rpcClient struct {
	opts Options
	pool *transport.Pool

	seq uint64
}

func (r *rpcClient) Name() string {
	return "grpc"
}

// next returns an iterator for the next nodes to call.

func (r *rpcClient) RPC(ctx context.Context, request micro.Request, response *micro.Response, opts ...CallOption) error {
	msg, err := r.Call(ctx, request, opts...)
	if err != nil {
		return err
	}
	response.Headers = msg.Header
	return codec.Unmarshal(request.Protocols().Response, msg.Body, response)
}

func (r *rpcClient) Call(ctx context.Context, request micro.Request, opts ...CallOption) (*transport.Message, error) {

	// make a copy of call opts
	callOpts := r.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	next, err := r.next(ctx, request, callOpts)
	if err != nil {
		return nil, err
	}

	var span oteltrace.Span
	tracer := tracing.GetTracer(CallScope, _version)
	name := fmt.Sprintf("%s.%s.%s", request.Method(), request.Service(), request.Endpoint())

	defer span.End()

	// check if we already have a deadline
	d, ok := ctx.Deadline()
	if !ok {
		ctx, span = tracer.Start(ctx, name,
			oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
			oteltrace.WithAttributes(
				attribute.String("rpc.transport", r.Name()),
			),
			oteltrace.WithAttributes(
				attribute.Int64("timeout", int64(callOpts.RequestTimeout)),
			))
		// no deadline so we create a new one
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, callOpts.RequestTimeout)
		defer cancel()
	} else {
		// got a deadline so no need to setup context
		// but we need to set the timeout we pass along
		ctx, span = tracer.Start(ctx, name,
			oteltrace.WithSpanKind(oteltrace.SpanKindInternal),
			oteltrace.WithAttributes(
				attribute.String("call.node", callOpts.Node),
			),
			oteltrace.WithAttributes(
				attribute.String("rpc.transport", r.Name()),
			),
			oteltrace.WithAttributes(
				attribute.Int64("deadline", int64(time.Until(d))),
			))
		opt := WithRequestTimeout(time.Until(d))
		opt(&callOpts)
	}

	// should we noop right here?
	select {
	case <-ctx.Done():
		return nil, exc.Timeout("micro.client.call", fmt.Sprintf("%v", ctx.Err()))
	default:
	}

	// make copy of call method
	rcall := r.call

	// wrap the call in reverse
	for i := len(callOpts.CallWrappers); i > 0; i-- {
		rcall = callOpts.CallWrappers[i-1](rcall)
	}

	var res *transport.Message

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

		// make the call
		res, err = rcall(ctx, node, request, callOpts)
		r.opts.Selector.Mark(service, node, err)

		return err
	}

	// get the retries
	retries := callOpts.Retries

	ch := make(chan error, retries+1)

	for i := 0; i <= retries; i++ {
		go func(i int) {
			ch <- call(i)
		}(i)

		select {
		case <-ctx.Done():
			return nil, exc.Timeout("go.micro.client", fmt.Sprintf("call timeout: %v", ctx.Err()))
		case err = <-ch:
			// if the call succeeded lets bail early
			if err == nil {
				return res, nil
			}

			retry, rErr := callOpts.Retry(ctx, request, i, err)
			if rErr != nil {
				return nil, rErr
			}

			if !retry {
				return nil, err
			}
			log.Debugf(ctx, "Retrying request. Previous attempt failed with: %v", err)
		}
	}

	return nil, err
}

func (r *rpcClient) Publish(ctx context.Context, request micro.Request, opts ...CallOption) error {
	if request.Protocols().Reqeust == "" {
		return fmt.Errorf("content type required on publish")
	}
	if r.opts.Broker != nil {
		return r.publish(ctx, request, opts...)
	}
	_, err := r.Call(ctx, request, opts...)
	return err
}

func (r *rpcClient) Stream(ctx context.Context, request micro.Request, opts ...CallOption) (micro.Stream, error) {

	// make a copy of call opts
	callOpts := r.opts.CallOptions
	for _, opt := range opts {
		opt(&callOpts)
	}

	var span oteltrace.Span
	tracer := tracing.GetTracer(StreamScope, _version)
	name := fmt.Sprintf("%s.%s.%s", request.Method(), request.Service(), request.Endpoint())
	ctx, span = tracer.Start(ctx, name,
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(
			attribute.String("stream.node", callOpts.Node),
		),
		oteltrace.WithAttributes(
			attribute.String("rpc.transport", r.Name()),
		),
	)
	defer span.End()

	next, err := r.next(ctx, request, callOpts)
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
