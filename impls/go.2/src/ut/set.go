package ut

import "iter"

// Set is an unordered set of unique elements of type T.
type Set[T comparable] map[T]struct{}

// NewSetI returns a Set containing the given elements.
func NewSetI[T comparable](elements iter.Seq[T]) Set[T] {
	res := Set[T]{}
	for el := range elements {
		res[el] = struct{}{}
	}
	return res
}

// Contains reports whether el is in the set.
func (set Set[T]) Contains(el T) bool {
	_, ok := set[el]
	return ok
}
