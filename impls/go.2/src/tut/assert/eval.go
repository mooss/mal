package assert

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/eq"
	"github.com/mooss/mal/src/typ"
)

// Eval evaluates input and checks strict equality with want.
func Eval(t *testing.T, env core.Env, input string, want core.Value) {
	t.Helper()
	got := Ok(t, tut.MustEval(env, input))
	if !eq.Mal(got, want) {
		t.Fatalf("MustEval(%q) = %v, want %v", input, got, want)
	}
}
func (tr T) Eval(env core.Env, input string, want core.Value) {
	tr.T.Helper()
	Eval(tr.T, env, input, want)
}

// EvalFatal checks that each input evaluates to Fatal.
func EvalFatal(t *testing.T, env core.Env, inputs ...string) {
	t.Helper()
	for _, in := range inputs {
		_ = Fatal(t, tut.MustEval(env, in))
	}
}
func (tr T) EvalFatal(env core.Env, inputs ...string) {
	tr.T.Helper()
	EvalFatal(tr.T, env, inputs...)
}

// EvalFatalContains checks that input evaluates to Fatal containing sub.
func EvalFatalContains(t *testing.T, env core.Env, input, sub string) typ.Fatal {
	t.Helper()
	f := Fatal(t, tut.MustEval(env, input))
	StringContains(t, f.String(), sub)
	return f
}
func (tr T) EvalFatalContains(env core.Env, input, sub string) typ.Fatal {
	tr.T.Helper()
	return EvalFatalContains(tr.T, env, input, sub)
}

// EvalSame evaluates input and wantSrc and checks Mal equality.
// The environment is cloned when evaluating to protect against side effects.
func EvalSame(t *testing.T, env core.Env, input, wantSrc string) {
	t.Helper()
	got := Ok(t, tut.MustEval(core.CatEnv(env), input))
	want := Ok(t, tut.MustEval(core.CatEnv(env), wantSrc))
	if !got.Equal(want) {
		t.Fatalf("MustEval(%q) = %v, want %v from %q", input, got, want, wantSrc)
	}
}
func (tr T) EvalSame(env core.Env, input, wantSrc string) {
	tr.T.Helper()
	EvalSame(tr.T, env, input, wantSrc)
}

// Print evaluates input, checks it returns Nil and prints wantOutput.
func Print(t *testing.T, env core.Env, input, wantOutput string) {
	t.Helper()
	output := tut.Capture(t, func() { Eval(t, env, input, typ.Nil{}) })
	Equal(t, output, wantOutput)
}
func (tr T) Print(env core.Env, input, wantOutput string) {
	tr.T.Helper()
	Print(tr.T, env, input, wantOutput)
}

// EvalValue checks v.Eval(env) equals want.
func EvalValue(t *testing.T, env core.Env, v, want core.Value) {
	t.Helper()
	if got := v.Eval(env); !got.Equal(want) {
		t.Fatalf("%v (%T).Eval() = %v, want %v", v, v, got, want)
	}
}
func (tr T) EvalValue(env core.Env, v, want core.Value) {
	tr.T.Helper()
	EvalValue(tr.T, env, v, want)
}

// SelfEval checks v.Eval(env) equals itself.
func SelfEval(t *testing.T, env core.Env, v core.Value) {
	t.Helper()
	EvalValue(t, env, v, v)
}
func (tr T) SelfEval(env core.Env, v core.Value) {
	tr.T.Helper()
	SelfEval(tr.T, env, v)
}

// Repr checks v has type T and its String equals want.
func Repr[T core.Value](t *testing.T, v core.Value, want string) {
	t.Helper()
	if v == nil {
		t.Fatalf("got nil, want %q", want)
	}
	Is[T](t, v)
	StringRepr(t, v, want)
}
func (tr T) Repr[T core.Value](v core.Value, want string) {
	tr.T.Helper()
	Repr[T](tr.T, v, want)
}

// Trivial checks Repr and SelfEval.
func Trivial[T core.Value](t *testing.T, env core.Env, v core.Value, want string) {
	t.Helper()
	Repr[T](t, v, want)
	SelfEval(t, env, v)
}
func (tr T) Trivial[T core.Value](env core.Env, v core.Value, want string) {
	tr.T.Helper()
	Trivial[T](tr.T, env, v, want)
}
