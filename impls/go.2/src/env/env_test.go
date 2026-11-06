package env

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func TestEvalIntegerReturnsItself(t *testing.T) {
	got := typ.Integer(42).Eval(Initial)
	if got != core.Value(typ.Integer(42)) {
		t.Fatalf("Eval(42) = %v, want 42", got)
	}
}

func TestEvalEmptyListReturnsItself(t *testing.T) {
	got := typ.List{}.Eval(Initial)
	lst, ok := got.(typ.List)
	if !ok {
		t.Fatalf("Eval(empty list) type = %T, want %T", got, typ.List{})
	}
	if len(lst) != 0 {
		t.Fatalf("Eval(empty list) = %v, want empty list", lst)
	}
}

func TestEvalSymbolLookup(t *testing.T) {
	for _, sym := range []string{"+", "-", "*", "/"} {
		got := typ.Symbol(sym).Eval(Initial)
		if !ut.Is[typ.Function](got) {
			t.Fatalf("Eval(%q) = %#v, want %T", sym, got, typ.Function(nil))
		}
	}
}

func TestEvalUnknownSymbolReturnsError(t *testing.T) {
	got := typ.Symbol("nope").Eval(Initial)
	if !ut.Is[typ.Fatal](got) {
		t.Fatalf("Eval(unknown symbol) = %#v, want %T", got, typ.Fatal{})
	}
}

func TestEvalUsesGivenEnv(t *testing.T) {
	custom := core.Env{
		"foo": _fun(func(_, _ typ.Integer) typ.Integer { return 99 }),
	}

	got := typ.Symbol("foo").Eval(custom)
	if !ut.Is[typ.Function](got) {
		t.Fatalf("Eval(foo) with custom env = %#v, want %T", got, typ.Function(nil))
	}

	got = typ.Symbol("+").Eval(custom)
	if !ut.Is[typ.Fatal](got) {
		t.Fatalf("Eval(+) with custom env = %#v, want %T", got, typ.Fatal{})
	}

	got = typ.List{typ.Symbol("foo"), typ.Integer(1), typ.Integer(2)}.Eval(custom)
	if got != core.Value(typ.Integer(99)) {
		t.Fatalf("Eval((foo 1 2)) = %v, want 99", got)
	}
}

func TestEvalArithmetic(t *testing.T) {
	tests := []struct {
		name string
		op   string
		a, b typ.Integer
		want typ.Integer
	}{
		{"add", "+", 1, 2, 3},
		{"subtract", "-", 5, 3, 2},
		{"multiply", "*", 4, 3, 12},
		{"divide", "/", 8, 2, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast := typ.List{typ.Symbol(tt.op), tt.a, tt.b}
			got := ast.Eval(Initial)
			if got != core.Value(tt.want) {
				t.Fatalf("Eval((%s %d %d)) = %v, want %d", tt.op, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestEvalNestedArithmetic(t *testing.T) {
	// (+ (* 2 3) (- 10 4)) == 12
	ast := typ.List{
		typ.Symbol("+"),
		typ.List{typ.Symbol("*"), typ.Integer(2), typ.Integer(3)},
		typ.List{typ.Symbol("-"), typ.Integer(10), typ.Integer(4)},
	}
	got := ast.Eval(Initial)
	if got != core.Value(typ.Integer(12)) {
		t.Fatalf("Eval nested = %v, want 12", got)
	}
}
