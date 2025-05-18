package zcmp

// True returns the yes when the b is true.
// Otherwise returns the value of no.
func True[T comparable](b bool, yes, no T) T {
	if b {
		return yes
	}
	return no
}

// False returns the yes when the b is false.
// Otherwise returns the value of no.
func False[T comparable](b bool, yes, no T) T {
	if !b {
		return yes
	}
	return no
}

// OrSlice returns the first of its arguments that is not empty.
// If no argument has element, it returns the nil slice.
func OrSlice[S []V, V any](vals ...S) S {
	for _, v := range vals {
		if len(v) > 0 {
			return v
		}
	}
	return nil
}

// OrMap returns the first of its arguments that is not empty.
// If no argument has element, it returns the nil map.
func OrMap[M map[K]V, K comparable, V any](vals ...M) M {
	for _, v := range vals {
		if len(v) > 0 {
			return v
		}
	}
	return nil
}
