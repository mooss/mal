package typ_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/tut/assert"
	"github.com/mooss/mal/src/typ"
)

func TestBoolTrivial(t *testing.T) {
	tr := assert.New(t)
	var envr core.Env

	for _, tc := range []struct {
		val typ.Bool
		str string
	}{
		{typ.True, "true"},
		{typ.False, "false"},
	} {
		tr.Trivial[typ.Bool](envr, tc.val, tc.str)
		tr.SelfEval(envr, tc.val)
	}
}

func TestIntegerMustEval(t *testing.T) {
	tr := assert.New(t)
	var envr core.Env

	for _, tc := range []struct {
		input string
		want  typ.Integer
	}{
		{"42", typ.Integer(42)},
		{"0", typ.Integer(0)},
		{"-7", typ.Integer(-7)},
	} {
		got := tr.MustEvalAs[typ.Integer](envr, tc.input)
		tr.Equal(got, tc.want)
		tr.Trivial[typ.Integer](envr, got, tc.input)
		tr.SelfEval(envr, got)
		tr.SelfEval(envr, tc.want)
	}
}

func TestNilMustEval(t *testing.T) {
	tr := assert.New(t)
	var envr core.Env

	got := tr.MustEvalAs[typ.Nil](envr, "nil")
	tr.Trivial[typ.Nil](envr, got, "nil")
	tr.SelfEval(envr, got)
}

func TestStringMustEval(t *testing.T) {
	tr := assert.New(t)
	var envr core.Env

	for _, tc := range []struct {
		input string
		want  typ.String
	}{
		{`"hello"`, typ.String("hello")},
		{`""`, typ.String("")},
		{`"a b"`, typ.String("a b")},
	} {
		got := tr.MustEvalAs[typ.String](envr, tc.input)
		tr.Equal(got, tc.want)
		tr.Trivial[typ.String](envr, got, tc.input)
		tr.SelfEval(envr, got)
		tr.SelfEval(envr, tc.want)
	}

	// String prints quoted.
	s := typ.String("hi")
	tr.Trivial[typ.String](envr, s, `"hi"`)
	tr.SelfEval(envr, s)
}

func TestFatalTrivial(t *testing.T) {
	tr := assert.New(t)

	f := typ.Fatalf("oops")
	tr.Repr[typ.Fatal](f, `(fatal "oops")`)
	if !f.NotNil() {
		tr.Fatalf("expected f.NotNil() to be true (f=%v)", f)
	}
	var empty typ.Fatal
	if empty.NotNil() {
		tr.Fatalf("expected f.NotNil() to be false (f=%v)", f)
	}

	qualified := f.Qualify("test")
	if !qualified.NotNil() {
		tr.Fatalf("expected qualified.NotNil() to be true (qualified=%v)", qualified)
	}
	tr.StringContains(qualified.Error(), "test: oops")
	tr.StringContains(qualified.String(), "test: oops")

	var envr core.Env
	tr.Run("Fatal represents itself as a fatal list", func(tr assert.T) {
		tr.StringRepr(f.Eval(envr), `(fatal "oops")`)
		tr.StringRepr(qualified.Eval(envr), `(fatal "test: oops")`)
	})
}
