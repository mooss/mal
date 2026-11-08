package typ

import (
	"iter"

	"github.com/mooss/mal/src/core"
)

// As attempts to assert that val has the intended type T.
// When the type is not expected, a Fatal "expected ..." error is returned.
// When the unexpected type is already a Fatal error, it is returned as is.
// Intended to be used for type assertions and to preserve and forward the errors encountered during
// evaluation.
func As[T core.Value](val core.Value) (T, Fatal) {
	res, isT := val.(T)
	if !isT {
		err, isFatal := val.(Fatal)
		if !isFatal {
			err = expected(res.Type(), val)
		}
		return res, err
	}
	return res, Fatal{}
}

// AsSeq returns an iterator yielding all values as the type T.
//
// The iterator dies at the first unexpected type.
// It is stateful and resumable.
// When iteration is stopped, it will resume at the same place, not restart at the beginning.
//
// The caller must check the error getter to see if the iterator died.
// When the type is not expected, it will return a Fatal "expected ..." error.
// When the unexpected type is already a Fatal error, it is returned as is by the error getter.
//
// Intended to be used for type assertions and to preserve and forward the errors encountered during
// evaluation.
//
//revive:disable:cognitive-complexity
func AsSeq[T core.Value](values ...core.Value) (iter.Seq[T], func() Fatal) {
	var fatal Fatal
	i := 0
	return func(yield func(T) bool) {
		for i < len(values) {
			val := values[i]
			typed, isT := val.(T)
			if !isT {
				err, isFatal := val.(Fatal)
				if !isFatal {
					err = expected(typed.Type(), val)
				}
				fatal = err
				i = len(values)
				return
			}
			i++
			if !yield(typed) {
				return
			}
		}
	}, func() Fatal { return fatal }
}

// EvalArgs evaluates all its arguments in order, returning them as a slice.
// The first Fatal error encountered triggers a short-circuit and is returned as an error.
func EvalArgs(envr core.Env, args ...core.Value) ([]core.Value, Fatal) {
	res := make([]core.Value, len(args))
	for i, arg := range args {
		evald := arg.Eval(envr)
		if err, isFatal := evald.(Fatal); isFatal {
			return nil, err
		}
		res[i] = evald
	}
	return res, Fatal{}
}
