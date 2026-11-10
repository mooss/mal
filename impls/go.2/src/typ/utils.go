package typ

import (
	"iter"
	"slices"
	"strconv"

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
		return res, Expected(res.Type(), val)
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
				fatal = Expected(typed.Type(), val)
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

// AsSlice asserts all the value as T.
//
// Intended to be used for type assertions and to preserve and forward the errors encountered during
// evaluation.
func AsSlice[T core.Value](values ...core.Value) ([]T, Fatal) {
	seq, geterr := AsSeq[T](values...)
	res := slices.Collect(seq)
	if err := geterr(); err.NotNil() {
		return nil, err
	}
	return res, Fatal{}
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

// ExpectNArgs returns a nil Fatal if got contains exactly expected values, otherwise a Fatal
// "expected ... args, got ..." error.
func ExpectNArgs(expected int, got []core.Value) Fatal {
	if expected == len(got) {
		return Fatal{}
	}
	return FmtExpectedNArgs(strconv.Itoa(expected), len(got))
}

// Expected format a Fatal error explaining that an unexpected type was given.
// When the unexpected type is already a Fatal error, it is returned as-is.
func Expected[T ~string](expected T, got core.Value) Fatal {
	err, isFatal := got.(Fatal)
	if isFatal {
		return err
	}
	return Fatalf("expected %s but got %s (%v)", expected, got.Type(), got)
}

// FmtExpectedNArgs format a Fatal "expected ... args, got ..." error.
func FmtExpectedNArgs(expected string, got int) Fatal {
	return Fatalf("expected %s args, got %d", expected, got)
}

func TrivialEq[T comparable](left T, right any) bool {
	rt, ok := right.(T)
	return ok && left == rt
}
