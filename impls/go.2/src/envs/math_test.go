package envs_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
	"github.com/mooss/mal/src/typ"
)

func TestEvalIntegerReturnsItself(t *testing.T) {
	assert.MalEqual(t, typ.Integer(42).Eval(envs.Math), tut.Num(42))
}

func TestEvalEmptyListReturnsItself(t *testing.T) {
	got := assert.Is[typ.List](t, tut.Lst().Eval(envs.Math))
	assert.Equal(t, len(got), 0)
}

func TestEvalSymbolLookup(t *testing.T) {
	for _, sym := range []string{"+", "-", "*", "/"} {
		assert.Is[typ.Function](t, tut.Sym(sym).Eval(envs.Math))
	}
}

func TestEvalUnknownSymbolReturnsError(t *testing.T) {
	assert.Fatal(t, tut.Sym("nope").Eval(envs.Math))
}

func TestEvalUsesGivenEnv(t *testing.T) {
	custom := core.NewEnv()
	custom.Set("foo", typ.Function(func(...core.Value) core.Value { return tut.Num(99) }))

	assert.Is[typ.Function](t, tut.Sym("foo").Eval(custom))
	assert.Fatal(t, tut.Sym("+").Eval(custom))
	assert.MalEqual(t, tut.Lst(tut.Sym("foo"), tut.Num(1), tut.Num(2)).Eval(custom), tut.Num(99))
}

func TestEvalArithmetic(t *testing.T) {
	tut.Run(t, []tut.Case[typ.Integer]{
		{Name: "add", Input: "(+ 1 2)", Want: 3},
		{Name: "subtract", Input: "(- 5 3)", Want: 2},
		{Name: "multiply", Input: "(* 4 3)", Want: 12},
		{Name: "divide", Input: "(/ 8 2)", Want: 4},
	}, func(t *testing.T, c tut.Case[typ.Integer]) {
		assert.Eval(t, envs.Math, c.Input, c.Want)
	})
}

func TestEvalNestedArithmetic(t *testing.T) {
	// (+ (* 2 3) (- 10 4)) == 12
	assert.Eval(t, envs.Math, "(+ (* 2 3) (- 10 4))", tut.Num(12))
}

func TestVariadicZeroArgsReturnsFatal(t *testing.T) {
	assert.EvalFatal(t, envs.Math, "(+)", "(-)", "(*)", "(/)")
}

func TestVariadicOneArg(t *testing.T) {
	// + and * with one arg return identity, - negates.
	tut.Run(t, []tut.Case[typ.Integer]{
		{Name: "plus 42", Input: "(+ 42)", Want: 42},
		{Name: "plus minus", Input: "(+ -7)", Want: -7},
		{Name: "plus zero", Input: "(+ 0)", Want: 0},
		{Name: "mult 42", Input: "(* 42)", Want: 42},
		{Name: "mult minus", Input: "(* -7)", Want: -7},
		{Name: "mult zero", Input: "(* 0)", Want: 0},
		{Name: "minus 5", Input: "(- 5)", Want: -5},
		//nolint:dupword
		{Name: "minus minus", Input: "(- -5)", Want: 5},
		{Name: "minus zero", Input: "(- 0)", Want: 0},
	}, func(t *testing.T, c tut.Case[typ.Integer]) {
		assert.Eval(t, envs.Math, c.Input, c.Want)
	})

	// / with one arg is an arity error.
	assert.EvalFatal(t, envs.Math, "(/ 42)", "(/ -7)", "(/ 0)")
}

func TestVariadicManyArgsLeftFold(t *testing.T) {
	tut.Run(t, []tut.Case[typ.Integer]{
		{Name: "(+ 1 2 3)", Input: "(+ 1 2 3)", Want: 6},
		{Name: "(+ 1 2 3 4)", Input: "(+ 1 2 3 4)", Want: 10},
		{Name: "(+ 1 2 3 4 5)", Input: "(+ 1 2 3 4 5)", Want: 15},
		{Name: "(+ -1 2 -3 4)", Input: "(+ -1 2 -3 4)", Want: 2},
		{Name: "(* 2 3 4)", Input: "(* 2 3 4)", Want: 24},
		{Name: "(* 2 3 4 5)", Input: "(* 2 3 4 5)", Want: 120},
		{Name: "(* 1 2 3 4 5)", Input: "(* 1 2 3 4 5)", Want: 120},
		{Name: "(* -2 3 4)", Input: "(* -2 3 4)", Want: -24},
		{Name: "(- 10 3 2)", Input: "(- 10 3 2)", Want: 5},
		{Name: "(- 20 5 3 2)", Input: "(- 20 5 3 2)", Want: 10},
		{Name: "(- 30 10 5 3 2)", Input: "(- 30 10 5 3 2)", Want: 10},
		{Name: "(- 100 20 30)", Input: "(- 100 20 30)", Want: 50},
		{Name: "(- 0 5 3)", Input: "(- 0 5 3)", Want: -8},
		{Name: "(/ 100 5 2)", Input: "(/ 100 5 2)", Want: 10},
		{Name: "(/ 120 6 2)", Input: "(/ 120 6 2)", Want: 10},
		{Name: "(/ 120 3 2 2)", Input: "(/ 120 3 2 2)", Want: 10},
		{Name: "(/ 1000 10 5 2 2)", Input: "(/ 1000 10 5 2 2)", Want: 5},
		{Name: "(/ 64 2 2 2)", Input: "(/ 64 2 2 2)", Want: 8},
	}, func(t *testing.T, c tut.Case[typ.Integer]) {
		assert.Eval(t, envs.Math, c.Input, c.Want)
	})
}

func TestMixedIntStringFatalReturnsFatal(t *testing.T) {
	assert.EvalFatal(t, envs.Math,
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
	)
}

func TestDeepNestedCombinations(t *testing.T) {
	tut.Run(t, []tut.Case[typ.Integer]{
		{Name: "(+ (* 2 3) (- 10 4))", Input: "(+ (* 2 3) (- 10 4))", Want: 12},
		{Name: "(+ (* 2 3) (- 10 4) (/ 20 4))", Input: "(+ (* 2 3) (- 10 4) (/ 20 4))", Want: 17},
		{Name: "(* (+ 1 2) (+ 3 4))", Input: "(* (+ 1 2) (+ 3 4))", Want: 21},
		{Name: "(* (+ 1 2 3) (- 10 (- 5 2)))", Input: "(* (+ 1 2 3) (- 10 (- 5 2)))", Want: 42},
		{Name: "(- (* (+ 1 2 3) 4) (/ 20 2))", Input: "(- (* (+ 1 2 3) 4) (/ 20 2))", Want: 14},
		{Name: "(/ (* 2 3 5) (+ 2 3))", Input: "(/ (* 2 3 5) (+ 2 3))", Want: 6},
		{Name: "(- (+ 10 20) (* 2 3))", Input: "(- (+ 10 20) (* 2 3))", Want: 24},
		{Name: "(+ (* (- 10 3) (/ 12 4)) (- (/ 100 5 2) (+ 1 2 3)))",
			Input: "(+ (* (- 10 3) (/ 12 4)) (- (/ 100 5 2) (+ 1 2 3)))", Want: 25},
		{Name: "(* (- (/ 100 2) 10) (+ (- 5 2) (/ 9 3)))",
			Input: "(* (- (/ 100 2) 10) (+ (- 5 2) (/ 9 3)))", Want: 240},
		{Name: "(/ (- (* 10 10) 20) (+ 2 2 2 2))",
			Input: "(/ (- (* 10 10) 20) (+ 2 2 2 2))", Want: 10},
		{Name: "(+ (- (* (/ 64 2) 2) 12) (* (+ 1 1 1) (- 10 5)))",
			Input: "(+ (- (* (/ 64 2) 2) 12) (* (+ 1 1 1) (- 10 5)))", Want: 67},
		{Name: "(* (+ (/ 20 2) (- 8 3)) (- (/ 30 3) (+ 1 2)))",
			Input: "(* (+ (/ 20 2) (- 8 3)) (- (/ 30 3) (+ 1 2)))", Want: 105},
	}, func(t *testing.T, c tut.Case[typ.Integer]) {
		assert.Eval(t, envs.Math, c.Input, c.Want)
	})
}

func TestNumericComparisons(t *testing.T) {
	tut.Run(t, []tut.Case[typ.Bool]{
		{Name: "(< 1 2)", Input: "(< 1 2)", Want: true},
		{Name: "(< 2 1)", Input: "(< 2 1)", Want: false},
		{Name: "(< 1 1)", Input: "(< 1 1)", Want: false},
		{Name: "(< 0 0)", Input: "(< 0 0)", Want: false},
		{Name: "(< -5 0)", Input: "(< -5 0)", Want: true},
		{Name: "(< 0 -5)", Input: "(< 0 -5)", Want: false},
		{Name: "(< -10 -5)", Input: "(< -10 -5)", Want: true},
		{Name: "(< -5 -10)", Input: "(< -5 -10)", Want: false},
		{Name: "(<= 1 2)", Input: "(<= 1 2)", Want: true},
		{Name: "(<= 2 1)", Input: "(<= 2 1)", Want: false},
		{Name: "(<= 1 1)", Input: "(<= 1 1)", Want: true},
		{Name: "(<= 0 0)", Input: "(<= 0 0)", Want: true},
		{Name: "(<= -5 -5)", Input: "(<= -5 -5)", Want: true},
		{Name: "(<= -5 0)", Input: "(<= -5 0)", Want: true},
		{Name: "(<= 0 -5)", Input: "(<= 0 -5)", Want: false},
		{Name: "(> 2 1)", Input: "(> 2 1)", Want: true},
		{Name: "(> 1 2)", Input: "(> 1 2)", Want: false},
		{Name: "(> 1 1)", Input: "(> 1 1)", Want: false},
		{Name: "(> 0 0)", Input: "(> 0 0)", Want: false},
		{Name: "(> 0 -5)", Input: "(> 0 -5)", Want: true},
		{Name: "(> -5 0)", Input: "(> -5 0)", Want: false},
		{Name: "(> -5 -10)", Input: "(> -5 -10)", Want: true},
		{Name: "(> -10 -5)", Input: "(> -10 -5)", Want: false},
		{Name: "(>= 2 1)", Input: "(>= 2 1)", Want: true},
		{Name: "(>= 1 2)", Input: "(>= 1 2)", Want: false},
		{Name: "(>= 1 1)", Input: "(>= 1 1)", Want: true},
		{Name: "(>= 0 0)", Input: "(>= 0 0)", Want: true},
		{Name: "(>= -5 -5)", Input: "(>= -5 -5)", Want: true},
		{Name: "(>= 0 -5)", Input: "(>= 0 -5)", Want: true},
		{Name: "(>= -5 0)", Input: "(>= -5 0)", Want: false},
		{Name: "(< (+ 1 2) (* 2 2))", Input: "(< (+ 1 2) (* 2 2))", Want: true},
		{Name: "(> (* 3 3) (+ 4 4))", Input: "(> (* 3 3) (+ 4 4))", Want: true},
		{Name: "(<= (- 5 5) 0)", Input: "(<= (- 5 5) 0)", Want: true},
		{Name: "(>= (/ 10 2) 5)", Input: "(>= (/ 10 2) 5)", Want: true},
		{Name: "(< (- 0 5) 0)", Input: "(< (- 0 5) 0)", Want: true},
		{Name: "(> (- 0 5) 0)", Input: "(> (- 0 5) 0)", Want: false},
		{Name: "(<= (+ 1 2 3) (* 2 3))", Input: "(<= (+ 1 2 3) (* 2 3))", Want: true},
		{Name: "(>= (* 2 3) (+ 1 2 3))", Input: "(>= (* 2 3) (+ 1 2 3))", Want: true},
	}, func(t *testing.T, c tut.Case[typ.Bool]) {
		assert.Eval(t, envs.Math, c.Input, c.Want)
	})
}

func TestNumericComparisonErrors(t *testing.T) {
	assert.EvalFatal(t, envs.Math,
		`(<)`,
		`(<=)`,
		`(>)`,
		`(>=)`,
		`(< 1)`,
		`(<= 1)`,
		`(> 1)`,
		`(>= 1)`,
		`(< 1 2 3)`,
		`(<= 1 2 3)`,
		`(> 1 2 3)`,
		`(>= 1 2 3)`,
		// Wrong types.
		`(< "a" 1)`,
		`(< 1 "a")`,
		`(< "a" "b")`,
		`(<= "a" 1)`,
		`(<= 1 "a")`,
		`(> "a" 1)`,
		`(> 1 "a")`,
		`(>= "a" 1)`,
		`(>= 1 "a")`,
		// Fatal propagation from nested failing calls.
		`(< 1 (+))`,
		`(< (+ "a" 1) 2)`,
		`(<= 1 (/ 1))`,
		`(> (+ "a" 1) 2)`,
		`(>= 1 (+ 1 (/ 1)))`,
		`(< (/ 100 5 "hi") 10)`,
	)
}

func TestLogicLiterals(t *testing.T) {
	// true and false evaluate to themselves.
	tut.Run(t, []tut.Case[typ.Bool]{
		{Name: "true", Input: "true", Want: true},
		{Name: "false", Input: "false", Want: false},
	}, func(t *testing.T, c tut.Case[typ.Bool]) {
		assert.Eval(t, envs.Math, c.Input, c.Want)
	})

	// Symbol lookup returns expected types.
	for _, sym := range []string{"<", "<=", ">", ">="} {
		assert.Is[typ.Function](t, tut.Sym(sym).Eval(envs.Math))
	}
	assert.Is[typ.Function](t, tut.Sym("not").Eval(envs.Math))
	for _, sym := range []string{"true", "false"} {
		assert.Is[typ.Bool](t, tut.Sym(sym).Eval(envs.Math))
	}
}

func TestNotFunction(t *testing.T) {
	tut.Run(t, []tut.Case[typ.Bool]{
		{Name: "(not true)", Input: "(not true)", Want: false},
		{Name: "(not false)", Input: "(not false)", Want: true},
		{Name: "(not nil)", Input: "(not nil)", Want: true},
		{Name: "(not 0)", Input: "(not 0)", Want: false},
		{Name: "(not 1)", Input: "(not 1)", Want: false},
		{Name: "(not -5)", Input: "(not -5)", Want: false},
		{Name: `(not "hi")`, Input: `(not "hi")`, Want: false},
		{Name: `(not "")`, Input: `(not "")`, Want: false},
		{Name: "(not ())", Input: "(not ())", Want: false},
		{Name: "(not (< 1 2))", Input: "(not (< 1 2))", Want: false},
		{Name: "(not (< 2 1))", Input: "(not (< 2 1))", Want: true},
		{Name: "(not (> 2 1))", Input: "(not (> 2 1))", Want: false},
		{Name: "(not (> 1 2))", Input: "(not (> 1 2))", Want: true},
		{Name: "(not (not true))", Input: "(not (not true))", Want: true},
		{Name: "(not (not false))", Input: "(not (not false))", Want: false},
		{Name: "(not (not nil))", Input: "(not (not nil))", Want: false},
		{Name: "(not (not 42))", Input: "(not (not 42))", Want: true},
		{Name: "(not (+ 1 2))", Input: "(not (+ 1 2))", Want: false},
		{Name: "(not (- 5 5))", Input: "(not (- 5 5))", Want: false},
	}, func(t *testing.T, c tut.Case[typ.Bool]) {
		assert.Eval(t, envs.Math, c.Input, c.Want)
	})
}

func TestNotErrors(t *testing.T) {
	assert.EvalFatal(t, envs.Math,
		`(not)`,
		`(not true false)`,
		`(not true false true)`,
		`(not nil nil)`,
		`(not (+ "a" 1))`,
		`(not (/ 1))`,
		`(not (+))`,
		`(not (< "a" 1))`,
		`(not (not (+ "a" 1)))`,
	)
}
