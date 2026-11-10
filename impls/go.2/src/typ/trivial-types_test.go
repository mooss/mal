package typ_test

import (
	"strings"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/typ"
)

func assertType[T core.Value](t *testing.T, v core.Value) {
	t.Helper()
	_, err := typ.As[T](v)
	if err.NotNil() {
		t.Fatalf("%v: %v", v, err)
	}
}

func assertString(t *testing.T, v core.Value, want string) {
	t.Helper()
	if got := v.String(); got != want {
		t.Fatalf("%v (%T).String() = %q, want %q", v, v, got, want)
	}
}

func assertRepr[T core.Value](t *testing.T, v core.Value, wantStr string) {
	t.Helper()
	assertType[T](t, v)
	assertString(t, v, wantStr)
}

func assertSelfEval(t *testing.T, env core.Env, v core.Value) {
	t.Helper()
	if got := v.Eval(env); !typ.Equal(got, v) {
		t.Fatalf("%v (%T).Eval() = %v, want %v", v, v, got, v)
	}
}

func assertTrivial[T core.Value](t *testing.T, env core.Env, v core.Value, wantStr string) {
	t.Helper()
	assertRepr[T](t, v, wantStr)
	assertSelfEval(t, env, v)
}

func assertEval(t *testing.T, env core.Env, v, want core.Value) {
	t.Helper()
	if got := v.Eval(env); !typ.Equal(got, want) {
		t.Fatalf("%v (%T).Eval() = %v, want %v", v, v, got, want)
	}
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%q does not contain %q", got, want)
	}
}

func assertNotNil(t *testing.T, f typ.Fatal, want bool) {
	t.Helper()
	if got := f.NotNil(); got != want {
		t.Fatalf("Fatal.NotNil() = %v, want %v", got, want)
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func mustEvalAs[T core.Value](t *testing.T, env core.Env, input string) T {
	t.Helper()
	got := tut.MustEval(env, input)
	v, ok := any(got).(T)
	if !ok {
		var zero T
		t.Fatalf("MustEval(%q) type = %T, want %T", input, got, zero)
	}
	return v
}

func TestBoolTrivial(t *testing.T) {
	var envr core.Env

	for _, tc := range []struct {
		val typ.Bool
		str string
	}{
		{typ.True, "true"},
		{typ.False, "false"},
	} {
		assertTrivial[typ.Bool](t, envr, tc.val, tc.str)
	}
}

func TestIntegerMustEval(t *testing.T) {
	var envr core.Env

	for _, tc := range []struct {
		input string
		want  typ.Integer
	}{
		{"42", typ.Integer(42)},
		{"0", typ.Integer(0)},
		{"-7", typ.Integer(-7)},
	} {
		got := mustEvalAs[typ.Integer](t, envr, tc.input)
		assertEqual(t, got, tc.want)
		assertTrivial[typ.Integer](t, envr, got, tc.input)
	}
}

func TestNilMustEval(t *testing.T) {
	var envr core.Env

	got := mustEvalAs[typ.Nil](t, envr, "nil")
	assertTrivial[typ.Nil](t, envr, got, "nil")
}

func TestStringMustEval(t *testing.T) {
	var envr core.Env

	for _, tc := range []struct {
		input string
		want  typ.String
	}{
		{`"hello"`, typ.String("hello")},
		{`""`, typ.String("")},
		{`"a b"`, typ.String("a b")},
	} {
		got := mustEvalAs[typ.String](t, envr, tc.input)
		assertEqual(t, got, tc.want)
		assertTrivial[typ.String](t, envr, got, tc.input)
	}

	// String prints quoted.
	s := typ.String("hi")
	assertTrivial[typ.String](t, envr, s, `"hi"`)
}

func TestQuoteMustEval(t *testing.T) {
	var envr core.Env

	assertEqual(t, mustEvalAs[typ.Integer](t, envr, "'42"), typ.Integer(42))
	_ = mustEvalAs[typ.Nil](t, envr, "'nil")
	assertEqual(t, mustEvalAs[typ.String](t, envr, `'"hi"`), typ.String("hi"))

	// Quote type itself: String, Type and Eval return raw value.
	q := typ.Quote{Value: typ.Integer(1)}
	assertRepr[typ.Quote](t, q, "'1")
	assertEval(t, envr, q, typ.Integer(1))
}

func TestSymbolViaQuoteMustEval(t *testing.T) {
	var envr core.Env

	// Quoting prevents symbol lookup, so this works with an empty env
	// and proves Quote returns the raw unevaluated form.
	got := mustEvalAs[typ.Symbol](t, envr, "'my-sym")
	assertEqual(t, string(got), "my-sym")
	assertRepr[typ.Symbol](t, got, "my-sym")

	// Direct Symbol trivial behaviour (no env lookup).
	direct := typ.Symbol("foo")
	assertRepr[typ.Symbol](t, direct, "foo")
}

func TestFatalTrivial(t *testing.T) {
	var envr core.Env

	f := typ.Fatalf("oops")
	assertRepr[typ.Fatal](t, f, `(fatal "oops")`)
	assertNotNil(t, f, true)
	var empty typ.Fatal
	assertNotNil(t, empty, false)
	assertString(t, f.Eval(envr), `(fatal "oops")`)

	qualified := f.Qualify("test")
	assertNotNil(t, qualified, true)
	assertContains(t, qualified.Error(), "test: oops")
	assertContains(t, qualified.String(), "test: oops")
}
