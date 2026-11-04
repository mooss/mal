// Package env implements environments.
package env

import (
	"github.com/mooss/mal/src/core"
)

type _sym = core.Symbol
type _lst = core.List
type _fun = core.Function
type _int = core.Integer
type _err = core.Error

type Env map[string]core.Function

var Initial Env = map[string]core.Function{
	"+": func(l _int, r _int) _int { return l + r },
	"-": func(l _int, r _int) _int { return l - r },
	"*": func(l _int, r _int) _int { return l * r },
	"/": func(l _int, r _int) _int { return l / r }, // Will panic on r == 0.
}

//revive:disable:cognitive-complexity
func Eval(ast core.Type, e Env) core.Type {
	switch typed := ast.(type) {
	case _sym:
		value, present := e[string(typed)]
		if !present {
			return core.Errorf("symbol %#q not found", typed)
		}
		return value
	case _lst:
		if len(typed) == 0 {
			return typed
		}
		if len(typed) != 3 {
			msg := "only list of 3 symbols can be evaluated right now (got %d)"
			return core.Errorf(msg, len(typed))
		}

		// Apply phase.
		fun, cerr := as[_fun](Eval(typed[0], e))
		if cerr != nil {
			return *cerr
		}

		arg1, cerr := as[_int](Eval(typed[1], e))
		if cerr != nil {
			return *cerr
		}

		arg2, cerr := as[_int](Eval(typed[2], e))
		if cerr != nil {
			return *cerr
		}

		return fun(arg1, arg2)
	default:
		return ast
	}
}

func as[T core.Type](val core.Type) (T, *_err) {
	res, isT := val.(T)
	if !isT {
		err, isErr := val.(_err)
		if !isErr {
			err = core.Errorf("expected %T but got %T (%v)", res, val, val)
		}
		return res, &err
	}
	return res, nil
}
