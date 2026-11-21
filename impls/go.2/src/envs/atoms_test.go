package envs_test

import (
	"testing"

	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
	"github.com/mooss/mal/src/typ"
)

func TestAtomCreation(t *testing.T) {
	tr := assert.New(t)
	tr.Run("creates atom holding value", func(tr assert.T) {
		cases := []tut.Case[string]{
			{Name: "int", Input: `(atom 1)`, Want: ""},
			{Name: "string", Input: `(atom "hi")`, Want: ""},
			{Name: "nil", Input: `(atom nil)`, Want: ""},
			{Name: "list", Input: `(atom (list 1 2))`, Want: ""},
		}
		tr.RunCases("creates atom holding value", func(tr assert.T, c tut.Case[string]) {
			tr.MustEvalAs[*typ.Atom](envs.MalFun, c.Input)
		}, cases...)
	})

	tr.Run("deref returns initial value", func(tr assert.T) {
		cases := []tut.Case[string]{
			{Name: "int", Input: `(deref (atom 42))`, Want: `42`},
			{Name: "string", Input: `(deref (atom "hi"))`, Want: `"hi"`},
			{Name: "nil", Input: `(deref (atom nil))`, Want: `nil`},
			{Name: "list", Input: `(deref (atom (list 1 2)))`, Want: `(list 1 2)`},
			{Name: "vector", Input: `(deref (atom [1 2]))`, Want: `[1 2]`},
			{Name: "nested atom", Input: `(atom? (deref (atom (atom 1))))`, Want: `true`},
		}
		tr.RunCases("deref returns initial value", func(tr assert.T, c tut.Case[string]) {
			tr.EvalSame(envs.MalStandard, c.Input, c.Want)
		}, cases...)
	})

	tr.Run("is atom", func(tr assert.T) {
		tr.Eval(envs.MalStandard, `(atom? (atom 1))`, typ.Bool(true))
		tr.Eval(envs.MalStandard, `(atom? (atom nil))`, typ.Bool(true))
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(atom)`, `(atom 1 2)`, `(atom nil nil)`)
	})
}

func TestAtomQ(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[bool]{
		{Name: "atom of int", Input: `(atom? (atom 1))`, Want: true},
		{Name: "atom of string", Input: `(atom? (atom "hi"))`, Want: true},
		{Name: "atom of nil", Input: `(atom? (atom nil))`, Want: true},
		{Name: "atom of list", Input: `(atom? (atom (list 1)))`, Want: true},
		{Name: "atom of vector", Input: `(atom? (atom [1 2]))`, Want: true},
		{Name: "atom of atom", Input: `(atom? (atom (atom 1)))`, Want: true},
		{Name: "int false", Input: `(atom? 1)`, Want: false},
		{Name: "string false", Input: `(atom? "hi")`, Want: false},
		{Name: "nil false", Input: `(atom? nil)`, Want: false},
		{Name: "bool false", Input: `(atom? true)`, Want: false},
		{Name: "list false", Input: `(atom? (list 1 2))`, Want: false},
		{Name: "empty list false", Input: `(atom? (list))`, Want: false},
		{Name: "vector false", Input: `(atom? [1 2])`, Want: false},
		{Name: "empty vector false", Input: `(atom? [])`, Want: false},
		{Name: "map false", Input: `(atom? {})`, Want: false},
		{Name: "deref value false", Input: `(atom? (deref (atom 1)))`, Want: false},
	}
	tr.RunCases("atom?", func(tr assert.T, c tut.Case[bool]) {
		tr.Eval(envs.MalStandard, c.Input, typ.Bool(c.Want))
	}, cases...)

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(atom?)`, `(atom? (atom 1) (atom 2))`)
	})
}

func TestDeref(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[string]{
		{Name: "int", Input: `(deref (atom 1))`, Want: `1`},
		{Name: "negative int", Input: `(deref (atom -7))`, Want: `-7`},
		{Name: "string", Input: `(deref (atom "foo"))`, Want: `"foo"`},
		{Name: "empty string", Input: `(deref (atom ""))`, Want: `""`},
		{Name: "bool", Input: `(deref (atom true))`, Want: `true`},
		{Name: "nil", Input: `(deref (atom nil))`, Want: `nil`},
		{Name: "list", Input: `(deref (atom (list 1 2)))`, Want: `(list 1 2)`},
		{Name: "vector", Input: `(deref (atom [1 2]))`, Want: `[1 2]`},
		{Name: "map", Input: `(deref (atom {"a" 1}))`, Want: `{"a" 1}`},
	}
	tr.RunCases("returns held value", func(tr assert.T, c tut.Case[string]) {
		tr.EvalSame(envs.MalStandard, c.Input, c.Want)
	}, cases...)

	tr.Run("sees reset value via let", func(tr assert.T) {
		tr.Eval(envs.MalStandard,
			`(= (let* [a (atom 1)] (do (reset! a 2) (deref a))) 2)`,
			typ.Bool(true))
	})

	tr.Run("sees swap value via let", func(tr assert.T) {
		tr.Eval(envs.MalStandard,
			`(= (let* [a (atom 1)] (do (swap! a list 2) (deref a))) (list 1 2))`,
			typ.Bool(true))
	})

	tr.Run("non-atoms return Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun,
			`(deref 1)`, `(deref "hi")`, `(deref true)`, `(deref nil)`,
			`(deref (list 1))`, `(deref [1 2])`, `(deref {})`,
		)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(deref)`, `(deref (atom 1) (atom 2))`)
	})
}

func TestResetOp(t *testing.T) {
	tr := assert.New(t)
	tr.Run("returns new value", func(tr assert.T) {
		cases := []tut.Case[string]{
			{Name: "int to int", Input: `(reset! (atom 1) 2)`, Want: `2`},
			{Name: "int to string", Input: `(reset! (atom 1) "hi")`, Want: `"hi"`},
			{Name: "string to int", Input: `(reset! (atom "hi") 42)`, Want: `42`},
			{Name: "to nil", Input: `(reset! (atom 1) nil)`, Want: `nil`},
			{Name: "nil to int", Input: `(reset! (atom nil) 1)`, Want: `1`},
			{Name: "to list", Input: `(reset! (atom 1) (list 1 2))`, Want: `(list 1 2)`},
			{Name: "to vector", Input: `(reset! (atom 1) [1 2])`, Want: `[1 2]`},
		}
		tr.RunCases("returns new value", func(tr assert.T, c tut.Case[string]) {
			tr.EvalSame(envs.MalFun, c.Input, c.Want)
		}, cases...)
	})

	tr.Run("mutates atom", func(tr assert.T) {
		tr.Eval(envs.MalStandard,
			`(= (let* [a (atom 1)] (do (reset! a 2) (deref a))) 2)`,
			typ.Bool(true))
		tr.Eval(envs.MalStandard,
			`(= (let* [a (atom 1)] (do (reset! a "hi") (deref a))) "hi")`,
			typ.Bool(true))
		tr.Eval(envs.MalStandard,
			`(= (let* [a (atom 1)] (do (reset! a nil) (deref a))) nil)`,
			typ.Bool(true))
		tr.Eval(envs.MalStandard,
			`(= (let* [a (atom 1)] (do (reset! a (list 1 2)) (deref a))) (list 1 2))`,
			typ.Bool(true))
	})

	tr.Run("non-atoms return Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun,
			`(reset! 1 2)`, `(reset! "hi" 2)`, `(reset! nil 2)`,
			`(reset! (list 1) 2)`, `(reset! [1 2] 2)`, `(reset! {} 2)`,
		)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(reset!)`, `(reset! (atom 1))`, `(reset! (atom 1) 2 3)`)
	})
}

func TestSwapOp(t *testing.T) {
	tr := assert.New(t)
	tr.Run("returns new value", func(tr assert.T) {
		cases := []tut.Case[string]{
			{Name: "list no extra args", Input: `(swap! (atom 1) list)`, Want: `(list 1)`},
			{Name: "list with extra args",
				Input: `(swap! (atom 1) list 2 3)`, Want: `(list 1 2 3)`},
			{Name: "str no extra args", Input: `(swap! (atom "a") str)`, Want: `"a"`},
			{Name: "str with extra args", Input: `(swap! (atom "a") str "b" "c")`, Want: `"abc"`},
			{Name: "str int atom", Input: `(swap! (atom 1) str "a")`, Want: `"1a"`},
			{Name: "pr-str", Input: `(swap! (atom 1) pr-str 2)`, Want: `"1 2"`},
		}
		tr.RunCases("returns new value", func(tr assert.T, c tut.Case[string]) {
			tr.EvalSame(envs.MalFun, c.Input, c.Want)
		}, cases...)
	})

	tr.Run("mutates atom", func(tr assert.T) {
		tr.Eval(envs.MalStandard,
			`(= (let* [a (atom 1)] (do (swap! a list 2 3) (deref a))) (list 1 2 3))`,
			typ.Bool(true))
		tr.Eval(envs.MalStandard,
			`(= (let* [a (atom "a")] (do (swap! a str "b") (deref a))) "ab")`,
			typ.Bool(true))
	})

	tr.Run("too few args return Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(swap!)`, `(swap! (atom 1))`)
	})

	tr.Run("non-atom first arg returns Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun,
			`(swap! 1 list)`, `(swap! "hi" list)`, `(swap! nil list)`,
			`(swap! (list 1) list)`, `(swap! [1 2] list)`,
		)
	})

	tr.Run("non-function second arg returns Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun,
			`(swap! (atom 1) 1)`, `(swap! (atom 1) "hi")`,
			`(swap! (atom 1) nil)`, `(swap! (atom 1) [1 2])`,
			`(swap! (atom 1) 1 2)`,
		)
	})
}
