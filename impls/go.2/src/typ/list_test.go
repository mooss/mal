package typ_test

import (
	"strings"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

func TestListEvalRejectsNonFunctionHead(t *testing.T) {
	requireError(t, typ.List{typ.Integer(1), typ.Integer(2), typ.Integer(3)}.Eval(core.NewEnv()))
}

func TestListEvalPropagatesFatalFromHeadAndArguments(t *testing.T) {
	fatal := typ.Fatalf("nested failure")
	envr := core.NewEnv()
	envr.Set("failure", fatal)
	fun := typ.Function(func(...core.Value) core.Value { return typ.Integer(0) })

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

// logger records evaluation order and returns val.
type logger struct {
	id  string
	log *[]string
	val core.Value
}

func (logger) Type() core.TypeName { return core.LIST }
func (p logger) String() string    { return p.id }
func (p logger) Eval(core.Env) core.Value {
	*p.log = append(*p.log, p.id)
	return p.val
}

func TestListEvalEmptyReturnsItself(t *testing.T) {
	envr := core.NewEnv()
	lst := typ.List{}
	got := lst.Eval(envr)
	gotLst, ok := got.(typ.List)
	if !ok {
		t.Fatalf("Eval() = %T (%s), want typ.List", got, got.String())
	}
	if len(gotLst) != 0 {
		t.Errorf("Eval() len = %d, want 0", len(gotLst))
	}
	if got.String() != "()" {
		t.Errorf("Eval() = %s, want ()", got.String())
	}
}

func TestListEvalFunctionEvaluatesArgsInOrder(t *testing.T) {
	envr := core.NewEnv()
	var log []string
	var gotArgs []core.Value
	fun := typ.Function(func(args ...core.Value) core.Value {
		gotArgs = args
		return typ.Integer(0)
	})
	lst := typ.List{
		fun,
		logger{id: "a", log: &log, val: typ.Integer(1)},
		logger{id: "b", log: &log, val: typ.Integer(2)},
	}
	got := lst.Eval(envr)
	if got.String() != "0" {
		t.Errorf("Eval() = %s, want 0", got.String())
	}
	if len(log) != 2 || log[0] != "a" || log[1] != "b" {
		t.Errorf("eval order = %v, want [a b]", log)
	}
	if len(gotArgs) != 2 || gotArgs[0].String() != "1" || gotArgs[1].String() != "2" {
		t.Errorf("function args = %v, want [1 2]", gotArgs)
	}
}

func TestListEvalFunctionStopsOnFirstFatal(t *testing.T) {
	fatal := typ.Fatalf("boom")
	envr := core.NewEnv()

	var log []string
	called := false
	fun := typ.Function(func(...core.Value) core.Value {
		called = true
		return typ.Integer(0)
	})
	lst := typ.List{
		fun,
		logger{id: "first", log: &log, val: fatal},
		logger{id: "second", log: &log, val: typ.Integer(1)},
	}
	got := lst.Eval(envr)
	requireError(t, got)
	if got.String() != fatal.String() {
		t.Errorf("Eval() = %s, want propagated fatal %s", got.String(), fatal.String())
	}
	if len(log) != 1 || log[0] != "first" {
		t.Errorf("eval log = %v, want [first] (stop on first fatal)", log)
	}
	if called {
		t.Error("function was called despite fatal argument")
	}

	// Fatal in second position: first is evaluated, function still not called.
	var log2 []string
	called = false
	lst2 := typ.List{
		fun,
		logger{id: "first", log: &log2, val: typ.Integer(1)},
		logger{id: "second", log: &log2, val: fatal},
	}
	got2 := lst2.Eval(envr)
	requireError(t, got2)
	if got2.String() != fatal.String() {
		t.Errorf("Eval() = %s, want propagated fatal %s", got2.String(), fatal.String())
	}
	if len(log2) != 2 {
		t.Errorf("eval log = %v, want 2 evaluations before fatal", log2)
	}
	if called {
		t.Error("function was called despite fatal argument")
	}
}

func TestListEvalSpecialReceivesEnvAndRawArgs(t *testing.T) {
	envr := core.NewEnv()
	envr.Set("x", typ.Integer(42))
	var gotEnv core.Env
	var gotArgs []core.Value
	spe := typ.Special(func(e core.Env, args ...core.Value) core.Value {
		gotEnv = e
		gotEnv.Set("SENTINEL", typ.Nil{})
		gotArgs = append([]core.Value(nil), args...)
		return typ.Integer(0)
	})
	rawFirst := typ.Symbol("x")
	rawSecond := typ.Integer(7)
	lst := typ.List{spe, rawFirst, rawSecond}
	lst.Eval(envr)

	if len(gotArgs) != 2 {
		t.Fatalf("special got %d args, want 2", len(gotArgs))
	}
	if _, ok := gotArgs[0].(typ.Symbol); !ok {
		msg := "special first arg = %T (%s), want raw typ.Symbol"
		t.Errorf(msg, gotArgs[0], gotArgs[0].String())
	}
	if gotArgs[0].String() != "x" {
		t.Errorf("special first arg = %s, want raw x (not 42)", gotArgs[0].String())
	}
	if gotArgs[1].String() != "7" {
		t.Errorf("special second arg = %s, want 7", gotArgs[1].String())
	}

	// Env is forwarded: evaluating the raw arg in the received env resolves.
	if gotEnv.Get("SENTINEL") == nil {
		t.Fatal("special did not receive env")
	}
	if evaluated := gotArgs[0].Eval(gotEnv); evaluated.String() != "42" {
		t.Errorf("raw arg Eval(env) = %s, want 42 (env not forwarded)", evaluated.String())
	}
}

func TestListEvalFatalOperatorForwards(t *testing.T) {
	fatal := typ.Fatalf("head failure")
	envr := core.NewEnv()
	envr.Set("bad", fatal)
	lst := typ.List{typ.Symbol("bad"), typ.Integer(1), typ.Integer(2)}
	got := lst.Eval(envr)
	requireError(t, got)
	if got.String() != fatal.String() {
		t.Errorf("Eval() = %s, want forwarded fatal %s", got.String(), fatal.String())
	}
}

func TestListEvalNonCallableReturnsCallableError(t *testing.T) {
	envr := core.NewEnv()
	lst := typ.List{typ.Integer(1), typ.Integer(2), typ.Integer(3)}
	got := lst.Eval(envr)
	requireError(t, got)
	if !strings.Contains(got.String(), "callable") {
		t.Errorf("Eval() = %s, want error mentioning callable", got.String())
	}
}
