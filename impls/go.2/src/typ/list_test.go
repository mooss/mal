package typ_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

func TestListEvalRejectsNonFunctionHead(t *testing.T) {
	requireError(t, typ.List{typ.Integer(1), typ.Integer(2), typ.Integer(3)}.Eval(core.Env{}))
}

func TestListEvalPropagatesFatalFromHeadAndArguments(t *testing.T) {
	fatal := typ.Fatalf("nested failure")
	envr := core.Env{"failure": fatal}
	fun := typ.Function(func(typ.Integer, typ.Integer) typ.Integer { return 0 })

	for _, lst := range []typ.List{
		{typ.Symbol("failure"), typ.Integer(1), typ.Integer(2)},
		{fun, typ.Symbol("failure"), typ.Integer(2)},
		{fun, typ.Integer(1), typ.Symbol("failure")},
	} {
		got := lst.Eval(envr)
		requireError(t, got)
		if got.String() != fatal.String() {
			t.Errorf("Eval() = %s, want propagated fatal %s", got.String(), fatal.String())
		}
	}
}
