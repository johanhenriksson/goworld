package util

import (
	"sync"
)

// Type-safe sync.Map implementation
// Read sync.Map documentation for caveats
type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: sync.Map{},
	}
}

func (m *SyncMap[K, V]) Load(key K) (value V, exists bool) {
	var v any
	v, exists = m.m.Load(key)
	if exists {
		value = v.(V)
	}
	return
}

func (m *SyncMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}
