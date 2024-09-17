package utils

import "sync"

type KeyLocks struct {
	mutexes *SyncMap[string, *sync.RWMutex]
}

func (m *KeyLocks) Get(key string) *sync.RWMutex {
	lock, ok := m.mutexes.Load(key)
	if !ok {
		lock, _ = m.mutexes.LoadOrStore(key, new(sync.RWMutex))
	}
	return lock
}

func (m *KeyLocks) Lock(key string) {
	m.Get(key).Lock()
}

func (m *KeyLocks) UnLock(key string) {
	m.Get(key).Unlock()
}

func (m *KeyLocks) RLock(key string) {
	m.Get(key).RLock()
}

func (m *KeyLocks) UnRLock(key string) {
	m.Get(key).Unlock()
}

func NewKeyLocks() *KeyLocks {
	return &KeyLocks{
		mutexes: NewSyncMap[string, *sync.RWMutex](),
	}
}
