package jsonschema

import (
	"encoding/json"
	"errors"
	"github.com/lolizeppelin/micro/utils"
	"github.com/lolizeppelin/micro/utils/ast"
	"github.com/swaggest/jsonschema-go"
	"reflect"
	"strings"
)

func gettype(typ *jsonschema.Type) jsonschema.SimpleType {
	if typ == nil {
		return ""
	}
	if typ.SimpleTypes == nil {
		if len(typ.SliceOfSimpleTypeValues) > 0 {
			for _, _t := range typ.SliceOfSimpleTypeValues {
				if _t != jsonschema.Null {
					return _t
				}

			}
		}
	} else {
		return *typ.SimpleTypes
	}
	return ""
}

func WithAdditionalProperties(schema *jsonschema.Schema) {
	if schema.AdditionalProperties != nil {
		return
	}
	typ := gettype(schema.Type)
	if typ == jsonschema.Object {
		_false := new(bool)
		*_false = false
		schema.WithAdditionalProperties(*&jsonschema.SchemaOrBool{TypeBoolean: _false})
		for _, s := range schema.Properties {
			WithAdditionalProperties(s.TypeObject)
		}
	} else if typ == jsonschema.Array {
		WithAdditionalProperties(schema.Items.SchemaOrBool.TypeObject)
	}

}

func Marshal(value reflect.Type, additional bool) ([]byte, error) {
	r := jsonschema.Reflector{}
	if value == utils.TypeOfBytes {
		s, _ := r.Reflect([]byte(nil), jsonschema.InlineRefs)
		return s.MarshalJSON()
	}

	if value.Kind() == reflect.Slice {
		iType := value.Elem()

		var es jsonschema.Schema
		var err error
		switch iType.Kind() {
		case reflect.Pointer:
			es, err = r.Reflect(reflect.New(iType.Elem()).Interface(), jsonschema.InlineRefs)
			break
		case reflect.Struct:
			es, err = r.Reflect(reflect.New(iType).Interface(), jsonschema.InlineRefs)
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
			es, err = r.Reflect(reflect.Zero(iType).Interface(), jsonschema.InlineRefs)
			break
		default:
			err = errors.New("error jsonschema type")
		}
		if err != nil {
			return nil, err
		}

		p := &es
		if !additional {
			WithAdditionalProperties(&es)
		}
		m := map[string]any{
			"type":        "array",
			"items":       p,
			"description": "no description",
		}
		return json.Marshal(m)
	}
	s, _ := r.Reflect(reflect.New(value.Elem()).Interface(), jsonschema.InlineRefs)
	if !additional {
		WithAdditionalProperties(&s)
	}
	return s.MarshalJSON()
}

func Schema(target reflect.Type, additional bool) (map[string]any, error) {
	b, _ := Marshal(target, additional)
	m := map[string]any{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetComment(component reflect.Type, method reflect.Method) *Comment {
	comment, _ := ast.GetComment(component, method)
	if comment == "" {
		return nil
	}
	comments := strings.Split(comment, "\n")
	if len(comments) == 0 {
		return nil
	}
	return &Comment{
		Summary:     strings.TrimSpace(comments[0]),
		Description: strings.TrimSpace(strings.Join(comments[1:], "\n")),
	}
}

// General 通用api
func General(path, name string, metadata map[string]string,
	query reflect.Type, req reflect.Type, res reflect.Type,
	component reflect.Type, method reflect.Method, comment *Comment) (api APIPath, err error) {

	p := APIPath{
		Path: path,
		Request: &Request{
			Method: name,
		},
		Responses: []*Response{
			{
				Code:        200,
				Description: "success",
			},
		},
	}
	if comment == nil {
		p.Comment = GetComment(component, method)
	} else {
		p.Comment = &Comment{Summary: comment.Summary, Description: comment.Description}
	}
	var m map[string]any
	if query != nil {
		m, err = Schema(query, true)
		if err != nil {
			return
		}
		p.Request.Query = m
	}
	if req != nil {
		m, err = Schema(req, false)
		if err != nil {
			return
		}
		p.Request.Body = ContentBody{
			Type:   metadata["req"],
			Schema: m,
		}
	}
	if res != nil {
		m, err = Schema(res, false)
		if err != nil {
			return
		}
		p.Responses[0].Body = &ContentBody{
			Type:   metadata["res"],
			Schema: m,
		}
	}
	api = p
	return
}
