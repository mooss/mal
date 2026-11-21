//revive:disable:line-length-limit,function-length
package envs_test

import (
	"testing"

	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
	"github.com/mooss/mal/src/typ"
)

// Tests for try*/catch* and throw using envs.MalStandard and tut/assert with Mal/Lisp code.
//
//nolint:dupword
func TestTryStarSpecial(t *testing.T) {
	tr := assert.New(t)

	tr.Run("no exception returns attempt", func(tr assert.T) {
		tr.Eval(envs.MalStandard, `(try* 42 (catch* e 99))`, tut.Int(42))
		tr.Eval(envs.MalStandard, `(try* "ok" (catch* e "caught"))`, tut.Str("ok"))
		tr.Eval(envs.MalStandard, `(try* nil (catch* e 1))`, typ.Nil{})
		tr.Eval(envs.MalStandard, `(try* true (catch* e false))`, typ.Bool(true))
		tr.EvalSame(envs.MalStandard, `(try* (list 1 2) (catch* e 99))`, `(list 1 2)`)
		tr.EvalSame(envs.MalStandard, `(try* [1 2] (catch* e 99))`, `[1 2]`)
	})

	tr.Run("catches thrown values", func(tr assert.T) {
		tr.Eval(envs.MalStandard, `(try* (throw "oops") (catch* e e))`, tut.Str("oops"))
		tr.Eval(envs.MalStandard, `(try* (throw 42) (catch* e e))`, tut.Int(42))
		tr.Eval(envs.MalStandard, `(try* (throw true) (catch* e e))`, typ.Bool(true))
		tr.Eval(envs.MalStandard, `(try* (throw nil) (catch* e e))`, typ.Nil{})
		tr.EvalSame(envs.MalStandard, `(try* (throw (list 1 2)) (catch* e e))`, `(list 1 2)`)
		tr.EvalSame(envs.MalStandard, `(try* (throw [1 2]) (catch* e e))`, `[1 2]`)
		tr.EvalSame(envs.MalStandard, `(try* (throw {"a" 1}) (catch* e e))`, `{"a" 1}`)
	})

	tr.Run("handler sees bound exception", func(tr assert.T) {
		tr.Eval(envs.MalStandard, `(try* (throw "a") (catch* err (str "caught:" err)))`, tut.Str("caught:a"))
		tr.EvalSame(envs.MalStandard, `(try* (throw 2) (catch* e (list e e)))`, `(list 2 2)`)
		tr.EvalSame(envs.MalStandard, `(let* [x 10] (try* (throw 2) (catch* e (list x e))))`, `(list 10 2)`)
		tr.Eval(envs.MalStandard, `(let* [x 10] (try* (throw 1) (catch* e x)))`, tut.Int(10))
		tr.Eval(envs.MalStandard, `(let* [x 10] (try* x (catch* e 99)))`, tut.Int(10))
	})

	tr.Run("handler not evaluated on success", func(tr assert.T) {
		tr.Eval(envs.MalStandard,
			`(let* [a (atom 0)] (do (try* 1 (catch* e (reset! a 1))) (deref a)))`,
			tut.Int(0))
	})

	tr.Run("nested try", func(tr assert.T) {
		tr.Eval(envs.MalStandard, `(try* (try* (throw 1) (catch* a (throw 2))) (catch* b b))`, tut.Int(2))
		tr.Eval(envs.MalStandard, `(try* (try* 1 (catch* a 2)) (catch* b 3))`, tut.Int(1))
		tr.Eval(envs.MalStandard, `(try* (try* (throw 1) (catch* a a)) (catch* b 99))`, tut.Int(1))
	})

	tr.Run("rethrown from handler is observable", func(tr assert.T) {
		tr.Eval(envs.MalStandard, `(try* (try* (throw "inner") (catch* a (throw (str a "-outer")))) (catch* b b))`, tut.Str("inner-outer"))
	})

	tr.RunCases("arity and type errors are caught", func(tr assert.T, c string) {
		tr.Eval(envs.MalStandard, c, tut.Int(42))
	}, `(try* (map) (catch* e 42))`,
		`(try* (throw 1 2) (catch* e 42))`,
		`(try* (throw) (catch* e 42))`,
		`(try* (count 1) (catch* e 42))`,
	)

	tr.Run("exception binding does not leak", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(do (try* (throw 5) (catch* myerr myerr)) myerr)`,
		)
		tr.Eval(envs.MalStandard, `(let* [e 10] (try* (throw 5) (catch* e e)))`, tut.Int(5))
		tr.Eval(envs.MalStandard, `(let* [e 10] (do (try* (throw 5) (catch* e e)) e))`, tut.Int(10))
	})

	tr.Run("invalid catch clause returns Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(try* 1 2)`,
			`(try* 1 (list 1 2 3))`,
			`(try* 1 (catch* e))`,
			`(try* 1 (catch* e 1 2))`,
			`(try* 1 (catch* 1 2))`,
			`(try* 1 "not-a-list")`,
			`(try* 1 nil)`,
		)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(try*)`,
			`(try* 1)`,
			`(try* 1 2 3)`,
		)
	})

	tr.Run("bare catch star is not bound", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard, `(catch* e e)`)
	})
}

func TestThrowFunction(t *testing.T) {
	tr := assert.New(t)

	tr.Run("throw always returns Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(throw 1)`,
			`(throw "x")`,
			`(throw nil)`,
			`(throw true)`,
			`(throw (list 1 2))`,
			`(throw [1 2])`,
			`(throw {})`,
		)
	})

	tr.Run("thrown values are caught by try", func(tr assert.T) {
		cases := []tut.Case[string]{
			{Name: "string", Input: `(try* (throw "boom") (catch* e e))`, Want: `"boom"`},
			{Name: "int", Input: `(try* (throw 7) (catch* e e))`, Want: `7`},
			{Name: "bool", Input: `(try* (throw false) (catch* e e))`, Want: `false`},
			{Name: "nil", Input: `(try* (throw nil) (catch* e e))`, Want: `nil`},
			{Name: "list", Input: `(try* (throw (list 1 2)) (catch* e e))`, Want: `(list 1 2)`},
			{Name: "vector", Input: `(try* (throw [1 2]) (catch* e e))`, Want: `[1 2]`},
			{Name: "map", Input: `(try* (throw {"a" 1}) (catch* e e))`, Want: `{"a" 1}`},
		}
		tr.RunCases("thrown values are caught by try", func(tr assert.T, c tut.Case[string]) {
			tr.EvalSame(envs.MalStandard, c.Input, c.Want)
		}, cases...)
	})

	tr.Run("throw inside do and let", func(tr assert.T) {
		tr.Eval(envs.MalStandard, `(try* (do 1 2 (throw "late")) (catch* e e))`, tut.Str("late"))
		tr.Eval(envs.MalStandard, `(let* [x 1] (try* (throw x) (catch* e e)))`, tut.Int(1))
		tr.Eval(envs.MalStandard, `(let* [f (fn* [x] (throw x))] (try* (f 42) (catch* e e)))`, tut.Int(42))
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard, `(throw)`, `(throw 1 2)`, `(throw nil nil)`)
	})
}
