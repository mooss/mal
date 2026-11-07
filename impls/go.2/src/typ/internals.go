package typ

import "github.com/mooss/mal/src/core"

func _map[From, To any](fun func(From) To, source []From) []To {
	res := make([]To, len(source))
	for i, el := range source {
		res[i] = fun(el)
	}
	return res
}

func _mapflat[From comparable, To any](fun func(From) To, source map[From]From) []To {
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

// as attempts to assert that val has the intended type T.
// When the type is not expected, a Fatal "expected ..." error is returned.
// When the unexpected type is already a Fatal error, it is returned as is.
// Intended to be used for type assertions and to preserve and forward the errors encountered during
// evaluation.
func as[T core.Value](val core.Value) (T, Fatal) {
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

// expected format a Fatal error explaining that a type other than the expected one was found.
func expected[T ~string](expected T, got core.Value) Fatal {
	return Fatalf("expected %s but got %s (%v)", expected, got.Type(), got)
}

// validKey returns true when the given value can be used for a map key.
func validKey(value core.Value) bool {
	switch value.(type) {
	case Integer, Nil, String, Symbol:
		return true
	}
	return false
}
