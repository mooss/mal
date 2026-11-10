package typ

import "github.com/mooss/mal/src/core"

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

func as2Args[A, B core.Value](args []core.Value) (a A, b B, err Fatal) {
	if err = ExpectNArgs(2, args); err.NotNil() {
		return a, b, err
	}

	a, err = As[A](args[0])
	if err.NotNil() {
		return a, b, err
	}
	b, err = As[B](args[1])
	if err.NotNil() {
		return a, b, err
	}

	return a, b, Fatal{}
}

func as3Args[A, B, C core.Value](args []core.Value) (a A, b B, c C, err Fatal) {
	if err = ExpectNArgs(3, args); err.NotNil() {
		return a, b, c, err
	}

	a, err = As[A](args[0])
	if err.NotNil() {
		return a, b, c, err
	}
	b, err = As[B](args[1])
	if err.NotNil() {
		return a, b, c, err
	}
	c, err = As[C](args[2])
	if err.NotNil() {
		return a, b, c, err
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
