package server

import (
	"context"
	"errors"
	"github.com/lolizeppelin/micro"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/transport"
	tp "github.com/lolizeppelin/micro/transport/grpc/proto"
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

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic recovered: %v, stack: %s", r, string(debug.Stack()))
			err = exc.InternalServerError("go.micro.server", "panic recovered: %v", r)
		}
	}()

	endpoint, ok := request.Header[transport.Endpoint]
	if !ok {
		return exc.InternalServerError("go.micro.server", "endpoint not found from header")
	}
	serviceName, methodName, err := serviceMethod(endpoint)
	if err != nil {
		return exc.InternalServerError("go.micro.server", err.Error())
	}
	// get grpc metadata
	gmd, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		gmd = metadata.MD{}
	}
	// copy the metadata to go-micro.metadata
	log.Debugf("request %s", endpoint)
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
	md[transport.Method] = request.Header[transport.Method]

	// create new context
	_ctx := transport.NewContext(ctx, md)

	// get peer from context
	if p, ok := peer.FromContext(ctx); ok {
		md["Remote"] = p.Addr.String()
		_ctx = peer.NewContext(_ctx, p)
	}

	// set the timeout if we have it
	if timeout > 0 {
		var cancel context.CancelFunc
		_ctx, cancel = context.WithTimeout(_ctx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	handler := g.service.Handler(serviceName, methodName)
	if handler == nil {
		return status.New(codes.Unimplemented, "unknown service or method").Err()
	}

	protocol, ok := request.Header[micro.ContentType]
	accept, ok := request.Header[micro.Accept]

	if !handler.Match(protocol, accept) {
		log.Warnf("reuqet protocol/accept not match")
	}

	var args []reflect.Value
	args, err = handler.BuildArgs(_ctx, request.Header[micro.ContentType], request.Body)
	if err != nil {
		return err
	}

	results := handler.Method.Func.Call(args)
	if handler.Response == nil {
		return
	}
	if e := results[1].Interface(); e != nil {
		err = e.(error)
		return
	}
	resp := results[0].Interface()
	codec := encoding.GetCodec(request.Header[micro.Accept])
	if codec == nil {
		err = exc.InternalServerError("go.micro.server",
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
