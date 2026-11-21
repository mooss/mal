package fun

import (
	"iter"
)

func Filter[T any](pred func(T) bool, source iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for el := range source {
			if pred(el) && !yield(el) {
				return
			}
		}
	}
}

func Map[From, To any](fun func(From) To, source []From) []To {
	res := make([]To, len(source))
	for i, el := range source {
		res[i] = fun(el)
	}
	return res
}

func Reduce[T any](fun func(T, T) T, acc T, rest iter.Seq[T]) T {
	for el := range rest {
		acc = fun(acc, el)
	}
	return acc
}
