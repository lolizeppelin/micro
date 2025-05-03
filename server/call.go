package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/lolizeppelin/micro"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/tracing"
	"github.com/lolizeppelin/micro/transport"
	tp "github.com/lolizeppelin/micro/transport/grpc/proto"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	HandlerScope = "micro/server/handler"
)

var (
	_version, _ = micro.NewVersion("1.0.0")
)

func (g *RPCServer) handler(ctx context.Context, msg *tp.Message) (*tp.Message, error) {

	statusCode := codes.OK
	statusDesc := ""

	response := new(tp.Message)
	err := g.processRequest(ctx, msg, response)
	if err != nil {
		var errStatus *status.Status
		var vErr *exc.Error
		switch {
		case errors.As(err, &vErr):
			// micro.Error now proto based and we can attach it to grpc status
			statusCode = exc.GrpcCodeFromMicroError(vErr)
			statusDesc = vErr.Error()
			vErr.Detail = strings.ToValidUTF8(vErr.Detail, "")
			errStatus, err = status.New(statusCode, statusDesc).WithDetails(vErr)
			if err != nil {
				return nil, err
			}
		default:
			// default case user pass own error type that not proto based
			statusCode = exc.ConvertCode(vErr)
			statusDesc = vErr.Error()
			errStatus = status.New(statusCode, statusDesc)
		}
		return nil, errStatus.Err()
	}
	return response, nil

}

func (g *RPCServer) processRequest(ctx context.Context, request, response *tp.Message) (err error) {

	var span oteltrace.Span
	tracer := tracing.GetTracer(HandlerScope, _version)
	ctx, span = tracer.Start(ctx, "process.request",
		oteltrace.WithAttributes(
			attribute.String("service", g.opts.Name),
			attribute.String("version", g.opts.Version.Version()),
		),
	)

	defer func() {
		if r := recover(); r != nil {
			span.AddEvent("panic")
			log.Errorf(ctx, "panic recovered: %v, stack: %s", r, string(debug.Stack()))
			err = exc.InternalServerError("go.micro.server", "panic recovered: %v", r)
		}

		span.End()
	}()

	endpoint, ok := request.Header[transport.Endpoint]
	if !ok {
		span.AddEvent("errors", oteltrace.WithAttributes(attribute.String("headers", "endpoint")))
		return exc.InternalServerError("go.micro.server", "endpoint not found from header")
	}
	serviceName, methodName, err := serviceMethod(endpoint)
	if err != nil {
		span.AddEvent("errors", oteltrace.WithAttributes(attribute.String("endpoint", "split")))
		return exc.InternalServerError("go.micro.server", err.Error())
	}
	// copy the metadata to go-micro.metadata
	log.Debugf(ctx, "request %s", endpoint)

	timeout := int64(0)
	md := make(transport.Metadata)
	for k, v := range request.Header {
		if k == micro.ContentType {
			continue
		}
		md[k] = v
	}
	// get grpc metadata
	if gmd, find := metadata.FromIncomingContext(ctx); find {
		for k, v := range gmd {
			if k == "timeout" && len(v) > 0 {
				timeout, _ = strconv.ParseInt(v[0], 10, 64)
			}
			md[k] = strings.Join(v, ", ")
		}
	}
	// get peer from context
	if p, find := peer.FromContext(ctx); find {
		md["Remote"] = p.Addr.String()
	}
	ctx = transport.NewContext(ctx, md)

	// set the timeout if we have it
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	handler := g.service.Handler(serviceName, methodName)
	if handler == nil {
		span.AddEvent("errors", oteltrace.WithAttributes(attribute.String("handler", "none")))
		return status.New(codes.Unimplemented, "unknown service or method").Err()
	}

	protocol, ok := request.Header[micro.ContentType]
	accept, ok := request.Header[micro.Accept]

	span.SetAttributes(attribute.String("protocol", protocol),
		attribute.String("accept", accept),
		attribute.String("name", handler.Name),
		attribute.String("resource", handler.Resource),
		attribute.Bool("internal", handler.Internal))

	if !handler.Match(protocol, accept) {
		span.AddEvent("missmatch",
			oteltrace.WithAttributes(
				attribute.String("protocol", protocol),
				attribute.String("accept", accept)))
		log.Warnf(ctx, "reuqet protocol/accept not match")
	}

	var args []reflect.Value
	ctx, args, err = handler.BuildArgs(ctx, request.Header[micro.ContentType], request.QueryParams(), request.Body)
	if err != nil {
		return err
	}

	results := handler.Method.Func.Call(args)
	if handler.Response == nil {
		span.AddEvent("success")
		return
	}
	if e := results[1].Interface(); e != nil {
		var match bool
		err, match = e.(error)
		if !match {
			err = fmt.Errorf("unknown handler call resulst")
			span.RecordError(err)
		}
		return
	}
	resp := results[0].Interface()
	codec := encoding.GetCodec(request.Header[micro.Accept])
	if codec == nil {
		err = exc.InternalServerError("micro.handler",
			"response codec '%s' not found", request.Header[micro.Accept])
		return
	}
	var buff []byte
	buff, err = codec.Marshal(resp)
	if err != nil {
		return
	}
	response.Body = buff
	return
}

func (g *RPCServer) Stream(stream tp.Transport_StreamServer) error {
	g.wg.Add(1)
	defer g.wg.Done()
	return status.New(codes.Unimplemented, "stream message not implemented").Err()
}

func (g *RPCServer) Call(ctx context.Context, msg *tp.Message) (*tp.Message, error) {
	g.wg.Add(1)
	defer g.wg.Done()
	return g.handler(ctx, msg)
}
