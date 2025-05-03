package server

import (
	"context"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
	exc "github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/log"
	"github.com/lolizeppelin/micro/tracing"
	"github.com/lolizeppelin/micro/utils"
	"github.com/xeipuuv/gojsonschema"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/encoding"
	"net/http"
	"net/url"
	"reflect"
)

type Handler struct {
	Resource       string                 // resource name
	Collection     string                 // collection name
	Name           string                 // method name
	Rtype          reflect.Type           // 结构体
	Receiver       reflect.Value          // receiver of method
	Method         reflect.Method         // method stub
	Query          reflect.Type           // 请求url query pram参数校验器
	Request        reflect.Type           // 请求参数
	QueryValidator *gojsonschema.Schema   // 请求参数校验器
	BodyValidator  *gojsonschema.Schema   // 请求载荷校验器
	Response       reflect.Type           // 返回参数
	Metadata       map[string]string      // 元数据
	Internal       bool                   // 内部rpc，不对外, 文档接口与http api需要屏蔽
	Hooks          []micro.PreExecuteHook // 执行前
}

func (handler *Handler) Hook(ctx context.Context, query url.Values, body []byte) (context.Context, error) {
	saved := ctx
	var err error
	for i, hook := range handler.Hooks {
		ctx, err = hook(ctx, query, body)
		if err != nil {
			if ctx == nil { // 确保ctx不会替换为nil
				ctx = saved
			}
			log.Debugf(ctx, "request execute %s hook %d failed", handler.Resource, i)
			return ctx, err
		}
	}
	return ctx, nil
}

/*
BuildArgs rpc转发将请求转数据转化为反射调用参数
*/
func (handler *Handler) BuildArgs(ctx context.Context,
	protocol string, query url.Values, body []byte) (context.Context, []reflect.Value, error) {

	var err error
	var span oteltrace.Span

	tracer := tracing.GetTracer(HandlerScope, _version)
	ctx, span = tracer.Start(ctx, "handler.reflect",
		oteltrace.WithAttributes(
			attribute.String("resource", handler.Resource),
			attribute.String("name", handler.Name),
		),
	)

	defer func() {
		if err != nil && span.IsRecording() {
			span.RecordError(err)
		}
		span.End()
	}()

	if handler.Query == nil && handler.Request == nil {
		ctx, err = handler.Hook(ctx, query, body)
		if err != nil {
			return ctx, nil, err
		}
		span.AddEvent("skip")
		return ctx, []reflect.Value{handler.Receiver, reflect.ValueOf(ctx)}, nil
	}

	var _query reflect.Value
	if handler.Query == nil {
		_query = reflect.ValueOf(nil)
	} else { // validator query parameters
		span.AddEvent("query.decode")
		_query = reflect.New(handler.Query.Elem())
		endpoint := fmt.Sprintf("%s.%s", handler.Resource, handler.Name)
		// url.Values转结构体
		p := _query.Interface()
		if err = codec.UnmarshalQuery(endpoint, query, p); err != nil {
			log.Debug(ctx, "unmarshal request query error")
			err = exc.BadRequest("micro.server", "Unmarshal request query failed")
			return ctx, nil, err.(error)
		}
		span.AddEvent("query.validate")
		// 进行json schema校验
		var result *gojsonschema.Result
		result, err = handler.QueryValidator.Validate(gojsonschema.NewGoLoader(p))
		if err != nil {
			log.Debug(ctx, "validate request query error")
			return ctx, nil, err
		}
		if !result.Valid() {
			log.Debug(ctx, "validate request query failed")
			msg := fmt.Sprintf("Validate request query failed")
			for _, desc := range result.Errors() {
				msg = fmt.Sprintf("%s %s", msg, desc)
			}
			err = exc.BadRequest("micro.server", msg)
			return ctx, nil, err.(error)
		}
	}

	if handler.Request == nil {
		span.AddEvent("query.hook")
		ctx, err = handler.Hook(ctx, query, body)
		if err != nil {
			return ctx, nil, err
		}
		return ctx, []reflect.Value{handler.Receiver, reflect.ValueOf(ctx), _query}, nil
	}

	span.AddEvent("body.protocol")
	// 请求协议
	_codec := encoding.GetCodec(protocol)
	if _codec == nil {
		err = exc.BadRequest("micro.server", "codec not found: '%s'", protocol)
		return ctx, nil, err.(error)
	}
	// body使用jsonschema校验 TODO 非json无法校验,需要处理
	if handler.BodyValidator != nil {
		span.AddEvent("body.validate")
		var result *gojsonschema.Result
		result, err = handler.BodyValidator.Validate(gojsonschema.NewBytesLoader(body))
		if err != nil {
			log.Debugf(ctx, "validate request body error: %v", err)
			err = exc.BadRequest("micro.server", "decode body failed")
			return ctx, nil, err.(error)
		}
		if !result.Valid() {
			log.Debug(ctx, "validate request body failed")
			msg := fmt.Sprintf("validate request body failed")
			for _, desc := range result.Errors() {
				msg = fmt.Sprintf("%s %s", msg, desc)
			}
			err = exc.BadRequest("micro.server", msg)
			return ctx, nil, err.(error)
		}
	}
	// 按协议解析数据流
	var arg reflect.Value
	if handler.Request == utils.TypeOfBytes {
		arg = reflect.Zero(handler.Request)
	} else {
		arg = reflect.New(handler.Request.Elem())
	}
	span.AddEvent("body.decode")
	if err = _codec.Unmarshal(body, arg.Interface()); err != nil {
		return ctx, nil, exc.BadRequest("micro.server", "codec unmarshal failed: %s", err.Error())
	}

	ctx, err = handler.Hook(ctx, query, body)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, []reflect.Value{handler.Receiver, reflect.ValueOf(ctx), _query, arg}, nil
}

/*
Match 判断请求头类型是否匹配
*/
func (handler *Handler) Match(request, response string) bool {
	//return protocol == handler.Metadata["res"] && accept == handler.Metadata["req"]
	return micro.MatchCodec(request, handler.Metadata["req"]) &&
		micro.MatchCodec(response, handler.Metadata["res"])
}

/*
UrlPath 获取api path
*/
func (handler *Handler) UrlPath() (resource, path, method string) {
	path = "/restful/"
	switch handler.Name {
	case "Get":
		method = http.MethodGet
		resource = handler.Resource
		path += fmt.Sprintf("%s/:id", handler.Resource)
	case "List":
		method = http.MethodGet
		resource = handler.Collection
		path += handler.Collection
		return
	case "Create":
		method = http.MethodPost
		resource = handler.Collection
		path += handler.Collection
		return
	case "Update":
		method = http.MethodPut
		resource = handler.Resource
		path += fmt.Sprintf("%s/:id", handler.Resource)
		return
	case "Patch":
		method = http.MethodPatch
		resource = handler.Collection
		path += handler.Collection
		return
	case "Delete":
		method = http.MethodDelete
		resource = handler.Resource
		path += fmt.Sprintf("%s/:id", handler.Resource)
		return
	default: // 非restful接口
		resource = fmt.Sprintf("%s/%s", handler.Resource, handler.Name)
		path = fmt.Sprintf("/%s", resource)
		matches := curdPrefix.FindStringSubmatch(handler.Method.Name)
		if matches == nil { // 这是一个网关接口
			return
		}
		method = matches[1]
	}
	return
}
