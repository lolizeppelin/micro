package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
)

// CopySlice 拷贝列表
func CopySlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

// DeepCopyMap 通过序列号深拷贝
func DeepCopyMap[T any](src map[string]T) (map[string]T, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	if err := enc.Encode(src); err != nil {
		return nil, err
	}
	dst := new(map[string]T)
	if err := dec.Decode(&dst); err != nil {
		return nil, err
	}
	return *dst, nil
}

// MergeMaps 通过序列号深拷贝后合并
func MergeMaps[T any](sources ...map[string]T) (map[string]T, error) {
	dst := make(map[string]T)
	for _, s := range sources {
		tmp, err := DeepCopyMap(s)
		if err != nil {
			return nil, err
		}
		for k, v := range tmp {
			dst[k] = v
		}
	}
	return dst, nil
}

// MergeSliceMaps 通过序列号深拷贝后合并
func MergeSliceMaps[T any](key string, sources ...any) (map[string]T, error) {
	tmp := make([]map[string]T, 0)
	for _, s := range sources {
		val := reflect.ValueOf(s)
		// Check if the source is a map
		if val.Kind() != reflect.Map {
			return nil, fmt.Errorf("each source must be a map")
		}
		// Check if the map's key is a string
		if val.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map keys must be strings")
		}
		// Create a new map for this element
		newMap := make(map[string]T)
		// Iterate over the map's entries
		for _, mapKey := range val.MapKeys() {
			if mapKey.String() == key {
				// Get the value associated with the key
				val := val.MapIndex(mapKey)
				// Convert the value to the desired type
				typedVal, ok := val.Interface().(T)
				if !ok {
					return nil, fmt.Errorf("map value type mismatch")
				}
				// Add the key-value pair to the new map
				newMap[mapKey.String()] = typedVal
			}
		}

		// Append the new map to the temporary slice
		tmp = append(tmp, newMap)

	}
	return MergeMaps(tmp...)
}
