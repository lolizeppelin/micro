package utils

import (
	"github.com/swaggest/jsonschema-go"
	"reflect"
)

// BuildJsonSchema 从结构体生成json schema
func BuildJsonSchema(value reflect.Type) ([]byte, error) {
	r := jsonschema.Reflector{}

	s, _ := r.Reflect(reflect.New(value.Elem()).Interface(), jsonschema.InlineRefs)
	return s.MarshalJSON()
}
