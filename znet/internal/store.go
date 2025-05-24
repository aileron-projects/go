package internal

import (
	"errors"
	"sync"
)

type ComparableCloser interface {
	comparable
	Close() error
}

// CloserStore stores [ComparableCloser].
// All implemented methods are safe for concurrent call.
type CloserStore[T ComparableCloser] struct {
	mu      sync.Mutex
	closers map[T]struct{}
}

// CloseAll call Close() method of all stored closers
// and delete from the store after closed.
func (m *CloserStore[T]) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var errs []error
	for closer := range m.closers {
		delete(m.closers, closer)
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Store stores the given closer to the store.
// It replaces the existing closer if the same
// closer is already exist.
func (m *CloserStore[T]) Store(closer T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closers == nil {
		m.closers = map[T]struct{}{}
	}
	m.closers[closer] = struct{}{}
}

// Delete deletes the closer from the store.
func (m *CloserStore[T]) Delete(closer T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.closers, closer)
}

// Length returns the number of closers.
func (m *CloserStore[T]) Length() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.closers)
}
