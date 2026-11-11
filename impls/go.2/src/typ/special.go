package typ

import "github.com/mooss/mal/src/core"

// Special represents a special form that can alter the environment and evaluate its arguments
// according to its own rules.
type Special func(core.Env, ...core.Value) core.Value

func (Special) Equal(core.Value) bool        { return false }
func (spe Special) Eval(core.Env) core.Value { return spe }
func (Special) String() string               { return "<internal special form>" }
func (Special) Type() core.TypeName          { return core.SPECIAL }

func Sp0[R core.Value](fun func(core.Env) R) Special {
	return func(env core.Env, args ...core.Value) core.Value {
		if err := as0Arg(args); err.NotNil() {
			return err
		}
		return fun(env)
	}
}

func SpII[A, B, R core.Value](fun func(core.Env, A, B) R) Special {
	return func(env core.Env, args ...core.Value) core.Value {
		a, b, err := as2Args[A, B](args)
		if err.NotNil() {
			return err
		}

		return fun(env, a, b)
	}
}

func SpIII[A, B, C, R core.Value](fun func(core.Env, A, B, C) R) Special {
	return func(env core.Env, args ...core.Value) core.Value {
		a, b, c, err := as3Args[A, B, C](args)
		if err.NotNil() {
			return err
		}

		return fun(env, a, b, c)
	}
}
