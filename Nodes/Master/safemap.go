package main

import "sync"

type SafeMap[K comparable, V any] struct {
	mu    sync.RWMutex
	nodes map[K]V
}

func NewSaveMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		nodes: make(map[K]V),
	}
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes[key] = value
}
func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.nodes[key]
	return val, ok

}
