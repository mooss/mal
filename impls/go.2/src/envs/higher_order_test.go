//revive:disable:line-length-limit
package envs_test

import (
	"testing"

	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
	"github.com/mooss/mal/src/typ"
)

func TestMap(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[bool]{
		{Name: "identity over list", Input: `(= (map (fn* [x] x) (list 1 2 3)) (list 1 2 3))`, Want: true},
		{Name: "identity over vector", Input: `(= (map (fn* [x] x) [1 2 3]) (list 1 2 3))`, Want: true},
		{Name: "empty list", Input: `(= (map (fn* [x] x) (list)) (list))`, Want: true},
		{Name: "empty vector", Input: `(= (map (fn* [x] x) []) (list))`, Want: true},
		{Name: "str over vector", Input: `(= (map str [1 2]) (list "1" "2"))`, Want: true},
		{Name: "str over list", Input: `(= (map str (list 1 2 3)) (list "1" "2" "3"))`, Want: true},
		{Name: "wrap with list", Input: `(= (map list [1 2 3]) (list (list 1) (list 2) (list 3)))`, Want: true},
		{Name: "builtin first", Input: `(= (map first [(list 1 2) (list 3 4)]) (list 1 3))`, Want: true},
		{Name: "builtin count", Input: `(= (map count [(list 1 2) [3] []]) (list 2 1 0))`, Want: true},
		{Name: "result is list even for vector input", Input: `(list? (map (fn* [x] x) [1 2 3]))`, Want: true},
		{Name: "result is list for list input", Input: `(list? (map (fn* [x] x) (list 1 2)))`, Want: true},
		{Name: "closure captures env", Input: `(= (let* [y 10] (map (fn* [x] (list x y)) [1 2])) (list (list 1 10) (list 2 10)))`, Want: true},
		{Name: "user fn transform", Input: `(= (map (fn* [x] (list x x)) [1 2]) (list (list 1 1) (list 2 2)))`, Want: true},
	}
	tr.RunCases("map", func(tr assert.T, c tut.Case[bool]) {
		tr.Eval(envs.MalStandard, c.Input, typ.Bool(c.Want))
	}, cases...)

	tr.Run("error in mapped fn propagates", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(map first [1 2])`,
			`(map (fn* [x] (throw "boom")) [1 2])`,
			`(map (fn* [x] (throw "boom")) (list 1))`,
		)
	})

	tr.Run("non-function first arg returns Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(map 1 [1 2])`,
			`(map "hi" [1 2])`,
			`(map nil [1 2])`,
			`(map [1 2] [1 2])`,
			`(map {} [1 2])`,
		)
	})

	tr.Run("non-sequence second arg returns Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(map list 1)`,
			`(map list "hi")`,
			`(map list true)`,
			`(map list {})`,
			`(map list nil)`,
		)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(map)`,
			`(map list)`,
			`(map list [1] [2])`,
			`(map list [1] [2] [3])`,
		)
	})
}

func TestApply(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[bool]{
		{Name: "list with vector", Input: `(= (apply list [1 2 3]) (list 1 2 3))`, Want: true},
		{Name: "list with list", Input: `(= (apply list (list 1 2 3)) (list 1 2 3))`, Want: true},
		{Name: "empty vector", Input: `(= (apply list []) (list))`, Want: true},
		{Name: "empty list", Input: `(= (apply list (list)) (list))`, Want: true},
		{Name: "one prepended arg plus vector", Input: `(= (apply list 0 [1 2]) (list 0 1 2))`, Want: true},
		{Name: "one prepended arg plus list", Input: `(= (apply list 0 (list 1 2)) (list 0 1 2))`, Want: true},
		{Name: "two prepended args", Input: `(= (apply list 0 1 [2 3]) (list 0 1 2 3))`, Want: true},
		{Name: "prepended arg with empty vector", Input: `(= (apply list 0 []) (list 0))`, Want: true},
		{Name: "prepended arg with empty list", Input: `(= (apply list 0 (list)) (list 0))`, Want: true},
		{Name: "str with vector", Input: `(= (apply str ["a" "b" "c"]) "abc")`, Want: true},
		{Name: "str with prepended", Input: `(= (apply str "a" ["b" "c"]) "abc")`, Want: true},
		{Name: "str with list", Input: `(= (apply str (list "a" "b")) "ab")`, Want: true},
		{Name: "concat spliced", Input: `(= (apply concat [(list 1 2) (list 3 4)]) (list 1 2 3 4))`, Want: true},
		{Name: "user fn with vector", Input: `(= (apply (fn* [a b] (list a b)) [1 2]) (list 1 2))`, Want: true},
		{Name: "user fn with prepended", Input: `(= (apply (fn* [a b c] (list a b c)) 1 [2 3]) (list 1 2 3))`, Want: true},
		{Name: "variadic user fn with rest", Input: `(= (apply (fn* [a & rest] rest) [1 2 3]) (list 2 3))`, Want: true},
		{Name: "closure captures env", Input: `(= (let* [y 10] (apply (fn* [x] (list x y)) [1])) (list 1 10))`, Want: true},
	}
	tr.RunCases("apply", func(tr assert.T, c tut.Case[bool]) {
		tr.Eval(envs.MalStandard, c.Input, typ.Bool(c.Want))
	}, cases...)

	tr.Run("error in applied fn propagates", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(apply (fn* [x] (throw "boom")) [1])`,
			`(apply first [1])`,
		)
	})

	tr.Run("non-function first arg returns Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(apply 1 [1 2])`,
			`(apply "hi" [1 2])`,
			`(apply nil [1 2])`,
			`(apply [1 2] [1 2])`,
		)
	})

	tr.Run("non-sequence last arg returns Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(apply list 1)`,
			`(apply list "hi")`,
			`(apply list true)`,
			`(apply list nil)`,
			`(apply list [1] 1)`,
			`(apply list 1 2)`,
		)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalStandard,
			`(apply)`,
			`(apply list)`,
		)
	})
}
