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

/////////////////////
// Type assertions //

func as2Args[A, B core.Value](args []core.Value) (A, B, Fatal) {
	var a A
	var b B

	if len(args) != 2 {
		return a, b, Fatalf("expected 2 args, got %d", len(args))
	}

	a, err := As[A](args[0])
	if err.NotNil() {
		return a, b, err
	}
	b, err = As[B](args[1])
	if err.NotNil() {
		return a, b, err
	}

	return a, b, Fatal{}
}

// validKey returns true when the given value can be used for a map key.
func validKey(value core.Value) bool {
	switch value.(type) {
	case Integer, Nil, String, Symbol:
		return true
	}
	return false
}
