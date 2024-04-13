package server

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lolizeppelin/micro"
	"google.golang.org/grpc/encoding"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	typeOfError    = reflect.TypeOf((*error)(nil)).Elem()
	typeOfBytes    = reflect.TypeOf(([]byte)(nil))
	typeOfContext  = reflect.TypeOf(new(context.Context)).Elem()
	typeOfProtoMsg = reflect.TypeOf(new(proto.Message)).Elem()
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
	Internal bool
	Name     string            // method name
	Receiver reflect.Value     // receiver of method
	Method   reflect.Method    // method stub
	Request  reflect.Type      // 请求参数
	Response reflect.Type      // 返回参数
	Metadata map[string]string // 元数据
}

func (handler *Handler) BuildArgs(ctx context.Context, protocol string, body []byte) ([]reflect.Value, error) {
	args := []reflect.Value{handler.Receiver, reflect.ValueOf(ctx)}
	if handler.Request == nil {
		return args, nil
	}
	codec := encoding.GetCodec(protocol)
	if codec == nil {
		return nil, fmt.Errorf("content type codec not found")
	}
	arg := reflect.New(handler.Request.Elem())
	if err := codec.Unmarshal(body, arg.Interface()); err != nil {
		return nil, err
	}
	args = append(args, arg)
	return args, nil
}

func (handler *Handler) Match(protocol, accept string) bool {
	return protocol == handler.Metadata["res"] && accept == handler.Metadata["req"]
}

func isHandlerMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}
	// 参数必须是2个或3个
	if mt.NumIn() != 2 && mt.NumIn() != 3 {
		return false
	}
	// 返回参数必须是0个或者2个
	if mt.NumOut() != 0 && mt.NumOut() != 2 {
		return false
	}
	// 第一个参数必须是context(0号参数是实例本身)
	if t1 := mt.In(1); !t1.Implements(typeOfContext) {
		return false
	}
	// 第二个参数必须是指针类型或者bytes或者函数
	if mt.NumIn() == 3 && mt.In(2).Kind() != reflect.Ptr && mt.In(2).Kind() != reflect.Func && mt.In(2) != typeOfBytes {
		return false
	}
	if mt.NumOut() == 2 {
		// 流式传输不允许返回值
		if mt.NumIn() == 3 && mt.In(2).Kind() == reflect.Func {
			return false
		}
		// 返回参数必须是 prt/bytes, error结构
		if mt.Out(1) != typeOfError || mt.Out(0) != typeOfBytes && mt.Out(0).Kind() != reflect.Ptr {
			return false
		}
	}
	return true
}

func extractComponent(component micro.Component) map[string]*Handler {
	typ := reflect.TypeOf(component)
	methods := make(map[string]*Handler)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mt := method.Type
		name := strings.ToLower(method.Name)
		if isHandlerMethod(method) {
			metadata := make(map[string]string)
			handler := &Handler{
				Method: method,
			}
			if mt.NumIn() == 3 {
				handler.Request = mt.In(2)
				if handler.Request == typeOfBytes {
					metadata["req"] = "bytes"
				} else if handler.Request.Implements(typeOfProtoMsg) {
					metadata["req"] = "proto"
				} else if handler.Request.Kind() == reflect.Func {
					// TODO 流式rpc支持
					panic("no support stream rpc")
					//metadata["req"] = "stream"
				} else {
					metadata["req"] = "json"
				}
			}
			if mt.NumOut() == 2 {
				handler.Response = mt.Out(2)
				if handler.Response == typeOfBytes {
					metadata["res"] = "bytes"
				} else if handler.Response.Implements(typeOfProtoMsg) {
					metadata["res"] = "proto"
				} else {
					metadata["res"] = "json"
				}
			}
			handler.Metadata = metadata
			methods[name] = handler
		}
	}
	return methods
}

func ExtractComponents(components []micro.Component) map[string]map[string]*Handler {
	services := make(map[string]map[string]*Handler)
	for _, c := range components {
		value := reflect.ValueOf(c)
		name := reflect.Indirect(value).Type().Name()
		if !isExported(name) {
			continue
		}

		// 服务名默认使用结构体名的小写
		field := value.Elem().FieldByName("Name")
		if field.IsValid() && field.Kind() != reflect.String {
			name = field.String()
		} else {
			name = strings.ToLower(name)
		}

		// 方法默认内部接口
		internal := true
		field = value.Elem().FieldByName("Internal")
		if field.IsValid() && field.Kind() != reflect.Bool {
			internal = field.Bool()
		}

		methods := extractComponent(c)
		for method, handler := range methods {
			handler.Receiver = value
			handler.Internal = internal
			m, ok := services[name]
			if !ok {
				m = make(map[string]*Handler)
				services[name] = m
			}
			m[method] = handler
		}
	}
	return services
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
