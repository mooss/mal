package typ

import "github.com/mooss/mal/src/core"

// TailCall is a deferred evaluation for the purpose of tail call optimization.
// It lets the language perform tail call optimization while keeping eval methods where each type is
// implemented.
type TailCall struct {
	Env  core.Env
	Form core.Value
}

func Tail(env core.Env, form core.Value) TailCall {
	return TailCall{env, form}
}

// Satisfy core.Value, but should never leak to the user.
func (TailCall) Equal(core.Value) bool    { return false }
func (TailCall) Eval(core.Env) core.Value { return Fatalf("attempt to evaluate a TailCall") }
func (TailCall) String() string           { return "<tailcall>" }
func (TailCall) Type() core.TypeName      { return core.TAILCALL }

// Resolve properly evaluates form in env, resolving all the tail calls.
// Useful resource: https://eli.thegreenplace.net/2017/on-recursion-continuations-and-trampolines/.
func Resolve(env core.Env, form core.Value) core.Value {
	for {
		val := form.Eval(env)
		tc, ok := val.(TailCall)
		if !ok {
			return val
		}
		form, env = tc.Form, tc.Env
	}
}
