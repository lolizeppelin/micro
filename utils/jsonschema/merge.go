package jsonschema

import (
	"encoding/json"
	"fmt"
)

func mergeJSONSchema(dst, src map[string]any) error {
	var err error
	for key, srcVal := range src {
		if _, exists := dst[key]; exists {
			switch key {
			case "required":
				err = mergeRequired(dst, srcVal)
				if err != nil {
					return err
				}
			case "properties":
				err = mergeProperties(dst, srcVal)
				if err != nil {
					return err
				}
			case "additionalProperties":
				err = mergeAdditionalProperties(dst, srcVal)
				if err != nil {
					return err
				}
			case "allOf", "anyOf", "oneOf":
				err = mergeCombinedSchemas(dst, key, srcVal)
				if err != nil {
					return err
				}
			default:
				dst[key] = srcVal
			}
		} else {
			dst[key] = srcVal
		}
	}
	return nil
}

func mergeRequired(dst map[string]any, srcVal any) error {
	dstRequired, ok1 := dst["required"].([]string)
	if !ok1 {
		return fmt.Errorf("jsonschema required type error")
	}
	srcRequired, ok2 := srcVal.([]string)
	if !ok2 {
		return fmt.Errorf("jsonschema required type error")
	}
	seen := make(map[string]bool)
	var combined []any
	for _, v := range dstRequired {
		if !seen[v] {
			seen[v] = true
			combined = append(combined, v)
		}
	}
	for _, v := range srcRequired {
		if !seen[v] {
			seen[v] = true
			combined = append(combined, v)
		}
	}
	dst["required"] = combined
	return nil
}

func mergeProperties(dst map[string]any, srcVal any) error {
	dstProps, ok1 := dst["properties"].(map[string]any)
	if !ok1 {
		return fmt.Errorf("jsonschema properties type error")
	}
	srcProps, ok2 := srcVal.(map[string]any)
	if !ok2 {
		return fmt.Errorf("jsonschema properties type error")
	}
	var err error
	for propKey, srcProp := range srcProps {
		srcPropMap, ok4 := srcProp.(map[string]any)
		if !ok4 {
			return fmt.Errorf("jsonschema properties type error")
		}
		if dstProp, exists := dstProps[propKey]; exists {
			dstPropMap, ok3 := dstProp.(map[string]any)
			if !ok3 {
				return fmt.Errorf("jsonschema properties type error")
			}
			err = mergeJSONSchema(dstPropMap, srcPropMap)
			if err != nil {
				return err
			}
			dstProps[propKey] = dstPropMap
		} else {
			dstProps[propKey] = srcPropMap
		}
	}
	return nil
}

func mergeAdditionalProperties(dst map[string]any, srcVal any) error {
	dstAP := dst["additionalProperties"]
	var err error
	if isSchemaObject(dstAP) && isSchemaObject(srcVal) {
		err = mergeJSONSchema(dstAP.(map[string]any), srcVal.(map[string]any))
		if err != nil {
			return err
		}
	} else {
		dst["additionalProperties"] = srcVal
	}
	return nil
}

func mergeCombinedSchemas(dst map[string]any, key string, srcVal any) error {
	dstCombined, ok1 := dst[key].([]map[string]any)
	if !ok1 {
		return fmt.Errorf("jsonschema %s type error", key)
	}
	srcCombined, ok2 := srcVal.([]map[string]any)
	if !ok2 {
		return fmt.Errorf("jsonschema %s type error", key)
	}
	dst[key] = append(dstCombined, srcCombined...)
	return nil
}

func isSchemaObject(v any) bool {
	_, ok := v.(map[string]any)
	return ok
}

/*
MergeJSONSchema 合并jsonschema, b覆盖a
*/
func MergeJSONSchema(a, b map[string]any) (map[string]any, error) {
	buff, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	output := new(map[string]any)
	err = json.Unmarshal(buff, output)
	if err != nil {
		return nil, err
	}
	dst := *output
	err = mergeJSONSchema(dst, b)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
