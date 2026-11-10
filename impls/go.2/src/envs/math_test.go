package envs_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func TestEvalIntegerReturnsItself(t *testing.T) {
	got := typ.Integer(42).Eval(envs.Math)
	if got != core.Value(typ.Integer(42)) {
		t.Fatalf("Eval(42) = %v, want 42", got)
	}
}

func TestEvalEmptyListReturnsItself(t *testing.T) {
	got := typ.List{}.Eval(envs.Math)
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
		got := typ.Symbol(sym).Eval(envs.Math)
		if !ut.Is[typ.Function](got) {
			t.Fatalf("Eval(%q) = %#v, want %T", sym, got, typ.Function(nil))
		}
	}
}

func TestEvalUnknownSymbolReturnsError(t *testing.T) {
	got := typ.Symbol("nope").Eval(envs.Math)
	if !ut.Is[typ.Fatal](got) {
		t.Fatalf("Eval(unknown symbol) = %#v, want %T", got, typ.Fatal{})
	}
}

func TestEvalUsesGivenEnv(t *testing.T) {
	custom := core.NewEnv()
	custom.Set("foo", typ.Function(func(...core.Value) core.Value { return typ.Integer(99) }))

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
			got := ast.Eval(envs.Math)
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
	got := ast.Eval(envs.Math)
	if got != core.Value(typ.Integer(12)) {
		t.Fatalf("Eval nested = %v, want 12", got)
	}
}

func TestVariadicZeroArgsReturnsFatal(t *testing.T) {
	for _, op := range []string{"+", "-", "*", "/"} {
		input := "(" + op + ")"
		got := tut.MustEval(envs.Math, input)
		if !ut.Is[typ.Fatal](got) {
			t.Fatalf("MustEval(%q) = %#v, want %T", input, got, typ.Fatal{})
		}
	}
}

func TestVariadicOneArg(t *testing.T) {
	// + and * with one arg return identity.
	for _, tc := range []struct {
		input string
		want  typ.Integer
	}{
		{"(+ 42)", 42},
		{"(+ -7)", -7},
		{"(+ 0)", 0},
		{"(* 42)", 42},
		{"(* -7)", -7},
		{"(* 0)", 0},
	} {
		got := tut.MustEval(envs.Math, tc.input)
		if got != core.Value(tc.want) {
			t.Fatalf("MustEval(%q) = %v, want %d", tc.input, got, tc.want)
		}
	}

	// - with one arg negates.
	for _, tc := range []struct {
		input string
		want  typ.Integer
	}{
		{"(- 5)", -5},
		{"(- -5)", 5},
		{"(- 0)", 0},
	} {
		got := tut.MustEval(envs.Math, tc.input)
		if got != core.Value(tc.want) {
			t.Fatalf("MustEval(%q) = %v, want %d", tc.input, got, tc.want)
		}
	}

	// / with one arg is an arity error.
	for _, input := range []string{"(/ 42)", "(/ -7)", "(/ 0)"} {
		got := tut.MustEval(envs.Math, input)
		if !ut.Is[typ.Fatal](got) {
			t.Fatalf("MustEval(%q) = %#v, want %T", input, got, typ.Fatal{})
		}
	}
}

func TestVariadicManyArgsLeftFold(t *testing.T) {
	tests := []struct {
		input string
		want  typ.Integer
	}{
		// + variadic.
		{"(+ 1 2 3)", 6},
		{"(+ 1 2 3 4)", 10},
		{"(+ 1 2 3 4 5)", 15},
		{"(+ -1 2 -3 4)", 2},
		// * variadic.
		{"(* 2 3 4)", 24},
		{"(* 2 3 4 5)", 120},
		{"(* 1 2 3 4 5)", 120},
		{"(* -2 3 4)", -24},
		// - left fold: (- a b c) == ((a - b) - c).
		{"(- 10 3 2)", 5},
		{"(- 20 5 3 2)", 10},
		{"(- 30 10 5 3 2)", 10},
		{"(- 100 20 30)", 50},
		{"(- 0 5 3)", -8},
		// / left fold: (/ a b c) == ((a / b) / c).
		{"(/ 100 5 2)", 10},
		{"(/ 120 6 2)", 10},
		{"(/ 120 3 2 2)", 10},
		{"(/ 1000 10 5 2 2)", 5},
		{"(/ 64 2 2 2)", 8},
	}

	for _, tt := range tests {
		got := tut.MustEval(envs.Math, tt.input)
		if got != core.Value(tt.want) {
			t.Fatalf("MustEval(%q) = %v, want %d", tt.input, got, tt.want)
		}
	}
}

func TestMixedIntStringFatalReturnsFatal(t *testing.T) {
	inputs := []string{
		// Int mixed with string.
		`(+ 1 "hi")`,
		`(+ "hi" 1)`,
		`(+ "a" "b")`,
		`(+ 1 2 "three")`,
		`(+ 1 2 3 4 "five")`,
		`(- 5 "hi")`,
		`(- "hi" 5)`,
		`(- 10 3 "x")`,
		`(- 20 5 3 "x")`,
		`(* 2 "hi")`,
		`(* "hi" 2)`,
		`(* 2 3 "hi")`,
		`(* 2 3 4 5 "hi")`,
		`(/ 8 "hi")`,
		`(/ "hi" 8)`,
		`(/ 100 5 "hi")`,
		`(/ 100 5 2 "hi")`,
		// Int mixed with Fatal from nested failing calls.
		`(+ 1 (/ 1))`,
		`(+ 1 (+))`,
		`(+ 1 (-))`,
		`(+ 1 (*))`,
		`(+ 1 (/))`,
		`(- 5 (+))`,
		`(- 10 3 (/ 1))`,
		`(* 2 (+))`,
		`(* 2 3 (/ 1))`,
		`(/ 10 (+))`,
		`(/ 10 (/) 2)`,
		`(+ (+ "a" 1) 2)`,
		`(* (- "a" 1) 2)`,
	}

	for _, input := range inputs {
		got := tut.MustEval(envs.Math, input)
		if !ut.Is[typ.Fatal](got) {
			t.Fatalf("MustEval(%q) = %#v, want %T", input, got, typ.Fatal{})
		}
	}
}

func TestDeepNestedCombinations(t *testing.T) {
	tests := []struct {
		input string
		want  typ.Integer
	}{
		{"(+ (* 2 3) (- 10 4))", 12},
		{"(+ (* 2 3) (- 10 4) (/ 20 4))", 17},
		{"(* (+ 1 2) (+ 3 4))", 21},
		{"(* (+ 1 2 3) (- 10 (- 5 2)))", 42},
		{"(- (* (+ 1 2 3) 4) (/ 20 2))", 14},
		{"(/ (* 2 3 5) (+ 2 3))", 6},
		{"(- (+ 10 20) (* 2 3))", 24},
		{"(+ (* (- 10 3) (/ 12 4)) (- (/ 100 5 2) (+ 1 2 3)))", 25},
		{"(* (- (/ 100 2) 10) (+ (- 5 2) (/ 9 3)))", 240},
		{"(/ (- (* 10 10) 20) (+ 2 2 2 2))", 10},
		{"(+ (- (* (/ 64 2) 2) 12) (* (+ 1 1 1) (- 10 5)))", 67},
		{"(* (+ (/ 20 2) (- 8 3)) (- (/ 30 3) (+ 1 2)))", 105},
	}

	for _, tt := range tests {
		got := tut.MustEval(envs.Math, tt.input)
		if got != core.Value(tt.want) {
			t.Fatalf("MustEval(%q) = %v, want %d", tt.input, got, tt.want)
		}
	}
}
