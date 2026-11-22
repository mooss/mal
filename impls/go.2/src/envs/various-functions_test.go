//revive:disable:cognitive-complexity

package envs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
	"github.com/mooss/mal/src/typ"
)

func TestEqualOp(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[bool]{
		{Name: "int same value", Input: `(= 1 1)`, Want: true},
		{Name: "int different value", Input: `(= 1 2)`, Want: false},
		{Name: "negative int same", Input: `(= -7 -7)`, Want: true},
		{Name: "string same value", Input: `(= "foo" "foo")`, Want: true},
		{Name: "string different value", Input: `(= "foo" "bar")`, Want: false},
		{Name: "empty string equal", Input: `(= "" "")`, Want: true},
		{Name: "bool true equal", Input: `(= true true)`, Want: true},
		{Name: "bool false equal", Input: `(= false false)`, Want: true},
		{Name: "bool true vs false", Input: `(= true false)`, Want: false},
		{Name: "nil equal", Input: `(= nil nil)`, Want: true},
		{Name: "empty list equal", Input: `(= (list) (list))`, Want: true},
		{Name: "list same values", Input: `(= (list 1 2) (list 1 2))`, Want: true},
		{Name: "list different values", Input: `(= (list 1 2) (list 1 3))`, Want: false},
		{Name: "list different length", Input: `(= (list 1) (list 1 2))`, Want: false},
		{Name: "empty vector equal", Input: `(= [] [])`, Want: true},
		{Name: "vector same values", Input: `(= [1 2] [1 2])`, Want: true},
		{Name: "vector different values", Input: `(= [1 2] [1 3])`, Want: false},
		{Name: "empty map equal", Input: `(= {} {})`, Want: true},
		{Name: "map same values", Input: `(= {"a" 1} {"a" 1})`, Want: true},
		{Name: "map different values", Input: `(= {"a" 1} {"a" 2})`, Want: false},
		{Name: "cross int vs string", Input: `(= 1 "1")`, Want: false},
		{Name: "cross int vs bool", Input: `(= 1 true)`, Want: false},
		{Name: "cross int vs nil", Input: `(= 0 nil)`, Want: false},
		{Name: "cross string vs bool", Input: `(= "true" true)`, Want: false},
		{Name: "cross string vs nil", Input: `(= "" nil)`, Want: false},
		{Name: "cross bool vs nil true", Input: `(= true nil)`, Want: false},
		{Name: "cross bool vs nil false", Input: `(= false nil)`, Want: false},
		{Name: "cross list vs vector", Input: `(= (list 1 2) [1 2])`, Want: true},
		{Name: "cross list vs nil", Input: `(= (list 1) nil)`, Want: false},
		{Name: "cross vector vs map", Input: `(= [1] {"a" 1})`, Want: false},
		{Name: "cross vector vs nil", Input: `(= [] nil)`, Want: false},
		{Name: "cross map vs nil", Input: `(= {} nil)`, Want: false},
		{Name: "cross int vs list", Input: `(= 1 (list 1))`, Want: false},
		{Name: "cross string vs vector", Input: `(= "a" ["a"])`, Want: false},
	}
	tr.RunCases("equal", func(tr assert.T, c tut.Case[bool]) {
		tr.Eval(envs.MalStandard, c.Input, typ.Bool(c.Want))
	}, cases...)

	//nolint:dupword
	tr.EvalFatal(envs.MalFun, `(=)`, `(= 1)`, `(= 1 2 3)`, `(= nil nil nil)`)
}

func TestReadString(t *testing.T) {
	tr := assert.New(t)
	tr.Run("one form only", func(tr assert.T) {
		cases := []tut.Case[string]{
			{Name: "integer", Input: `(read-string "1")`, Want: `1`},
			{Name: "string", Input: `(read-string "\"hi\"")`, Want: `"hi"`},
			{Name: "list form", Input: `(read-string "(1 2)")`, Want: `(list 1 2)`},
			{Name: "vector form", Input: `(read-string "[1 2]")`, Want: `[1 2]`},
		}
		tr.RunCases("one form only", func(tr assert.T, c tut.Case[string]) {
			tr.EvalSame(envs.MalFun, c.Input, c.Want)
		}, cases...)
	})

	tr.Run("multi-form empty or invalid return Fatal", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun,
			`(read-string "1 2")`,
			`(read-string "")`,
			`(read-string "   ")`,
			`(read-string "(1 2")`,
			`(read-string "[")`,
		)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(read-string)`, `(read-string "1" "2")`)
	})
}

func TestSlurp(t *testing.T) {
	tr := assert.New(t)
	tr.Run("round-trip file content", func(tr assert.T) {
		dir := tr.TempDir()
		path := filepath.Join(dir, "hello.txt")
		content := "hello world\nsecond line"
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			tr.Fatalf("WriteFile() error = %v", err)
		}
		input := fmt.Sprintf(`(slurp %q)`, path)
		got := tr.MustEvalAs[typ.String](envs.MalFun, input)
		tr.Equal(got.Unbox(), content)
	})

	tr.Run("missing file returns Fatal without panic", func(tr assert.T) {
		dir := tr.TempDir()
		missing := filepath.Join(dir, "does-not-exist.txt")
		input := fmt.Sprintf(`(slurp %q)`, missing)
		tr.EvalFatal(envs.MalFun, input)
	})

	tr.Run("arity", func(tr assert.T) {
		tr.EvalFatal(envs.MalFun, `(slurp)`, fmt.Sprintf(`(slurp %q %q)`, "a", "b"))
	})
}
