package ut

import (
	"iter"
)

func First[T any](seq iter.Seq[T]) (T, bool) {
	for el := range seq {
		return el, true
	}
	var zero T
	return zero, false
}
