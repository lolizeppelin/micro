package utils

import (
	"fmt"
	"golang.org/x/exp/maps"
	"reflect"
	"sync"
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

type SyncMap[K comparable, V any] struct {
	storages sync.Map
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.storages.Delete(key)
}

func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.storages.Load(key)
	if !ok {
		return value, ok
	}
	return v.(V), ok
}

func (m *SyncMap[K, V]) Len() int {
	var i int
	m.storages.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}

func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.storages.LoadAndDelete(key)
	if !loaded {
		return value, loaded
	}
	return v.(V), loaded
}

func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	a, loaded := m.storages.LoadOrStore(key, value)
	return a.(V), loaded
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.storages.Range(func(key, value any) bool { return f(key.(K), value.(V)) })
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.storages.Store(key, value)
}

func (m *SyncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	a, loaded := m.storages.Swap(key, value)
	return a.(V), loaded
}

func (m *SyncMap[K, V]) All() map[K]V {
	o := map[K]V{}
	m.storages.Range(func(key, value any) bool {
		o[key.(K)] = value.(V)
		return true
	})
	return o
}

func (m *SyncMap[K, V]) DeleteAll() {
	var keys []K
	m.storages.Range(func(key, value any) bool {
		keys = append(keys, key.(K))
		return true
	})
	for _, key := range keys {
		m.storages.Delete(key)
	}
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
