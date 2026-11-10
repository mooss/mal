package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

var (
	arityGt1Err = typ.Fatalf("expected at least 1 argument")
	arityGt2Err = typ.Fatalf("expected at least 2 argument")
)

func reduce[T core.Value](fun func(T, T) T) typ.Function {
	return func(args ...core.Value) core.Value {
		if len(args) == 0 {
			return arityGt1Err
		}

		typed, gerr := typ.AsSeq[T](args...)
		acc, _ := ut.First(typed)
		acc = ut.Reduce(fun, acc, typed)
		if err := gerr(); err.NotNil() {
			return err
		}

		return acc
	}
}

type numbers interface{ ~int64 | ~float64 }

func plus[T numbers](l, r T) T  { return l + r }
func minus[T numbers](l, r T) T { return l - r }
func times[T numbers](l, r T) T { return l * r }
func div[T numbers](l, r T) T   { return l / r } // Will panic on r == 0.

func minusInt(args ...core.Value) core.Value {
	switch len(args) {
	case 1:
		val, err := typ.As[_int](args[0])
		if err.NotNil() {
			return err
		}
		return -val
	default:
		return reduce(minus[_int])(args...)
	}
}

func divInt(args ...core.Value) core.Value {
	if len(args) < 2 {
		return arityGt2Err
	}

	return reduce(div[_int])(args...)
}

var Math core.Env = core.Bindings{
	"+":     reduce(plus[_int]),
	"-":     typ.Function(minusInt),
	"*":     reduce(times[_int]),
	"/":     typ.Function(divInt),
	"true":  typ.Bool(true),
	"false": typ.Bool(false),
}.AsEnv()
