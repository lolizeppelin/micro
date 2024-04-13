package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func (g *RPCServer) handler(srv interface{}, stream grpc.ServerStream) error {

	g.wg.Add(1)
	defer g.wg.Done()

	msg := new(transport.Message)
	if err := stream.RecvMsg(msg); err != nil {
		return status.New(codes.InvalidArgument, "decode message failed").Err()
	}
	endpoint, ok := msg.Header[transport.Endpoint]
	if !ok {
		return status.New(codes.InvalidArgument, "endpoint not found from header").Err()
	}
	serviceName, methodName, err := serviceMethod(endpoint)
	if err != nil {
		return status.New(codes.InvalidArgument, err.Error()).Err()
	}
	// get grpc metadata
	gmd, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		gmd = metadata.MD{}
	}
	// copy the metadata to go-micro.metadata

	timeout := int64(0)
	md := transport.Metadata{}
	for k, v := range gmd {
		if k == "x-content-type" {
			continue
		}
		if k == "timeout" && len(v) > 0 {
			timeout, _ = strconv.ParseInt(v[0], 10, 64)
		}
		md[k] = strings.Join(v, ", ")
	}
	md[transport.Method] = msg.Header[transport.Method]

	// create new context
	ctx := transport.NewContext(stream.Context(), md)

	// get peer from context
	if p, ok := peer.FromContext(stream.Context()); ok {
		md["Remote"] = p.Addr.String()
		ctx = peer.NewContext(ctx, p)
	}

	// set the timeout if we have it
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	handler := g.service.Handler(serviceName, methodName)
	protocol, ok := msg.Header[micro.ContentType]
	accept, ok := msg.Header[micro.Accept]
	if !codec.MatchCodec(protocol, accept, handler.Metadata["res"], handler.Metadata["req"]) {

	}

	if handler == nil {
		return status.New(codes.Unimplemented, "unknown service or method").Err()
	}

	//if handler.Metadata["req"] == "stream" {
	//	return g.processStream(stream, handler, ctx)
	//}

	return g.processRequest(ctx, stream, handler, msg)
}

func (g *RPCServer) processRequest(ctx context.Context, stream grpc.ServerStream,
	handler *Handler, msg *transport.Message) error {

	args, err := handler.BuildArgs(ctx, msg.Header[micro.ContentType], msg.Body)
	if err != nil {
		return err
	}

	// define the handler func
	fn := func(ctx context.Context) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("panic recovered: %v, stack: %s", r, string(debug.Stack()))
				err = exc.InternalServerError("go.micro.server", "panic recovered: %v", r)
			}
		}()
		results := handler.Method.Func.Call(args)
		if handler.Response == nil {
			return
		}
		resp = results[0].Interface()
		err = results[1].Interface().(error)
		return
	}

	statusCode := codes.OK
	statusDesc := ""
	// execute the handler
	result, err := fn(ctx)
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
				return err
			}
		default:
			// default case user pass own error type that not proto based
			statusCode = exc.ConvertCode(vErr)
			statusDesc = vErr.Error()
			errStatus = status.New(statusCode, statusDesc)
		}
		return errStatus.Err()
	}
	if handler.Response == nil {
		result = new(empty.Empty)
	}
	if err = stream.SendMsg(result); err != nil {
		return err
	}
	return status.New(statusCode, statusDesc).Err()

}
