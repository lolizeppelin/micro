package utils

import (
	"fmt"
	"reflect"
)

func GetField(object any, key string) (any, error) {
	elem := reflect.ValueOf(object)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("provided object is not a struct")
	}

	field := elem.FieldByName(key)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s not found", key)
	}
	if !field.CanInterface() {
		return nil, fmt.Errorf("field %s cannot be accessed", key)
	}
	return field.Interface(), nil
}

func GetStringsByField[T any](object []T, key string) ([]string, error) {
	var strings []string
	for _, elem := range object {
		v, err := GetField(elem, key)
		if err != nil {
			return nil, err
		}
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("field %s is not a string", key)
		}
		strings = append(strings, s)
	}
	return strings, nil
}

func UnSafeGetStringsByField[T any](object []T, key string) []string {
	strings, _ := GetStringsByField[T](object, key)
	return strings
}
