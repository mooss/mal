// Package env implements environments.
package env

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

type _sym = typ.Symbol
type _lst = typ.List
type _fun = typ.Function
type _int = typ.Integer
type _vec = typ.Vector
type _err = typ.Error

type Env map[string]core.Type

var Initial Env = map[string]core.Type{
	"+": _fun(func(l _int, r _int) _int { return l + r }),
	"-": _fun(func(l _int, r _int) _int { return l - r }),
	"*": _fun(func(l _int, r _int) _int { return l * r }),
	"/": _fun(func(l _int, r _int) _int { return l / r }), // Will panic on r == 0.
}

//revive:disable:cognitive-complexity,cyclomatic
func Eval(ast core.Type, env Env) core.Type {
	switch typed := ast.(type) {
	case _sym:
		value, present := env[string(typed)]
		if !present {
			return typ.Errorf("symbol %#q not found", typed)
		}
		return value

	case _lst:
		if len(typed) == 0 {
			return typed
		}
		if len(typed) != 3 {
			msg := "only list of 3 symbols can be evaluated right now (got %d)"
			return typ.Errorf(msg, len(typed))
		}

		// Apply phase.
		fun, cerr := as[_fun](Eval(typed[0], env))
		if cerr != nil {
			return *cerr
		}

		arg1, cerr := as[_int](Eval(typed[1], env))
		if cerr != nil {
			return *cerr
		}

		arg2, cerr := as[_int](Eval(typed[2], env))
		if cerr != nil {
			return *cerr
		}

		return fun(arg1, arg2)

	case _vec:
		res := make(_vec, len(typed))
		for i, el := range typed {
			value := Eval(el, env)
			if ut.Is[typ.Error](value) {
				return value
			}
			res[i] = value
		}
		return res

	default:
		return ast
	}
}

func as[T core.Type](val core.Type) (T, *_err) {
	res, isT := val.(T)
	if !isT {
		err, isErr := val.(_err)
		if !isErr {
			err = typ.Errorf("expected %T but got %T (%v)", res, val, val)
		}
		return res, &err
	}
	return res, nil
}
