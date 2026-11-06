package typ

import "github.com/mooss/mal/src/core"

type _err = Error
type _fun = Function
type _int = Integer
type _vec = Vector

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

func as[T core.Type](val core.Type) (T, *_err) {
	res, isT := val.(T)
	if !isT {
		err, isErr := val.(_err)
		if !isErr {
			err = Errorf("expected %T but got %T (%v)", res, val, val)
		}
		return res, &err
	}
	return res, nil
}
