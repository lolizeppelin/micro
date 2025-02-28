package utils

import (
	"fmt"
	"golang.org/x/exp/maps"
	"reflect"
)

// MapToSliceByField 提取结构体中的自定字段转为列表
func MapToSliceByField[T any, K comparable](s []T, field ...string) (map[K]T, error) {
	result := make(map[K]T)

	key := "ID"
	if len(field) > 0 {
		key = field[0]
	}

	for _, item := range s {
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// 获取字段值
		fieldVal := val.FieldByName(key)
		if !fieldVal.IsValid() {
			return nil, fmt.Errorf("field %s not found", key)
		}

		// 确保字段可以比较
		if !fieldVal.Type().Comparable() {
			return nil, fmt.Errorf("field %s is not comparable", key)
		}

		// 转换字段值为键类型
		k, ok := fieldVal.Interface().(K)
		if !ok {
			return nil, fmt.Errorf("field %v cannot be converted to key type", k)
		}

		result[k] = item
	}

	return result, nil
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return new(SyncMap[K, V])
}

func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	return maps.Keys(m)
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	return maps.Values(m)
}

func MapClear[M ~map[K]V, K comparable, V any](m M) {
	maps.Clear(m)
}
