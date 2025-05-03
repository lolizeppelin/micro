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

func MergeJsonMaps(dest, src map[string]any) map[string]any {
	for key, srcVal := range src {
		// 检查目标map是否已存在该key
		if destVal, exists := dest[key]; exists {
			// 类型断言判断是否为嵌套map
			destMap, destIsMap := destVal.(map[string]any)
			srcMap, srcIsMap := srcVal.(map[string]any)
			// 如果双方都是map，则递归合并
			if destIsMap && srcIsMap {
				dest[key] = MergeJsonMaps(destMap, srcMap)
			} else {
				// 否则直接覆盖（包括类型不同的情况）
				dest[key] = srcVal
			}
		} else {
			// 目标map不存在该key时直接添加
			dest[key] = srcVal
		}
	}
	return dest
}
