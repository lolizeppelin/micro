package utils

import (
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	"reflect"
)

func _formatSchema(node *jsonschema.Schema, definitions map[string]*jsonschema.Schema) *jsonschema.Schema {
	if node.Type == "array" {
		if node.Items.Ref != "" {
			node.Items = definitions[node.Items.Ref]
			return node
		}
	} else if node.Type == "object" {
		keys := map[string]*jsonschema.Schema{}
		for pair := node.Properties.Newest(); pair != nil; pair = pair.Prev() {
			n := _formatSchema(pair.Value, definitions)
			if n != pair.Value {
				keys[pair.Key] = n
			}
		}
		for k, v := range keys {
			node.Properties.Set(k, v)
		}
	} else if node.Ref != "" {
		return definitions[node.Ref]
	}
	return node
}

func BuildJsonSchema(value reflect.Type) ([]byte, error) {
	s := jsonschema.Reflect(reflect.New(value.Elem()).Interface())
	definitions := map[string]*jsonschema.Schema{}
	for k, v := range s.Definitions {
		key := fmt.Sprintf("#/$defs/%s", k)
		definitions[key] = v
	}
	for name, def := range definitions {
		definitions[name] = _formatSchema(def, definitions)
	}
	root := definitions[s.Ref]
	_formatSchema(root, definitions)
	buff, err := json.Marshal(root)
	return buff, err
}
