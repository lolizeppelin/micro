package utils

import (
	"sync"
)

type Versioned interface {
	Version() int64
}

type VersionedMap[K comparable, V Versioned] struct {
	lock     sync.RWMutex
	storages map[K]V
}

func (m *VersionedMap[K, V]) LoadAndDelete(key K) (V, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.storages[key]
	if !ok {
		return v, false
	}
	delete(m.storages, key)
	return v, true
}

func (m *VersionedMap[K, V]) Load(key K) (V, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	v, ok := m.storages[key]
	return v, ok
}

func (m *VersionedMap[K, V]) Store(key K, value V) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.storages[key] = value
}

// LoadOrStore 存储一个值、find表示值存在,没有进行存储
func (m *VersionedMap[K, V]) LoadOrStore(key K, value V) (ret V, find bool) {
	m.lock.RLock()
	v, ok := m.storages[key]
	m.lock.RUnlock()
	if ok {
		return v, ok
	}
	m.lock.Lock()
	m.lock.Unlock()
	v, ok = m.storages[key]
	if ok {
		return v, ok
	}
	m.storages[key] = value
	return value, false
}

// StoreNew 存储一个新值、返回值表示是否存储成功
func (m *VersionedMap[K, V]) StoreNew(key K, value V) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.storages[key]
	if !ok {
		m.storages[key] = value
		return true
	}
	if v == nil || value.Version() > v.Version() {
		m.storages[key] = value
		return true
	}
	return false
}

func (m *VersionedMap[K, V]) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.storages)
}

func (m *VersionedMap[K, V]) Values() []V {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return MapValues(m.storages)
}

func (m *VersionedMap[K, V]) Keys() []K {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return MapKeys(m.storages)
}

func (m *VersionedMap[K, V]) Clone() map[K]V {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return CopyMap(m.storages)
}

func (m *VersionedMap[K, V]) Clear() {
	m.lock.RLock()
	defer m.lock.RUnlock()
	MapClear(m.storages)
}

func NewVersionedMap[K comparable, V Versioned]() *VersionedMap[K, V] {
	m := &VersionedMap[K, V]{
		storages: map[K]V{},
	}
	return m
}
