package utils

import "sync"

// sync.Map with type enforcement
type Map[K comparable, V any] struct{ m sync.Map }

func (m *Map[K, V]) Load(key K) (value V) {
	v, ok := m.m.Load(key)
	if !ok {
		return value
	}
	return v.(V)
}

func (m *Map[K, V]) Store(key K, value V) { m.m.Store(key, value) }

func (m *Map[K, V]) LoadOrStore(key K, value V) V {
	a, _ := m.m.LoadOrStore(key, value)
	return a.(V)
}

func (m *Map[K, V]) LoadAndStore(key K, value V) V {
	v := m.Load(key)
	m.Store(key, value)
	return v
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V) {
	v, loaded := m.m.LoadAndDelete(key)
	if !loaded {
		return value
	}
	return v.(V)
}

func (m *Map[K, V]) Delete(key K) { m.m.Delete(key) }

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool { return f(key.(K), value.(V)) })
}
