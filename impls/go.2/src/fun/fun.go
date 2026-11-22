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

// MapFlat applies fun to every key and value of source, returning a slice with the flattened
// results in key-value order.
func MapFlat[From comparable, To any](fun func(From) To, source map[From]From) []To {
	res := make([]To, len(source)*2)
	i := 0
	for k, v := range source {
		res[i] = fun(k)
		i++
		res[i] = fun(v)
		i++
	}
	return res
}

func Reduce[T any](fun func(T, T) T, acc T, rest iter.Seq[T]) T {
	for el := range rest {
		acc = fun(acc, el)
	}
	return acc
}
