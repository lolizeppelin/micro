package server

import (
	"fmt"
	"github.com/lolizeppelin/micro"
	"github.com/lolizeppelin/micro/utils"
	"github.com/lolizeppelin/micro/utils/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
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

	//if mt.NumIn() == 4 && mt.In(3).Kind() != reflect.Ptr && mt.In(3).Kind() != reflect.Func &&
	//	mt.In(3) != utils.TypeOfBytes {
	//	return false
	//}

	// 第三个参数必须是指针类型或者bytes或者函数或指针列表
	if mt.NumIn() == 4 && !(mt.In(3).Kind() == reflect.Ptr ||
		mt.In(3).Kind() == reflect.Func ||
		mt.In(3) == utils.TypeOfBytes ||
		(mt.In(3).Kind() == reflect.Slice && mt.In(3).Elem().Kind() == reflect.Ptr)) {
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

/*
ExtractComponent 解析组件
返回的map是rpcserver可用穿透的Handlers（可用网关代理）
返回的列表是所有Handlers
*/
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
			Hooks:      component.Hooks(method.Name),
		}
		if mt.NumIn() >= 3 {
			query := mt.In(2)
			handler.Query = query
			buff, _ := jsonschema.Marshal(handler.Query, false)
			loader := gojsonschema.NewBytesLoader(buff)
			validator, _ := gojsonschema.NewSchema(loader)
			handler.QueryValidator = validator
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
		handlers = append(handlers, handler)

		handler.Metadata = metadata
		switch method.Name {
		case "Get", "List", "Create", "Update", "Patch", "Delete": // restful curd接口
			handler.Name = method.Name
		default:
			matches := curdPrefix.FindStringSubmatch(method.Name)
			if matches == nil { // 没有前缀,网关接口
				name := strings.ToLower(method.Name)
				if _, ok := methods[name]; ok {
					panic(fmt.Sprintf("duplicate name %s.%s", component.Name(), name))
				}
				handler.Name = name
				methods[name] = handler
			} else {
				name := strings.ToLower(matches[2])
				if matches[1] == "RPC" { // 内部rpc接口
					if _, ok := methods[name]; ok {
						panic(fmt.Sprintf("duplicate name %s.%s", component.Name(), name))
					}
					handler.Internal = true
					methods[name] = handler
				}
				handler.Name = name
			}
		}
	}

	return methods, handlers
}

/*
ExtractComponents 解析组件
返回的map是rpcserver可用穿透的Handlers（可用网关代理）
返回的列表是所有Handlers
*/
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
		for _, h := range hs {
			h.Receiver = value
		}
		handlers = append(handlers, hs...)
		for method, handler := range methods {
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
				Name:       fmt.Sprintf("%s.%s", name, method),
				Metadata:   handler.Metadata,
				PrimaryKey: handler.Name == "Get" || handler.Name == "Update" || handler.Name == "Delete",
				Internal:   handler.Internal,
			}
			endpoints = append(endpoints, endpoint)
		}
	}
	return

}
