package envs_test

import (
	"testing"

	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
	"github.com/mooss/mal/src/tut/eq"
	"github.com/mooss/mal/src/typ"
)

func TestCount(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[int64]{
		{Name: "nil as zero", Input: `(count nil)`, Want: 0},
		{Name: "empty list", Input: `(count (list))`, Want: 0},
		{Name: "multi-element list", Input: `(count (list 1 2 3))`, Want: 3},
		{Name: "empty vector", Input: `(count [])`, Want: 0},
		{Name: "non-empty vector", Input: `(count [1 2])`, Want: 2},
		{Name: "empty map", Input: `(count {})`, Want: 0},
		{Name: "non-empty map", Input: `(count {"a" 1 "b" 2})`, Want: 2},
	}
	tr.RunCases("count", func(tr assert.T, c tut.Case[int64]) {
		tr.Eval(envs.MalFun, c.Input, tut.Int(c.Want))
	}, cases...)

	tr.Run("non-containers return Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(count 1)`, `(count 0)`, `(count true)`, `(count false)`)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(count)`, `(count [1] [2])`, `(count nil nil)`)
	})
}

func TestEmptyQ(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[bool]{
		{Name: "nil is empty", Input: `(empty? nil)`, Want: true},
		{Name: "empty list", Input: `(empty? (list))`, Want: true},
		{Name: "non-empty list", Input: `(empty? (list 1))`, Want: false},
		{Name: "multi-element list", Input: `(empty? (list 1 2 3))`, Want: false},
		{Name: "empty vector", Input: `(empty? [])`, Want: true},
		{Name: "non-empty vector", Input: `(empty? [1])`, Want: false},
		{Name: "empty map", Input: `(empty? {})`, Want: true},
		{Name: "non-empty map", Input: `(empty? {"a" 1})`, Want: false},
	}
	tr.RunCases("empty?", func(tr assert.T, c tut.Case[bool]) {
		tr.Eval(envs.MalStandard, c.Input, typ.Bool(c.Want))
	}, cases...)

	tr.Run("non-containers return Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(empty? 1)`, `(empty? 0)`, `(empty? true)`, `(empty? false)`)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(empty?)`, `(empty? [] [])`, `(empty? nil nil)`)
	})
}

func TestListConstruction(t *testing.T) {
	tr := assert.New(t)
	tr.Run("empty list", func(tr assert.T) {
		lst := tr.MustEvalAs[typ.List](envs.MalFun, `(list)`)
		tr.Equal(len(lst), 0)
		tr.Eval(envs.MalStandard, `(= (list) (list))`, typ.Bool(true))
	})

	tr.Run("multi-element list", func(tr assert.T) {
		lst := tr.MustEvalAs[typ.List](envs.MalFun, `(list 1 2 3)`)
		tr.Equal(len(lst), 3)
		for i, src := range []string{`1`, `2`, `3`} {
			want := tr.MustEval(envs.MalFun, src)
			if !eq.Mal(lst[i], want) {
				tr.Fatalf("element %d = %v, want %v", i, lst[i], want)
			}
		}
		tr.Eval(envs.MalStandard, `(= (list 1 2 3) (list 1 2 3))`, typ.Bool(true))
		tr.Eval(envs.MalStandard, `(= (list 1 2) (list 1 2 3))`, typ.Bool(false))
	})
}

func TestIsList(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[bool]{
		{Name: "empty list true", Input: `(list? (list))`, Want: true},
		{Name: "multi-element list true", Input: `(list? (list 1 2 3))`, Want: true},
		{Name: "vector false", Input: `(list? [1 2])`, Want: false},
		{Name: "empty vector false", Input: `(list? [])`, Want: false},
		{Name: "nil false", Input: `(list? nil)`, Want: false},
		{Name: "scalar int false", Input: `(list? 1)`, Want: false},
		{Name: "scalar string false", Input: `(list? "hi")`, Want: false},
		{Name: "scalar bool false", Input: `(list? true)`, Want: false},
		{Name: "map false", Input: `(list? {})`, Want: false},
	}
	tr.RunCases("list?", func(tr assert.T, c tut.Case[bool]) {
		tr.Eval(envs.MalStandard, c.Input, typ.Bool(c.Want))
	}, cases...)

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(list?)`, `(list? (list) (list))`)
	})
}
