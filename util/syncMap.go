package util

import (
	"fmt"
	"sync"
)

type SyncMap[K comparable, V any] struct {
	Lock sync.RWMutex
	M    map[K]V
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{M: make(map[K]V)}
}

func (m *SyncMap[K, V]) Load(key K) (V, bool) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	value, ok := m.M[key]
	return value, ok
}

func (m *SyncMap[K, V]) LoadAnd(key K, f func(V) error) error {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	value, ok := m.M[key]
	if !ok {
		return fmt.Errorf("")
	}
	return f(value)
}

func (m *SyncMap[K, V]) Modify(key K, f func(V) (V, error)) error {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	value, ok := m.M[key]
	if !ok {
		return fmt.Errorf("")
	}
	value, err := f(value)
	if err != nil {
		return err
	}
	m.M[key] = value
	return nil
}
