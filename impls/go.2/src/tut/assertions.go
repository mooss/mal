// Package tut provides test utilities.
package tut

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

// IsFatal reports whether v is a Fatal value.
func IsFatal(v core.Value) bool {
	if v == nil {
		return false
	}
	if _, ok := v.(typ.Fatal); ok {
		return true
	}
	if _, ok := v.(*typ.Fatal); ok {
		panic("got fatal pointer")
	}
	return false
}

// RequireAs asserts v has type T and returns it.
func RequireAs[T core.Value](t *testing.T, v core.Value) T {
	t.Helper()
	got, ok := v.(T)
	if !ok {
		var zero T
		t.Fatalf("type = %T, want %T", v, zero)
	}
	return got
}

// RequireFatal fails unless v is Fatal and returns it.
func RequireFatal(t *testing.T, v core.Value) typ.Fatal {
	t.Helper()
	if f, ok := v.(typ.Fatal); ok {
		return f
	}
	if _, ok := v.(*typ.Fatal); ok {
		panic("got fatal pointer")
	}

	//revive:disable:unreachable-code (>< I still have to return don't I?)
	t.Fatalf("got %v (%T), want Fatal", v, v)
	var zero typ.Fatal
	return zero
}

// RequireOK fails if v is Fatal and returns v.
func RequireOK(t *testing.T, v core.Value) core.Value {
	t.Helper()
	if IsFatal(v) {
		t.Fatalf("got Fatal %v, want success", v)
	}
	return v
}

// AssertEqual fails if got != want according to reflect.DeepEqual.
func AssertEqual[T any](t *testing.T, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// AssertWith fails when compare(got, want) returns false.
func AssertWith[T any](t *testing.T, got, want T, name string, compare func(a, b T) bool) {
	t.Helper()
	if !compare(got, want) {
		t.Fatalf("%s(%#v, %#v) failed", name, got, want)
	}
}

// AssertMalEqual fails unless MalEqual(got, want).
func AssertMalEqual(t *testing.T, got, want core.Value) {
	t.Helper()
	AssertWith(t, got, want, "MalEqual", MalEqual)
}

// AssertMalDifferent fails unless MalDifferent(got, want).
func AssertMalDifferent(t *testing.T, got, want core.Value) {
	t.Helper()
	AssertWith(t, got, want, "MalDifferent", MalDifferent)
}

// AssertAllMalEqual fails unless MalEqualSlices(got, want).
func AssertAllMalEqual(t *testing.T, got, want []core.Value) {
	t.Helper()
	AssertWith(t, got, want, "EqualSlices", MalEqualSlices)
}

// AssertNoErr fails if err != nil.
func AssertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertErr fails if err == nil.
func AssertErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// AssertContains fails if got lacks sub.
func AssertContains(t *testing.T, got, sub string) {
	t.Helper()
	if !strings.Contains(got, sub) {
		t.Fatalf("%q does not contain %q", got, sub)
	}
}

// AssertString fails if v.String() != want.
func AssertString(t *testing.T, v core.Value, want string) {
	t.Helper()
	if v == nil {
		t.Fatalf("got nil, want %q", want)
	}
	if got := v.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

// AssertGet fails unless env.Get(key) strictly equals want.
func AssertGet(t *testing.T, env core.Env, key string, want core.Value) {
	t.Helper()
	if got := env.Get(key); !MalEqual(got, want) {
		t.Fatalf("Get(%q) = %v, want %v", key, got, want)
	}
}

// AssertMissing fails if any key is bound.
func AssertMissing(t *testing.T, env core.Env, keys ...string) {
	t.Helper()
	for _, k := range keys {
		if got := env.Get(k); got != nil {
			t.Fatalf("Get(%q) = %v, want nil", k, got)
		}
	}
}

// AssertEval evaluates input and checks strict equality with want.
func AssertEval(t *testing.T, env core.Env, input string, want core.Value) {
	t.Helper()
	got := RequireOK(t, MustEval(env, input))
	if !MalEqual(got, want) {
		t.Fatalf("MustEval(%q) = %v, want %v", input, got, want)
	}
}

// AssertFatal checks that each input evaluates to Fatal.
func AssertFatal(t *testing.T, env core.Env, inputs ...string) {
	t.Helper()
	for _, in := range inputs {
		_ = RequireFatal(t, MustEval(env, in))
	}
}

// AssertEquivalent evaluates input and wantSrc and checks Mal equality.
// The environment is cloned when evaluating to protect against side effects.
func AssertEquivalent(t *testing.T, env core.Env, input, wantSrc string) {
	t.Helper()
	got := RequireOK(t, MustEval(core.CatEnv(env), input))
	want := RequireOK(t, MustEval(core.CatEnv(env), wantSrc))
	if !got.Equal(want) {
		t.Fatalf("MustEval(%q) = %v, want %v from %q", input, got, want, wantSrc)
	}
}

// AssertPrint evaluates input, checks it returns Nil and prints wantOutput.
func AssertPrint(t *testing.T, env core.Env, input, wantOutput string) {
	t.Helper()
	output := Capture(t, func() { AssertEval(t, env, input, Nil()) })
	AssertEqual(t, output, wantOutput)
}

// AssertFatalContains checks that input evaluates to Fatal containing sub.
func AssertFatalContains(t *testing.T, env core.Env, input, sub string) typ.Fatal {
	t.Helper()
	f := RequireFatal(t, MustEval(env, input))
	AssertContains(t, f.String(), sub)
	return f
}

// AssertNil fails unless v is nil or Nil.
func AssertNil(t *testing.T, v core.Value) {
	t.Helper()
	if v == nil {
		return
	}
	if _, ok := v.(typ.Nil); ok {
		return
	}
	t.Fatalf("got %v (%T), want nil", v, v)
}
