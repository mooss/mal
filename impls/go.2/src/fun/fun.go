package fun

import (
	"iter"
)

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
