package lazy

import "iter"

func Map[From, To any](
	fun func(From) To,
	source iter.Seq[From],
) iter.Seq[To] {
	return func(yield func(To) bool) {
		for el := range source {
			if !yield(fun(el)) {
				return
			}
		}
	}
}

func Map2[From1, From2, To any](
	fun func(From1, From2) To,
	source iter.Seq2[From1, From2],
) iter.Seq[To] {
	return func(yield func(To) bool) {
		for left, right := range source {
			if !yield(fun(left, right)) {
				return
			}
		}
	}
}

// CollectHint collects a sequence using a size hint to preallocate the results.
func CollectHint[T any](source iter.Seq[T], hint int) []T {
	res := make([]T, 0, hint)
	for el := range source {
		res = append(res, el)
	}
	return res
}
