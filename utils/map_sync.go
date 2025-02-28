package utils

import (
	"sync"
)

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
