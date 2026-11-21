package typ

import (
	"fmt"
	"iter"
	"slices"
	"strconv"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/damn"
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

// AsListOrVec converts a list, vector, or vector form to a slice of values, or returns a fatal
// error.
func AsListOrVec(elements core.Value) ([]core.Value, Fatal) {
	var bindlist []core.Value
	switch typed := elements.(type) {
	case List:
		bindlist = typed
	case Vector:
		bindlist = typed
	case VectorForm:
		bindlist = typed
	default:
		msg := fmt.Sprintf("%s, %s or %s", core.LIST, core.VECTOR, core.VECTOR_FORM)
		return nil, Expected(msg, elements)
	}
	return bindlist, Fatal{}
}

// AsFunOrMac converts a function or macro to its underlying Go function or returns a fatal error.
func AsFunOrMac(fun core.Value) (func(...core.Value) core.Value, Fatal) {
	switch typed := fun.(type) {
	case Function:
		return typed, Fatal{}
	case Macro:
		return typed, Fatal{}
	default:
		return nil, Expected("function or macro", fun)
	}
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
		evald := Resolve(envr, arg)
		if err, isFatal := evald.(Fatal); isFatal {
			return nil, err
		}
		res[i] = evald
	}
	return res, Fatal{}
}

// ExpectAtLeastNArgs returns a nil Fatal if got contains at least expected values, otherwise a
// Fatal "expected ... args, got ..." error.
func ExpectAtLeastNArgs(expected int, got []core.Value) Fatal {
	if len(got) >= expected {
		return Fatal{}
	}
	return FmtExpectedNArgs(fmt.Sprintf("at least %d", expected), len(got))
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
	if err, isFatal := got.(Fatal); isFatal {
		return err
	}
	return ThrowErr(damn.TypeMismatch(expected, got))
}

// FmtExpectedNArgs format a Fatal "expected ... args, got ..." error.
func FmtExpectedNArgs(expected string, got int) Fatal {
	return ThrowErr(&damn.ArityError{Want: expected, Got: got})
}

// IsTrue returns whether the value is considered truthful, discarding any Fatal error.
// Intended for cases where the error is handled elsewhere or irrelevant.
func IsTrue(arg core.Value) bool {
	res, _ := Truthful(arg)
	return res
}

// TrivialEq implements trivial equality comparison, when right has the same type as left, the two
// are compared with Go's == operator.
// Intended to be used as an equality predicate for comparable values of a known concrete type.
func TrivialEq[T comparable](left T, right any) bool {
	rt, ok := right.(T)
	return ok && left == rt
}

// Truthful returns whether the value is considered truthful in a boolean context.
// A nil value is falsy, a Fatal error is forwarded as is, a Bool is returned as is and any other
// value is considered truthful.
func Truthful(arg core.Value) (bool, Fatal) {
	if arg == nil {
		return false, Fatal{}
	}

	realised := false
	switch typed := arg.(type) {
	case Fatal:
		return false, typed
	case Nil:
	case Bool:
		realised = bool(typed)
	default:
		realised = true
	}
	return realised, Fatal{}
}
