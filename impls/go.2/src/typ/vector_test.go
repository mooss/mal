package typ_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/typ"
)

func TestVectorEvalEvaluatesElementsInOrder(t *testing.T) {
	envr := core.Bindings{
		"first":  tut.Num(10),
		"second": tut.Num(20),
		"third":  tut.Num(30),
	}.AsEnv()
	got := tut.AssertOk(t, tut.Vec(
		tut.Sym("first"), tut.Sym("second"), tut.Sym("third"),
	).Eval(envr))

	vec := tut.Is[typ.Vector](t, got)
	tut.AssertEqual(t, len(vec), 3)
	tut.AssertString(t, got, "[10 20 30]")
}

func TestVectorEvalPropagatesFirstFatalAndStops(t *testing.T) {
	firstFatal := typ.Fatalf("first failure")
	secondFatal := typ.Fatalf("second failure")
	var log []string
	envr := core.Bindings{
		"first":  firstFatal,
		"second": secondFatal,
	}.AsEnv()
	got := tut.Vec(
		tut.Sym("first"),
		tut.Sym("second"),
		tut.Probe{ID: "probe", Log: &log, Ret: tut.Num(1)},
	).Eval(envr)

	tut.AssertString(t, tut.AssertFatal(t, got), firstFatal.String())
	tut.AssertEqual(t, len(log), 0)
}
