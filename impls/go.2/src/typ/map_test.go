package typ_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func TestNewMapFormRejectsOddLength(t *testing.T) {
	if _, err := typ.NewMapForm(typ.Integer(1)); err == nil {
		t.Fatal("NewMapForm() accepted an odd number of elements")
	}

	form, err := typ.NewMapForm(typ.Integer(1), typ.Integer(2))
	if err != nil {
		t.Fatalf("NewMapForm() rejected an even number of elements: %v", err)
	}
	if len(form) != 1 || form[0].Key != typ.Integer(1) || form[0].Value != typ.Integer(2) {
		t.Errorf("NewMapForm() = %#v, want [{1 2}]", form)
	}
}

func TestMapFormValuesRoundTrip(t *testing.T) {
	vals := []core.Value{typ.Integer(1), typ.Integer(2), typ.String("a"), typ.Nil{}}
	form, err := typ.NewMapForm(vals...)
	if err != nil {
		t.Fatalf("NewMapForm() returned error: %v", err)
	}
	got := form.Flatten()
	if len(got) != len(vals) {
		t.Fatalf("Values() length = %d, want %d", len(got), len(vals))
	}
	for i := range vals {
		if got[i] != vals[i] {
			t.Errorf("Values()[%d] = %#v, want %#v", i, got[i], vals[i])
		}
	}
	if got := (typ.MapForm{}).Flatten(); len(got) != 0 {
		t.Errorf("Flatten() on empty form = %#v, want empty", got)
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
		got := typ.MapForm{{Key: key, Value: typ.Integer(2)}}.Eval(core.Env{})
		if ut.Is[typ.Fatal](got) {
			t.Errorf("Eval() rejected valid key %T: %s", key, got)
		}
	}

	for _, key := range []core.Value{
		typ.List{typ.Integer(1)},
		typ.VectorForm{typ.Integer(1)},
		typ.Map{typ.Integer(1): typ.Integer(2)},
		typ.MapForm{{Key: typ.Integer(1), Value: typ.Integer(2)}},
		typ.Function(func(...core.Value) core.Value { return nil }),
	} {
		got := typ.MapForm{{Key: key, Value: typ.Integer(2)}}.Eval(core.Env{})
		if !ut.Is[typ.Fatal](got) {
			t.Errorf("Eval() accepted invalid key type %T", key)
		}
		requireError(t, got)
	}
}

func TestMapLitEvalEvaluatesKeysAndValuesInOrder(t *testing.T) {
	var order []string
	literal := typ.MapForm{{
		Key:   evalOrder{"first key", typ.Integer(1), &order},
		Value: evalOrder{"first value", typ.Integer(10), &order}}, {
		Key:   evalOrder{"second key", typ.Integer(2), &order},
		Value: evalOrder{"second value", typ.Integer(20), &order}},
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
		t.Fatalf("Eval() type = %T, want %T", got, typ.Map{})
	}
	if values[typ.Integer(1)] != typ.Integer(10) || values[typ.Integer(2)] != typ.Integer(20) {
		t.Errorf("Eval() = %s, want {1 10 2 20}", values.String())
	}
}

func TestMapLitEvalPropagatesFatalFromKeyAndValue(t *testing.T) {
	fatal := typ.Fatalf("nested failure")
	envr := core.Env{"failure": fatal}

	got := typ.MapForm{{Key: typ.Symbol("failure"), Value: typ.Integer(1)}}.Eval(envr)
	requireError(t, got)
	if got.String() != fatal.String() {
		t.Errorf("key fatal = %s, want %s", got.String(), fatal.String())
	}

	got = typ.MapForm{{Key: typ.Integer(1), Value: typ.Symbol("failure")}}.Eval(envr)
	requireError(t, got)
	if got.String() != fatal.String() {
		t.Errorf("value fatal = %s, want %s", got.String(), fatal.String())
	}
}

func TestMapEvalReturnsItself(t *testing.T) {
	original := typ.Map{typ.Integer(1): typ.Integer(2)}
	got, ok := original.Eval(core.Env{}).(typ.Map)
	if !ok {
		t.Fatalf("Eval() type = %T, want %T", original.Eval(core.Env{}), typ.Map{})
	}

	got[typ.Integer(3)] = typ.Integer(4)
	if original[typ.Integer(3)] != typ.Integer(4) {
		t.Error("Eval() returned a different map")
	}
}
