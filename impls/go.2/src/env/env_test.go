package env

import (
	"testing"

	"github.com/mooss/mal/src/core"
)

func TestEvalIntegerReturnsItself(t *testing.T) {
	got := Eval(core.Integer(42), Initial)
	if got != core.Type(core.Integer(42)) {
		t.Fatalf("Eval(42) = %v, want 42", got)
	}
}

func TestEvalEmptyListReturnsItself(t *testing.T) {
	got := Eval(core.List{}, Initial)
	lst, ok := got.(core.List)
	if !ok {
		t.Fatalf("Eval(empty list) type = %T, want core.List", got)
	}
	if len(lst) != 0 {
		t.Fatalf("Eval(empty list) = %v, want empty list", lst)
	}
}

func TestEvalSymbolLookup(t *testing.T) {
	for _, sym := range []string{"+", "-", "*", "/"} {
		got := Eval(core.Symbol(sym), Initial)
		if got.Type() != core.FUNCTION {
			t.Fatalf("Eval(%q).Type() = %s, want %s", sym, got.Type(), core.FUNCTION)
		}
	}
}

func TestEvalUnknownSymbolReturnsError(t *testing.T) {
	got := Eval(core.Symbol("nope"), Initial)
	if got.Type() != core.ERROR {
		t.Fatalf("Eval(unknown symbol).Type() = %s, want %s", got.Type(), core.ERROR)
	}
}

func TestEvalUsesGivenEnv(t *testing.T) {
	custom := Env{
		"foo": func(_, _ core.Integer) core.Integer { return 99 },
	}

	got := Eval(core.Symbol("foo"), custom)
	if got.Type() != core.FUNCTION {
		t.Fatalf("Eval(foo) with custom env type = %s, want %s", got.Type(), core.FUNCTION)
	}

	got = Eval(core.Symbol("+"), custom)
	if got.Type() != core.ERROR {
		t.Fatalf("Eval(+) with custom env type = %s, want %s", got.Type(), core.ERROR)
	}

	got = Eval(core.List{core.Symbol("foo"), core.Integer(1), core.Integer(2)}, custom)
	if got != core.Type(core.Integer(99)) {
		t.Fatalf("Eval((foo 1 2)) = %v, want 99", got)
	}
}

func TestEvalArithmetic(t *testing.T) {
	tests := []struct {
		name string
		op   string
		a, b core.Integer
		want core.Integer
	}{
		{"add", "+", 1, 2, 3},
		{"subtract", "-", 5, 3, 2},
		{"multiply", "*", 4, 3, 12},
		{"divide", "/", 8, 2, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast := core.List{core.Symbol(tt.op), tt.a, tt.b}
			got := Eval(ast, Initial)
			if got != core.Type(tt.want) {
				t.Fatalf("Eval((%s %d %d)) = %v, want %d", tt.op, tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestEvalNestedArithmetic(t *testing.T) {
	// (+ (* 2 3) (- 10 4)) == 12
	ast := core.List{
		core.Symbol("+"),
		core.List{core.Symbol("*"), core.Integer(2), core.Integer(3)},
		core.List{core.Symbol("-"), core.Integer(10), core.Integer(4)},
	}
	got := Eval(ast, Initial)
	if got != core.Type(core.Integer(12)) {
		t.Fatalf("Eval nested = %v, want 12", got)
	}
}
