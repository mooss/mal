package typ

import (
	"maps"
	"slices"

	"github.com/mooss/mal/src/core"
)

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

/////////////////////
// Type assertions //

func as1Arg[A core.Value](args []core.Value) (a A, err Fatal) {
	if err = ExpectNArgs(1, args); err.NotNil() {
		return
	}

	a, err = As[A](args[0])
	if err.NotNil() {
		return
	}

	return a, Fatal{}
}

func as2Args[A, B core.Value](args []core.Value) (a A, b B, err Fatal) {
	if err = ExpectNArgs(2, args); err.NotNil() {
		return
	}

	a, err = As[A](args[0])
	if err.NotNil() {
		return
	}
	b, err = As[B](args[1])
	if err.NotNil() {
		return
	}

	return a, b, Fatal{}
}

func as3Args[A, B, C core.Value](args []core.Value) (a A, b B, c C, err Fatal) {
	if err = ExpectNArgs(3, args); err.NotNil() {
		return
	}

	a, err = As[A](args[0])
	if err.NotNil() {
		return
	}
	b, err = As[B](args[1])
	if err.NotNil() {
		return
	}
	c, err = As[C](args[2])
	if err.NotNil() {
		return
	}

	return a, b, c, Fatal{}
}

// validKey returns true when the given value can be used for a map key.
func validKey(value core.Value) bool {
	switch value.(type) {
	case Integer, Nil, String, Symbol:
		return true
	}
	return false
}

//////////////////////
// Equality helpers //

func assocEq(a, b Assoc) bool {
	return a.Key.Equal(b.Key) && a.Value.Equal(b.Value)
}

func mapEq(left Map, right core.Value) bool {
	r, ok := right.(Map)
	return ok && maps.EqualFunc(left, r, core.Value.Equal)
}

func sliceEq[E any, S ~[]E](left S, right core.Value, eq func(E, E) bool) bool {
	r, ok := right.(S)
	return ok && slices.EqualFunc(left, r, eq)
}
