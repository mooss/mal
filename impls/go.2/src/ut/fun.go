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

func Reduce[T any](fun func(T, T) T, acc T, rest iter.Seq[T]) T {
	for el := range rest {
		acc = fun(acc, el)
	}
	return acc
}
