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

// EvalFatal checks that each input evaluates to Fatal.
func EvalFatal(t *testing.T, env core.Env, inputs ...string) {
	t.Helper()
	for _, in := range inputs {
		_ = Fatal(t, tut.MustEval(env, in))
	}
}

// EvalFatalContains checks that input evaluates to Fatal containing sub.
func EvalFatalContains(t *testing.T, env core.Env, input, sub string) typ.Fatal {
	t.Helper()
	f := Fatal(t, tut.MustEval(env, input))
	StringContains(t, f.String(), sub)
	return f
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

// Print evaluates input, checks it returns Nil and prints wantOutput.
func Print(t *testing.T, env core.Env, input, wantOutput string) {
	t.Helper()
	output := tut.Capture(t, func() { Eval(t, env, input, typ.Nil{}) })
	Equal(t, output, wantOutput)
}
