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
	"github.com/mooss/mal/src/tut/eq"
	"github.com/mooss/mal/src/typ"
)

func strCase(tr assert.T, c tut.Case[string]) {
	tr.Eval(envs.MalFun, c.Input, tut.Str(c.Want))
}

func TestStr(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[string]{
		{Name: "no args", Input: `(str)`, Want: ""},
		{Name: "single string", Input: `(str "foo")`, Want: "foo"},
		{Name: "two strings concatenated without separator", Input: `(str "foo" "bar")`,
			Want: "foobar"},
		{Name: "strings with spaces", Input: `(str "hello" " " "world")`, Want: "hello world"},
		{Name: "empty string", Input: `(str "")`, Want: ""},
		{Name: "empty plus string", Input: `(str "" "a")`, Want: "a"},
		{Name: "single integer", Input: `(str 42)`, Want: "42"},
		{Name: "integers concatenated without separator", Input: `(str 1 2 3)`, Want: "123"},
		{Name: "string and integer", Input: `(str "a" 1)`, Want: "a1"},
		{Name: "integer and string", Input: `(str 1 "a")`, Want: "1a"},
		{Name: "multiple mixed", Input: `(str "a" 1 "b" 2)`, Want: "a1b2"},
		{Name: "negative integer", Input: `(str -7)`, Want: "-7"},
		{Name: "zero", Input: `(str 0)`, Want: "0"},
	}
	tr.RunCases("str", strCase, cases...)
}

func TestPrStr(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[string]{
		{Name: "no args", Input: `(pr-str)`, Want: ""},
		{Name: "single integer", Input: `(pr-str 42)`, Want: "42"},
		{Name: "integers joined with space", Input: `(pr-str 1 2 3)`, Want: "1 2 3"},
		{Name: "single string is readable", Input: `(pr-str "hi")`, Want: `"hi"`},
		{Name: "strings joined with space and readable",
			Input: `(pr-str "a" "b")`, Want: `"a" "b"`},
		{Name: "string and integer", Input: `(pr-str "a" 1)`, Want: `"a" 1`},
		{Name: "integer and string", Input: `(pr-str 1 "a")`, Want: `1 "a"`},
		{Name: "mixed", Input: `(pr-str "a" 1 "b" 2)`, Want: `"a" 1 "b" 2`},
		{Name: "empty string is quoted", Input: `(pr-str "")`, Want: `""`},
		{Name: "negative integer", Input: `(pr-str -7)`, Want: "-7"},
	}
	tr.RunCases("pr-str", strCase, cases...)
}

func printCase(tr assert.T, c tut.Case[string]) {
	tr.Print(envs.MalFun, c.Input, c.Want)
}

func TestPrn(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[string]{
		{Name: "no args prints newline", Input: `(prn)`, Want: "\n"},
		{Name: "single integer", Input: `(prn 1)`, Want: "1\n"},
		{Name: "integers joined with space", Input: `(prn 1 2 3)`, Want: "1 2 3\n"},
		{Name: "single string is readable", Input: `(prn "hi")`, Want: "\"hi\"\n"},
		{Name: "strings are readable", Input: `(prn "a" "b")`, Want: "\"a\" \"b\"\n"},
		{Name: "string and integer", Input: `(prn "a" 1)`, Want: "\"a\" 1\n"},
		{Name: "integer and string", Input: `(prn 1 "a")`, Want: "1 \"a\"\n"},
	}
	tr.RunCases("prn", printCase, cases...)
}

func TestPrintln(t *testing.T) {
	tr := assert.New(t)
	cases := []tut.Case[string]{
		{Name: "no args prints newline", Input: `(println)`, Want: "\n"},
		{Name: "single integer", Input: `(println 1)`, Want: "1\n"},
		{Name: "integers joined with space", Input: `(println 1 2 3)`, Want: "1 2 3\n"},
		{Name: "single string is raw", Input: `(println "hi")`, Want: "hi\n"},
		{Name: "strings are raw and joined with space", Input: `(println "a" "b")`, Want: "a b\n"},
		{Name: "string and integer raw", Input: `(println "a" 1)`, Want: "a 1\n"},
		{Name: "integer and string raw", Input: `(println 1 "a")`, Want: "1 a\n"},
		{Name: "sentence", Input: `(println "hello" "world")`, Want: "hello world\n"},
	}
	tr.RunCases("println", printCase, cases...)
}

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
		tr.Equal(string(got), content)
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
