package zsync

import (
	"sync"
)

var (
	_ Store[any, any] = &sync.Map{}
)

// Store stores sets of key and values and provide their operation.
// [sync.Map] satisfies the Store[any,any] interface.
type Store[K comparable, T any] interface {
	// Store sets the value for a key.
	// It the key already exists, it replaces the value.
	Store(key K, value T)
	// Load returns the stored value for a key.
	// If the kwy was not found, it returns zero value of type T and false.
	Load(key K) (value T, ok bool)
	// Delete deletes the value for a key.
	Delete(key K)
	// Clear deletes all the entries.
	Clear()
	// Swap swaps the value for a key and returns the previous value if any.
	// The loaded result reports whether the key was present.
	Swap(key K, value T) (previous T, loaded bool)
	// LoadAndDelete deletes the value for a key, returning the previous value if any.
	// The loaded result reports whether the key was present.
	LoadAndDelete(key K) (value T, loaded bool)
	// LoadOrStore returns the existing value for the key if present.
	// Otherwise, it stores and returns the given value.
	// The loaded result is true if the value was loaded, false if stored.
	LoadOrStore(key K, value T) (actual T, loaded bool)
	// CompareAndSwap swaps the old and new values for key
	// if the value stored is equal to old.
	// The old value must be of a comparable type.
	CompareAndSwap(key K, old, new T) (swapped bool)
	// CompareAndDelete deletes the entry for key if its value is equal to old.
	// The old value must be of a comparable type.
	// If there is no current value for key, CompareAndDelete returns false.
	CompareAndDelete(key K, old T) (deleted bool)
	// Range calls f sequentially for each key and value.
	// If f returns false, range stops the iteration.
	Range(f func(key K, value any) bool)
}

// Load is the helper function for [Store].
// It provides typed operation for the store with the any typed value.
// Load returns the stored value for a key from the store.
// If the key was not found, it returns zero value of type T and false.
// If type conversion failed, it returns zero value of type T and false.
// If nil was loaded, it always results in type assertion failure.
func Load[K comparable, T any](store Store[K, any], key K) (value T, ok bool) {
	v, ok := store.Load(key)
	if !ok {
		return value, false
	}
	value, ok = v.(T)
	return value, ok
}

// Swap is the helper function for [Store].
// It provides typed operation for the store with the any typed value.
// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
// If type conversion failed, it returns zero value of type T and false as ok.
// If nil was loaded, it always results in type assertion failure.
func Swap[K comparable, T any](store Store[K, any], key K, value any) (previous T, loaded, ok bool) {
	v, loaded := store.Swap(key, value)
	if !loaded {
		return previous, false, false
	}
	previous, ok = v.(T)
	return previous, true, ok
}

// LoadAndDelete is the helper function for [Store].
// It provides typed operation for the store with the any typed value.
// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
// If type conversion failed, it returns zero value of type T and false as ok.
// If nil was loaded, it always results in type assertion failure.
func LoadAndDelete[K comparable, T any](store Store[K, any], key K) (value T, loaded, ok bool) {
	v, loaded := store.LoadAndDelete(key)
	if !loaded {
		return value, false, false
	}
	value, ok = v.(T)
	return value, true, ok
}

// LoadOrStore is the helper function for [Store].
// It provides typed operation for the store with the any typed value.
// LoadOrStore returns the existing value for the key if present obtained from the store.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
// If type conversion failed, it returns zero value of type T and false as ok.
// If nil was loaded, it always results in type assertion failure.
func LoadOrStore[K comparable, T any](store Store[K, any], key K, value any) (actual T, loaded, ok bool) {
	v, loaded := store.LoadOrStore(key, value)
	if !loaded {
		return actual, false, false
	}
	actual, ok = v.(T)
	return actual, true, ok
}
