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

func as[T core.Value](val core.Value) (T, *Fatal) {
	res, isT := val.(T)
	if !isT {
		err, isErr := val.(Fatal)
		if !isErr {
			err = Fatalf("expected %T but got %T (%v)", res, val, val)
		}
		return res, &err
	}
	return res, nil
}

func validKey(value core.Value) bool {
	switch value.(type) {
	case Integer, Nil, String, Symbol:
		return true
	}
	return false
}
