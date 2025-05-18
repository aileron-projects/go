package internal

import (
	"iter"
	"sync"
)

// UniqueStore stores a comparable values.
// All implemented methods are safe for concurrent call.
type UniqueStore[T comparable] struct {
	mu     sync.Mutex
	values map[T]struct{}
}

// Value returns an iterator for values.
func (m *UniqueStore[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range m.values {
			if !yield(v) {
				return
			}
		}
	}
}

// Set sets the given value to the store.
// The value overwrites the existing value if exists.
func (m *UniqueStore[T]) Set(value T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.values == nil {
		m.values = map[T]struct{}{}
	}
	m.values[value] = struct{}{}
}

// Delete deletes the value from the store.
func (m *UniqueStore[T]) Delete(value T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.values, value)
}

// Length returns the number of values.
func (m *UniqueStore[T]) Length() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.values)
}
