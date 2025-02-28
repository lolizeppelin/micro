package utils

import (
	"golang.org/x/exp/maps"
	"sync"
	"sync/atomic"
	"time"
)

type ExpiredItem[V any] struct {
	ExpiredAt *int64
	Payload   V
}

type ExpireMap[K comparable, V any] struct {
	lock        sync.RWMutex
	storages    map[K]*ExpiredItem[V]
	fetch       func(...K) (map[K]V, error) // 新值获取方法
	delay       int64                       // 自动延期,默认0,不启用,启用最低延迟10s
	expire      int64                       // 过期时间,默认300s
	placeholder V                           // 默认空值

	pool chan K
}

func (m *ExpireMap[K, V]) Delete(key K) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.storages, key)

}

// Load 加载缓存, 允许过期值
func (m *ExpireMap[K, V]) Load(key K, expired bool) (V, bool) {
	m.lock.RLock()
	item, ok := m.storages[key]
	m.lock.RUnlock()

	now := NowUnix()
	if ok {
		at := atomic.LoadInt64(item.ExpiredAt)
		if at < now { // 已经过期
			if m.delay > 0 { // 允许自动延期
				atomic.CompareAndSwapInt64(item.ExpiredAt, at, at+m.delay)
				m.pool <- key // 触发异步更新
			}
			if !expired { // 不允许过期
				return m.placeholder, false
			}
		}
		return item.Payload, ok
	}
	return m.SyncLoad(key)
}

func (m *ExpireMap[K, V]) SyncLoad(key K) (V, bool) {
	m.lock.RLock()
	item, ok := m.storages[key]
	m.lock.RUnlock()
	now := NowUnix()
	if ok {
		at := atomic.LoadInt64(item.ExpiredAt)
		if at > now {
			return item.Payload, true
		}
	}

	v, err := m.fetch(key)
	if err != nil || len(v) == 0 {
		return m.placeholder, false
	}
	payload, find := v[key]
	if !find {
		return m.placeholder, false
	}
	expired := now + m.expire

	item = &ExpiredItem[V]{
		ExpiredAt: &expired,
		Payload:   payload,
	}
	m.lock.Lock()
	m.storages[key] = item
	m.lock.Unlock()
	return item.Payload, true

}

// LoadOrStore 存储或者加载( ok表示加载成功)
func (m *ExpireMap[K, V]) LoadOrStore(key K, value V, expire int64) (V, bool) {
	m.lock.RLock()
	item, ok := m.storages[key]
	m.lock.RUnlock()
	now := NowUnix()
	if ok {
		if atomic.LoadInt64(item.ExpiredAt) > now { // 未过期
			return item.Payload, ok
		}
	}
	// 已经过期或未加载
	m.lock.Lock()
	defer m.lock.Unlock()
	item, ok = m.storages[key]
	if ok && atomic.LoadInt64(item.ExpiredAt) > now { // 未过期
		return item.Payload, ok
	}
	expired := NowUnix()
	if expire > 0 {
		expired = expired + expire
	} else {
		expired = expired + m.expire
	}
	m.storages[key] = &ExpiredItem[V]{
		ExpiredAt: &expired,
		Payload:   value,
	}
	return value, false
}

func (m *ExpireMap[K, V]) Store(key K, value V, expire int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	expired := NowUnix()

	if expire > 0 {
		expired = expired + expire
	} else {
		expired = expired + m.expire
	}
	m.storages[key] = &ExpiredItem[V]{
		ExpiredAt: &expired,
		Payload:   value,
	}
}

func (m *ExpireMap[K, V]) DeleteAll() {
	m.lock.Lock()
	defer m.lock.Unlock()
	maps.Clear(m.storages)

}

func (m *ExpireMap[K, V]) All() map[K]V {
	m.lock.RLock()
	defer m.lock.RUnlock()
	o := map[K]V{}
	for k, item := range m.storages {
		o[k] = item.Payload
	}
	return o
}

func (m *ExpireMap[K, V]) Len() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.storages)
}

func (m *ExpireMap[K, V]) async() {
	timeout := 15 * time.Second

	pools := map[K]struct{}{}

	flush := func() {
		if len(pools) == 0 {
			return
		}
		values, err := m.fetch(maps.Keys(pools)...)
		maps.Clear(pools)
		if err == nil && len(values) > 0 {
			m.flush(values)
		}
	}

	for {
		select {
		case k := <-m.pool:
			pools[k] = struct{}{}
			if len(pools) >= 100 {
				flush()
			}
		case <-time.After(timeout):
			flush()
		}
	}
}

func (m *ExpireMap[K, V]) flush(values map[K]V) {
	m.lock.Lock()
	defer m.lock.Unlock()
	expired := NowUnix() + m.expire
	for k, v := range values {
		m.storages[k] = &ExpiredItem[V]{
			ExpiredAt: &expired,
			Payload:   v,
		}
	}
}

func NewExpireMap[K comparable, V any](fetch func(...K) (map[K]V, error),
	expire int64, delay ...int64) *ExpireMap[K, V] {
	if fetch == nil {
		panic("fetch function not found")
	}

	async := int64(0)
	if len(delay) > 0 && delay[0] > 10 {
		async = delay[0]
		if async < 10 {
			async = 10
		}
	}
	if expire < 3 {
		expire = 3
	}

	m := &ExpireMap[K, V]{
		fetch:    fetch,
		expire:   expire,
		delay:    async,
		storages: map[K]*ExpiredItem[V]{},
		pool:     make(chan K, 100),
	}

	if async > 0 {
		go m.async()
	}

	return m
}
