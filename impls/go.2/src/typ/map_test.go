package typ_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func TestMapLitEvalChecksEvenLength(t *testing.T) {
	for _, literal := range []typ.MapLit{
		{},
		{typ.Integer(1), typ.Integer(2)},
	} {
		got := literal.Eval(core.Env{})
		if ut.Is[typ.Fatal](got) {
			t.Errorf("Eval() for %s returned an error", literal.String())
		}
	}

	requireError(t, typ.MapLit{typ.Integer(1)}.Eval(core.Env{}))
}

func TestMapLitValidateEvenRejectsOddLength(t *testing.T) {
	if err := (typ.MapLit{typ.Integer(1)}).ValidateEven(); err == nil {
		t.Fatal("ValidateEven() accepted a map literal with an odd number of elements")
	}
}

func TestMapLitEvalAcceptsOnlySupportedKeyTypes(t *testing.T) {
	for _, key := range []core.Value{
		typ.Integer(1),
		typ.String("key"),
		//TODO: add Quote type or proper list application support to test symbols as map keys.
		// typ.List{typ.Symbol("quote"), typ.Symbol("sym")},
		typ.Nil{},
	} {
		got := typ.MapLit{key, typ.Integer(2)}.Eval(core.Env{})
		if ut.Is[typ.Fatal](got) {
			t.Errorf("Eval() rejected valid key %T: %s", key, got)
		}
	}

	for _, key := range []core.Value{
		typ.List{typ.Integer(1)},
		typ.Vector{typ.Integer(1)},
		typ.Map{typ.Integer(1): typ.Integer(2)},
		typ.MapLit{typ.Integer(1), typ.Integer(2)},
		typ.Function(func(a typ.Integer, b typ.Integer) typ.Integer { return a + b }),
	} {
		got := typ.MapLit{key, typ.Integer(2)}.Eval(core.Env{})
		if !ut.Is[typ.Fatal](got) {
			t.Errorf("Eval() accepted invalid key type %T", key)
		}
		requireError(t, got)
	}
}

func TestMapLitEvalEvaluatesKeysAndValuesInOrder(t *testing.T) {
	var order []string
	literal := typ.MapLit{
		evalOrder{"first key", typ.Integer(1), &order},
		evalOrder{"first value", typ.Integer(10), &order},
		evalOrder{"second key", typ.Integer(2), &order},
		evalOrder{"second value", typ.Integer(20), &order},
	}

	got := literal.Eval(core.Env{})
	if ut.Is[typ.Fatal](got) {
		t.Fatalf("Eval() returned error: %s", got.String())
	}
	if len(order) != 4 || order[0] != "first key" || order[1] != "first value" ||
		order[2] != "second key" || order[3] != "second value" {
		msg := "evaluation order = %v, want [first key first value second key second value]"
		t.Errorf(msg, order)
	}

	values, ok := got.(typ.Map)
	if !ok {
		t.Fatalf("Eval() type = %T, want Map", got)
	}
	if values[typ.Integer(1)] != typ.Integer(10) || values[typ.Integer(2)] != typ.Integer(20) {
		t.Errorf("Eval() = %s, want {1 10 2 20}", values.String())
	}
}

func TestMapLitEvalPropagatesFatalFromKeyAndValue(t *testing.T) {
	fatal := typ.Fatalf("nested failure")
	envr := core.Env{"failure": fatal}

	got := typ.MapLit{typ.Symbol("failure"), typ.Integer(1)}.Eval(envr)
	requireError(t, got)
	if got.String() != fatal.String() {
		t.Errorf("key fatal = %s, want %s", got.String(), fatal.String())
	}

	got = typ.MapLit{typ.Integer(1), typ.Symbol("failure")}.Eval(envr)
	requireError(t, got)
	if got.String() != fatal.String() {
		t.Errorf("value fatal = %s, want %s", got.String(), fatal.String())
	}
}

func TestMapEvalReturnsItself(t *testing.T) {
	original := typ.Map{typ.Integer(1): typ.Integer(2)}
	got, ok := original.Eval(core.Env{}).(typ.Map)
	if !ok {
		t.Fatalf("Eval() type = %T, want Map", original.Eval(core.Env{}))
	}

	got[typ.Integer(3)] = typ.Integer(4)
	if original[typ.Integer(3)] != typ.Integer(4) {
		t.Error("Eval() returned a different map")
	}
}
