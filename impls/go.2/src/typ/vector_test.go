package typ_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

func TestVectorEvalEvaluatesElementsInOrder(t *testing.T) {
	envr := core.Env{
		"first":  typ.Integer(10),
		"second": typ.Integer(20),
		"third":  typ.Integer(30),
	}
	got := typ.Vector{typ.Symbol("first"), typ.Symbol("second"), typ.Symbol("third")}.Eval(envr)

	values, ok := got.(typ.Vector)
	if !ok {
		t.Fatalf("Eval() type = %T, want Vector", got)
	}
	if len(values) != 3 {
		t.Fatalf("Eval() length = %d, want 3", len(values))
	}
	if values.String() != "[10 20 30]" {
		t.Errorf("Eval() = %s, want [10 20 30]", values.String())
	}
}
func TestVectorEvalPropagatesFirstFatalAndStops(t *testing.T) {
	firstFatal := typ.Fatalf("first failure")
	secondFatal := typ.Fatalf("second failure")
	count := 0
	envr := core.Env{
		"first":  firstFatal,
		"second": secondFatal,
	}
	got := typ.Vector{
		typ.Symbol("first"),
		typ.Symbol("second"),
		evalCounter{count: &count},
	}.Eval(envr)

	requireError(t, got)
	if got.String() != firstFatal.String() {
		t.Errorf("Eval() = %s, want first propagated fatal %s", got.String(), firstFatal.String())
	}
	if count != 0 {
		t.Errorf("elements evaluated after fatal = %d, want 0", count)
	}
}
