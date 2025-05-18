package zmath

import (
	"cmp"

	"slices"
)

// NewSet returns a new sorted set from the given slice x.
// It modifies the given slice of x.
// Callers should not modify the slice x after called NewSet.
func NewSet[S ~[]E, E cmp.Ordered](x S) Set[E] {
	slices.Sort(x)
	return Set[E](x)
}

// Set is the set for mathematical set operations.
// The slice []E must be sorted for correct calculation.
type Set[E cmp.Ordered] []E

// Has returns if the set s has at least one elem.
func (s *Set[E]) Has(elem E) bool {
	for _, e := range *s {
		switch {
		case e < elem:
			continue
		case e == elem:
			return true
		case e > elem:
			return false
		}
	}
	return false
}

// Add adds elem to the set.
func (s *Set[E]) Add(elem E) {
	s.addAfter(0, elem)
}

// AddElems adds multiple elements to the set.
func (s *Set[E]) AddElems(elems ...E) {
	slices.Sort(elems)
	index := 0
	for _, e := range elems {
		index = s.addAfter(index, e)
	}
}

func (s *Set[E]) addAfter(i int, elem E) int {
	if i >= len(*s) {
		*s = append((*s), elem)
		return len(*s)
	}
	for j := i; j < len(*s); j++ {
		v := (*s)[j]
		switch {
		case v < elem:
			continue
		case v > elem:
			*s = append((*s)[:j], append([]E{elem}, (*s)[j:]...)...)
			return j
		}
	}
	*s = append((*s), elem)
	return len(*s)
}

// Remove removes the given elem from the set.
func (s *Set[E]) Remove(elem E) {
	s.removeAfter(0, elem)
}

// RemoveElems removes multiple elements from the set.
func (s *Set[E]) RemoveElems(elems ...E) {
	slices.Sort(elems)
	index := 0
	for _, e := range elems {
		index = s.removeAfter(index, e)
	}
}

func (s *Set[E]) removeAfter(i int, elem E) int {
	if i >= len(*s) {
		return len(*s)
	}
	for j := i; j < len(*s); j++ {
		v := (*s)[j]
		switch {
		case v < elem:
			continue
		case v == elem:
			*s = append((*s)[:j], (*s)[j+1:]...)
			return j
		case v > elem:
			return j
		}
	}
	return len(*s)
}

// Equal returns the result of A=B.
// In other words, A⊆B and B⊆A.
// For empty set ∅ and non-empty set A,B, result will be
//
//   - ∅=∅ : true
//   - ∅=B : false
//   - A=∅ : false
func Equal[E cmp.Ordered](a, b Set[E]) bool {
	n, m := len(a), len(b)
	if n == 0 {
		return m == 0 // ∅=∅, ∅=B
	}
	if m == 0 {
		return n == 0 // ∅=∅, A=∅
	}
	aIter := uniqueIterator(a)
	bIter := uniqueIterator(b)
	for {
		an, af := aIter()
		bn, bf := bIter()
		if af != bf {
			return false
		}
		if !af && !bf {
			break
		}
		if an != bn {
			return false
		}
	}
	return true
}

// Superset returns the result of A⊇B.
// [Superset] is the inverse operation of the [Subset].
// For empty set ∅ and non-empty set A,B, result will be
//
//   - ∅⊇∅ : true
//   - ∅⊇B : false
//   - A⊇∅ : true
func Superset[E cmp.Ordered](a, b Set[E]) bool {
	return Subset(b, a)
}

// ProperSuperset returns the result of A⊃B.
// [ProperSuperset] is the inverse operation of the [ProperSubset].
// For empty set ∅ and non-empty set A,B, result will be
//
//   - ∅⊃∅ : false
//   - ∅⊃B : false
//   - A⊃∅ : true
func ProperSuperset[E cmp.Ordered](a, b Set[E]) bool {
	return ProperSubset(b, a)
}

// Subset returns the result of A⊆B.
// [Subset] is the inverse operation of the [Superset].
// For empty set ∅ and non-empty set A,B, result will be
//
//   - ∅⊆∅ : true
//   - ∅⊆B : true
//   - A⊆∅ : false
func Subset[E cmp.Ordered](a, b Set[E]) bool {
	n, m := len(a), len(b)
	if n == 0 {
		return true // ∅⊆∅, ∅⊆B
	}
	if m == 0 {
		return n == 0 // ∅⊆∅, A⊆∅
	}
	aIter := uniqueIterator(a)
	bIter := uniqueIterator(b)
	an, _ := aIter()
	bn, _ := bIter()
	var af, bf bool
	for {
		switch {
		case an == bn:
			bn, bf = bIter()
			an, af = aIter()
			if !af {
				return true
			}
			if !bf {
				return false
			}
		case an < bn:
			return false
		case an > bn:
			bn, bf = bIter()
			if !bf {
				return false
			}
		}
	}
}

// ProperSubset returns the result of A⊂B.
// [ProperSubset] is the inverse operation of the [ProperSuperset].
// For empty set ∅ and non-empty set A,B, result will be
//
//   - ∅⊂∅ : false
//   - ∅⊂B : true
//   - A⊂∅ : false
func ProperSubset[E cmp.Ordered](a, b Set[E]) bool {
	n, m := len(a), len(b)
	if n == 0 {
		return m != 0 // ∅⊂∅, ∅⊂B
	}
	if m == 0 {
		return false // ∅⊂∅ , A⊂∅
	}
	aIter := uniqueIterator(a)
	bIter := uniqueIterator(b)
	an, _ := aIter()
	bn, _ := bIter()
	var af, bf, notFound bool
	for {
		switch {
		case an == bn:
			an, af = aIter()
			bn, bf = bIter()
			if !af {
				return bf || notFound
			}
			if !bf {
				return false
			}
		case an < bn:
			return false
		case an > bn:
			bn, bf = bIter()
			if !bf {
				return false
			}
			notFound = true
		}
	}
}

// Union returns the result of A∪B.
// For empty set ∅ and non-empty set A,B, result will be
//
//   - ∅∪∅ : ∅
//   - ∅∪B : B
//   - A∪∅ : A
func Union[E cmp.Ordered](a, b Set[E]) Set[E] {
	n, m := len(a), len(b)
	if n == 0 && m == 0 {
		return nil // ∅∪∅
	}
	if n == 0 {
		return b // ∅∪B
	}
	if m == 0 {
		return a // A∪∅
	}
	set := Set[E]{}
	aIter := uniqueIterator(a)
	bIter := uniqueIterator(b)
	an, af := aIter()
	bn, bf := bIter()
	for {
		switch {
		case an == bn:
			set = append(set, an)
			an, af = aIter()
			bn, bf = bIter()
		case an < bn:
			set = append(set, an)
			an, af = aIter()
		case an > bn:
			set = append(set, bn)
			bn, bf = bIter()
		}
		if !af {
			for bf {
				set = append(set, bn)
				bn, bf = bIter()
			}
			return set
		}
		if !bf {
			for af {
				set = append(set, an)
				an, af = aIter()
			}
			return set
		}
	}
}

// Intersection returns the result of A∩B.
// For empty set ∅ and non-empty set A,B, result will be
//
//   - ∅∩∅ : ∅
//   - ∅∩B : ∅
//   - A∩∅ : ∅
func Intersection[E cmp.Ordered](a, b Set[E]) Set[E] {
	n, m := len(a), len(b)
	if n == 0 || m == 0 {
		return nil // ∅∩∅, ∅∩B, A∩∅
	}
	set := Set[E]{}
	aIter := uniqueIterator(a)
	bIter := uniqueIterator(b)
	an, af := aIter()
	bn, bf := bIter()
	for {
		switch {
		case an == bn:
			set = append(set, an)
			an, af = aIter()
			bn, bf = bIter()
		case an < bn:
			an, af = aIter()
		case an > bn:
			bn, bf = bIter()
		}
		if !af || !bf {
			return set
		}
	}
}

// Difference returns the result of A−B.
// For empty set ∅ and non-empty set A,B, result will be
//
//   - ∅−∅ : ∅
//   - ∅−B : ∅
//   - A−∅ : A
func Difference[E cmp.Ordered](a, b Set[E]) Set[E] {
	n, m := len(a), len(b)
	if n == 0 {
		return nil // ∅−∅, ∅−B
	}
	if m == 0 {
		return a // A−∅
	}
	set := Set[E]{}
	aIter := uniqueIterator(a)
	bIter := uniqueIterator(b)
	an, af := aIter()
	bn, bf := bIter()
	for {
		switch {
		case an == bn:
			an, af = aIter()
			bn, bf = bIter()
		case an < bn:
			set = append(set, an)
			an, af = aIter()
		case an > bn:
			bn, bf = bIter()
		}
		if !af {
			return set
		}
		if !bf {
			for af {
				set = append(set, an)
				an, af = aIter()
			}
			return set
		}
	}
}

func uniqueIterator[E cmp.Ordered](set Set[E]) func() (E, bool) {
	n := len(set)
	i := 0
	var zero E
	return func() (E, bool) {
		if i >= n {
			return zero, false
		}
		e := set[i]
		for i++; i < n && e == set[i]; i++ {
		}
		return e, true
	}
}
