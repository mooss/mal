package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/fun"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

var Math core.Env = core.Bindings{
	"+":     reduce(plus[_int]),
	"-":     typ.Function(minusInt),
	"*":     reduce(times[_int]),
	"/":     typ.Function(divInt),
	"true":  typ.True,
	"false": typ.False,
	"<":     typ.II(lt[_int]),
	"<=":    typ.II(leq[_int]),
	">":     typ.II(gt[_int]),
	">=":    typ.II(geq[_int]),
}.AsEnv()

func reduce[T core.Value](fn func(T, T) T) typ.Function {
	return func(args ...core.Value) core.Value {
		if len(args) == 0 {
			return arityGt1Err
		}

		typed, gerr := typ.AsSeq[T](args...)
		acc, _ := ut.First(typed)
		acc = fun.Reduce(fn, acc, typed)
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

func lt[T numbers](l, t T) typ.Bool  { return l < t }
func leq[T numbers](l, t T) typ.Bool { return l <= t }
func gt[T numbers](l, t T) typ.Bool  { return l > t }
func geq[T numbers](l, t T) typ.Bool { return l >= t }
