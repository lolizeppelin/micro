package jsonschema

import (
	"encoding/json"
	"github.com/swaggest/jsonschema-go"
	"net/http"
	"reflect"
)

func Marshal(value reflect.Type) ([]byte, error) {
	r := jsonschema.Reflector{}

	s, _ := r.Reflect(reflect.New(value.Elem()).Interface(), jsonschema.InlineRefs)
	return s.MarshalJSON()
}

func Schema(target reflect.Type) (map[string]any, error) {
	b, _ := Marshal(target)
	m := map[string]any{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// General 通用api
func General(path string, metadata map[string]string,
	query reflect.Type, req reflect.Type, res reflect.Type) (api APIPath, err error) {

	p := APIPath{
		Path: path,
		Request: &Request{
			Method: http.MethodPost,
		},
		Responses: []*Response{
			{
				Code:        200,
				Description: "success",
			},
		},
	}
	var m map[string]any
	if query != nil {
		m, err = Schema(query)
		if err != nil {
			return
		}
		p.Request.Query = m
	}
	if req != nil {
		m, err = Schema(req)
		if err != nil {
			return
		}
		p.Request.Body = ContentBody{
			Type:   metadata["res"],
			Schema: m,
		}
	}
	if res != nil {
		m, err = Schema(res)
		if err != nil {
			return
		}
		p.Responses[0].Body = &ContentBody{
			Type:   metadata["req"],
			Schema: m,
		}
	}
	api = p

	return
}
