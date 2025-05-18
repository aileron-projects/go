package zmaps

// Keys returns all keys of the given map.
// Unlike [maps.Keys], it does not return iter.Seq[K] but return []K.
// Keys creates a new slice of []K which means it does not suite for
// large map input.
func Keys[Map ~map[K]V, K comparable, V any](m Map) []K {
	n := len(m)
	if n == 0 {
		return nil
	}
	keys := make([]K, 0, n)
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// Values returns all values of the given map.
// Unlike [maps.Values], it does not return iter.Seq[V] but return []V.
// Values creates a new slice of []V which means it does not suite for
// large map input.
func Values[Map ~map[K]V, K comparable, V any](m Map) []V {
	n := len(m)
	if n == 0 {
		return nil
	}
	values := make([]V, 0, n)
	for _, value := range m {
		values = append(values, value)
	}
	return values
}
