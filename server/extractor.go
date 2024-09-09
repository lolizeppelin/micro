package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/codec"
	"github.com/lolizeppelin/micro/errors"
	"github.com/lolizeppelin/micro/utils"
	"github.com/lolizeppelin/micro/utils/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"google.golang.org/grpc/encoding"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	// 组件非Restful方法
	curdPrefix, _ = regexp.Compile(fmt.Sprintf("^(%s|%s|%s|%s|%s|RPC)_([A-Z].*)$",
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete))
)

func isExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func serviceMethod(endpoint string) (service string, method string, err error) {
	s := strings.Split(endpoint, ".")
	if len(s) != 2 {
		err = fmt.Errorf("endpoint value error")
		return
	}
	if s[0] == "" {
		err = fmt.Errorf("service value error")
		return
	}
	if s[0] == "" {
		err = fmt.Errorf("method value error")
		return
	}
	service = s[0]
	method = s[1]
	return
}

func ServiceMethod(m string) (string, string, error) {

	if len(m) == 0 {
		return "", "", fmt.Errorf("malformed method name: %q", m)
	}

	// grpc method
	if m[0] == '/' {
		// [ , Foo, Bar]
		// [ , package.Foo, Bar]
		// [ , a.package.Foo, Bar]
		parts := strings.Split(m, "/")
		if len(parts) != 3 || len(parts[1]) == 0 || len(parts[2]) == 0 {
			return "", "", fmt.Errorf("malformed method name: %q", m)
		}
		service := strings.Split(parts[1], ".")
		return service[len(service)-1], parts[2], nil
	}

	// non grpc method
	parts := strings.Split(m, ".")

	// expect [Foo, Bar]
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed method name: %q", m)
	}

	return parts[0], parts[1], nil
}

type Handler struct {
	Resource       string               // resource name
	Collection     string               // collection name
	Name           string               // method name
	Rtype          reflect.Type         // 结构体
	Receiver       reflect.Value        // receiver of method
	Method         reflect.Method       // method stub
	Query          reflect.Type         // 请求url query pram参数校验器
	Request        reflect.Type         // 请求参数
	QueryValidator *gojsonschema.Schema // 请求参数校验器
	BodyValidator  *gojsonschema.Schema // 请求载荷校验器
	Response       reflect.Type         // 返回参数
	Metadata       map[string]string    // 元数据
	Internal       bool                 // 内部rpc，不对外
}

/*
BuildArgs rpc转发将请求转数据转化为反射调用参数
*/
func (handler *Handler) BuildArgs(ctx context.Context, protocol string, query url.Values, body []byte) ([]reflect.Value, error) {
	args := []reflect.Value{handler.Receiver, reflect.ValueOf(ctx)}
	if handler.Query == nil && handler.Request == nil {
		return args, nil
	}

	var arg reflect.Value
	if handler.Query == nil {
		arg = reflect.ValueOf(nil)
	} else {
		// validator query parameters
		arg = reflect.New(handler.Query.Elem())
		endpoint := fmt.Sprintf("%s.%s", handler.Resource, handler.Name)
		// url.Values转结构体
		if err := codec.UnmarshalQuery(endpoint, query, arg.Interface()); err != nil {
			return nil, errors.BadRequest("micro.server", "Unmarshal request query failed")
		}
		// query结构体转bytes进行json schema校验
		buff, _ := json.Marshal(arg.Interface())
		result, err := handler.QueryValidator.Validate(gojsonschema.NewBytesLoader(buff))
		if err != nil {
			return nil, err
		}
		if !result.Valid() {
			msg := fmt.Sprintf("Validate request query failed")
			for _, desc := range result.Errors() {
				msg = fmt.Sprintf("%s %s", msg, desc)
			}
			return nil, errors.BadRequest("micro.server", msg)
		}
	}

	args = append(args, arg)
	if handler.Request == nil {
		return args, nil
	}
	_codec := encoding.GetCodec(protocol)
	if _codec == nil {
		return nil, errors.BadRequest("micro.server", "codec not found: '%s'", protocol)
	}

	if handler.BodyValidator != nil {
		result, err := handler.BodyValidator.Validate(gojsonschema.NewBytesLoader(body))
		if err != nil {
			return nil, err
		}
		if !result.Valid() {
			msg := fmt.Sprintf("Validate request body failed")
			for _, desc := range result.Errors() {
				msg = fmt.Sprintf("%s %s", msg, desc)
			}
			return nil, errors.BadRequest("micro.server", msg)
		}
	}

	if handler.Request == utils.TypeOfBytes {
		arg = reflect.Zero(handler.Request)
	} else {
		arg = reflect.New(handler.Request.Elem())
	}
	if err := _codec.Unmarshal(body, arg.Interface()); err != nil {
		return nil, errors.BadRequest("micro.server", "codec unmarshal failed: %s", err.Error())
	}
	args = append(args, arg)
	return args, nil
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
		matches := curdPrefix.FindStringSubmatch(handler.Method.Name)
		if matches == nil {
			path = ""
			return
		}
		method = matches[1]
		resource = fmt.Sprintf("%s/%s", handler.Resource, handler.Name)
		path = fmt.Sprintf("/%s", resource)
	}
	return
}

func isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}

	// 入参数个数 2 3 4
	if mt.NumIn() < 2 || mt.NumIn() > 4 {
		return false
	}

	// 返回参数必须是0个或者2个
	if mt.NumOut() != 0 && mt.NumOut() != 2 {
		return false
	}

	// 第一个参数必须是context(0号参数是实例本身)
	if t1 := mt.In(1); !t1.Implements(utils.TypeOfContext) {
		return false
	}
	// 第二个参数必须是指针类型
	if mt.NumIn() >= 3 && mt.In(2).Kind() != reflect.Ptr {
		return false
	}
	// 第三个参数必须是指针类型或者bytes或者函数
	if mt.NumIn() == 4 && mt.In(3).Kind() != reflect.Ptr && mt.In(3).Kind() != reflect.Func &&
		mt.In(3) != utils.TypeOfBytes {
		return false
	}

	if mt.NumOut() == 2 {
		// 流式传输不允许返回值
		if mt.NumIn() == 4 && mt.In(3).Kind() == reflect.Func {
			return false
		}
		// 返回参数必须是 prt/bytes, error结构
		if mt.Out(1) != utils.TypeOfError || mt.Out(0) != utils.TypeOfBytes && mt.Out(0).Kind() != reflect.Ptr &&
			mt.Out(0).Kind() != reflect.Slice {
			return false
		}
	}
	return true
}

func ExtractComponent(component micro.Component) (map[string]*Handler, []*Handler) {
	typ := reflect.TypeOf(component)
	methods := make(map[string]*Handler)
	var handlers []*Handler

	rtype := typ.Elem()

	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if !isHandlerMethod(method) {
			continue
		}
		mt := method.Type

		metadata := make(map[string]string)
		handler := &Handler{
			Resource:   component.Name(),
			Collection: component.Collection(),
			Method:     method,
			Rtype:      rtype,
		}
		if mt.NumIn() >= 3 {
			query := mt.In(2)
			if query.Elem().NumField() > 0 {
				handler.Query = query
				buff, _ := jsonschema.Marshal(handler.Query, false)
				loader := gojsonschema.NewBytesLoader(buff)
				validator, _ := gojsonschema.NewSchema(loader)
				handler.QueryValidator = validator
			}
		}
		if mt.NumIn() == 4 {
			handler.Request = mt.In(3)
			if handler.Request == utils.TypeOfBytes {
				metadata["req"] = "bytes"
			} else if handler.Request.Kind() == reflect.Func {
				// TODO 流式rpc支持
				panic("no support stream rpc")
				//metadata["req"] = "stream"
			} else {
				// 生成Validator
				buff, _ := jsonschema.Marshal(handler.Request, true)
				loader := gojsonschema.NewBytesLoader(buff)
				validator, _ := gojsonschema.NewSchema(loader)
				handler.BodyValidator = validator
				metadata["req"] = "json"
			}
		}
		if mt.NumOut() == 2 {
			handler.Response = mt.Out(0)
			if handler.Response == utils.TypeOfBytes {
				metadata["res"] = "bytes"
			} else if handler.Response.Implements(utils.TypeOfProtoMsg) {
				metadata["res"] = "proto"
			} else {
				metadata["res"] = "json"
			}
		}
		handler.Metadata = metadata
		switch method.Name {
		case "Get", "List", "Create", "Update", "Patch", "Delete": // restful curd接口
			handler.Name = method.Name
			handlers = append(handlers, handler)
		default:
			matches := curdPrefix.FindStringSubmatch(method.Name)
			if matches == nil { // 没有前缀,网关接口
				name := strings.ToLower(method.Name)
				handler.Name = name
				methods[name] = handler
			} else {
				name := strings.ToLower(matches[2])
				if _, ok := methods[name]; ok {
					panic(fmt.Sprintf("duplicate name %s.%s", component.Name(), name))
				}
				handler.Name = name
				if matches[1] == "RPC" { // 内部rpc接口
					handler.Internal = true
					methods[name] = handler
				} else { // 非标准请求接口
					handlers = append(handlers, handler)
				}
			}
		}
	}

	return methods, handlers
}

func ExtractComponents(components []micro.Component) (map[string]map[string]*Handler, []*Handler) {
	services := make(map[string]map[string]*Handler)
	var handlers []*Handler
	for _, c := range components {
		value := reflect.ValueOf(c)
		name := reflect.Indirect(value).Type().Name()
		if !isExported(name) {
			continue
		}
		methods, hs := ExtractComponent(c)
		handlers = append(handlers, hs...)
		for method, handler := range methods {
			handler.Receiver = value
			m, ok := services[c.Name()]
			if !ok {
				m = make(map[string]*Handler)
				services[c.Name()] = m
			}
			m[method] = handler
		}
	}
	return services, handlers
}

func extractEndpoints(services map[string]map[string]*Handler) (endpoints []*micro.Endpoint) {
	for name, service := range services {
		for method, handler := range service {
			endpoint := &micro.Endpoint{
				Name:     fmt.Sprintf("%s.%s", name, method),
				Metadata: handler.Metadata,
				Internal: handler.Internal,
			}
			endpoints = append(endpoints, endpoint)
		}
	}
	return

}
