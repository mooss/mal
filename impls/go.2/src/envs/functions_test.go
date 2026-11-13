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

func TestStr(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[string]{
		{Name: "no args", Prog: `(str)`, Want: ""},
		{Name: "single string", Prog: `(str "foo")`, Want: "foo"},
		{Name: "two strings concatenated without separator", Prog: `(str "foo" "bar")`,
			Want: "foobar"},
		{Name: "strings with spaces", Prog: `(str "hello" " " "world")`, Want: "hello world"},
		{Name: "empty string", Prog: `(str "")`, Want: ""},
		{Name: "empty plus string", Prog: `(str "" "a")`, Want: "a"},
		{Name: "single integer", Prog: `(str 42)`, Want: "42"},
		{Name: "integers concatenated without separator", Prog: `(str 1 2 3)`, Want: "123"},
		{Name: "string and integer", Prog: `(str "a" 1)`, Want: "a1"},
		{Name: "integer and string", Prog: `(str 1 "a")`, Want: "1a"},
		{Name: "multiple mixed", Prog: `(str "a" 1 "b" 2)`, Want: "a1b2"},
		{Name: "negative integer", Prog: `(str -7)`, Want: "-7"},
		{Name: "zero", Prog: `(str 0)`, Want: "0"},
	}, func(t *testing.T, c tut.SimpleCase[string]) {
		t.Helper()
		assert.Eval(t, envs.MalFun, c.Prog, tut.Str(c.Want))
	})
}

func TestPrStr(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[string]{
		{Name: "no args", Prog: `(pr-str)`, Want: ""},
		{Name: "single integer", Prog: `(pr-str 42)`, Want: "42"},
		{Name: "integers joined with space", Prog: `(pr-str 1 2 3)`, Want: "1 2 3"},
		{Name: "single string is readable", Prog: `(pr-str "hi")`, Want: `"hi"`},
		{Name: "strings joined with space and readable", Prog: `(pr-str "a" "b")`, Want: `"a" "b"`},
		{Name: "string and integer", Prog: `(pr-str "a" 1)`, Want: `"a" 1`},
		{Name: "integer and string", Prog: `(pr-str 1 "a")`, Want: `1 "a"`},
		{Name: "mixed", Prog: `(pr-str "a" 1 "b" 2)`, Want: `"a" 1 "b" 2`},
		{Name: "empty string is quoted", Prog: `(pr-str "")`, Want: `""`},
		{Name: "negative integer", Prog: `(pr-str -7)`, Want: "-7"},
	}, func(t *testing.T, c tut.SimpleCase[string]) {
		t.Helper()
		assert.Eval(t, envs.MalFun, c.Prog, tut.Str(c.Want))
	})
}

func TestPrn(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[string]{
		{Name: "no args prints newline", Prog: `(prn)`, Want: "\n"},
		{Name: "single integer", Prog: `(prn 1)`, Want: "1\n"},
		{Name: "integers joined with space", Prog: `(prn 1 2 3)`, Want: "1 2 3\n"},
		{Name: "single string is readable", Prog: `(prn "hi")`, Want: "\"hi\"\n"},
		{Name: "strings are readable", Prog: `(prn "a" "b")`, Want: "\"a\" \"b\"\n"},
		{Name: "string and integer", Prog: `(prn "a" 1)`, Want: "\"a\" 1\n"},
		{Name: "integer and string", Prog: `(prn 1 "a")`, Want: "1 \"a\"\n"},
	}, func(t *testing.T, c tut.SimpleCase[string]) {
		t.Helper()
		assert.Print(t, envs.MalFun, c.Prog, c.Want)
	})
}

func TestPrintln(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[string]{
		{Name: "no args prints newline", Prog: `(println)`, Want: "\n"},
		{Name: "single integer", Prog: `(println 1)`, Want: "1\n"},
		{Name: "integers joined with space", Prog: `(println 1 2 3)`, Want: "1 2 3\n"},
		{Name: "single string is raw", Prog: `(println "hi")`, Want: "hi\n"},
		{Name: "strings are raw and joined with space", Prog: `(println "a" "b")`, Want: "a b\n"},
		{Name: "string and integer raw", Prog: `(println "a" 1)`, Want: "a 1\n"},
		{Name: "integer and string raw", Prog: `(println 1 "a")`, Want: "1 a\n"},
		{Name: "sentence", Prog: `(println "hello" "world")`, Want: "hello world\n"},
	}, func(t *testing.T, c tut.SimpleCase[string]) {
		t.Helper()
		assert.Print(t, envs.MalFun, c.Prog, c.Want)
	})
}

func TestEqualOp(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[bool]{
		{Name: "int same value", Prog: `(= 1 1)`, Want: true},
		{Name: "int different value", Prog: `(= 1 2)`, Want: false},
		{Name: "negative int same", Prog: `(= -7 -7)`, Want: true},
		{Name: "string same value", Prog: `(= "foo" "foo")`, Want: true},
		{Name: "string different value", Prog: `(= "foo" "bar")`, Want: false},
		{Name: "empty string equal", Prog: `(= "" "")`, Want: true},
		{Name: "bool true equal", Prog: `(= true true)`, Want: true},
		{Name: "bool false equal", Prog: `(= false false)`, Want: true},
		{Name: "bool true vs false", Prog: `(= true false)`, Want: false},
		{Name: "nil equal", Prog: `(= nil nil)`, Want: true},
		{Name: "empty list equal", Prog: `(= (list) (list))`, Want: true},
		{Name: "list same values", Prog: `(= (list 1 2) (list 1 2))`, Want: true},
		{Name: "list different values", Prog: `(= (list 1 2) (list 1 3))`, Want: false},
		{Name: "list different length", Prog: `(= (list 1) (list 1 2))`, Want: false},
		{Name: "empty vector equal", Prog: `(= [] [])`, Want: true},
		{Name: "vector same values", Prog: `(= [1 2] [1 2])`, Want: true},
		{Name: "vector different values", Prog: `(= [1 2] [1 3])`, Want: false},
		{Name: "empty map equal", Prog: `(= {} {})`, Want: true},
		{Name: "map same values", Prog: `(= {"a" 1} {"a" 1})`, Want: true},
		{Name: "map different values", Prog: `(= {"a" 1} {"a" 2})`, Want: false},
		{Name: "cross int vs string", Prog: `(= 1 "1")`, Want: false},
		{Name: "cross int vs bool", Prog: `(= 1 true)`, Want: false},
		{Name: "cross int vs nil", Prog: `(= 0 nil)`, Want: false},
		{Name: "cross string vs bool", Prog: `(= "true" true)`, Want: false},
		{Name: "cross string vs nil", Prog: `(= "" nil)`, Want: false},
		{Name: "cross bool vs nil true", Prog: `(= true nil)`, Want: false},
		{Name: "cross bool vs nil false", Prog: `(= false nil)`, Want: false},
		{Name: "cross list vs vector", Prog: `(= (list 1 2) [1 2])`, Want: true},
		{Name: "cross list vs nil", Prog: `(= (list 1) nil)`, Want: false},
		{Name: "cross vector vs map", Prog: `(= [1] {"a" 1})`, Want: false},
		{Name: "cross vector vs nil", Prog: `(= [] nil)`, Want: false},
		{Name: "cross map vs nil", Prog: `(= {} nil)`, Want: false},
		{Name: "cross int vs list", Prog: `(= 1 (list 1))`, Want: false},
		{Name: "cross string vs vector", Prog: `(= "a" ["a"])`, Want: false},
	}, func(t *testing.T, c tut.SimpleCase[bool]) {
		t.Helper()
		assert.Eval(t, envs.MalDefault, c.Prog, typ.Bool(c.Want))
	})

	t.Run("arity", func(t *testing.T) {
		//nolint:dupword
		assert.EvalFatal(t, envs.MalFun, `(=)`, `(= 1)`, `(= 1 2 3)`, `(= nil nil nil)`)
	})
}

func TestCount(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[int]{
		{Name: "nil as zero", Prog: `(count nil)`, Want: 0},
		{Name: "empty list", Prog: `(count (list))`, Want: 0},
		{Name: "multi-element list", Prog: `(count (list 1 2 3))`, Want: 3},
		{Name: "empty vector", Prog: `(count [])`, Want: 0},
		{Name: "non-empty vector", Prog: `(count [1 2])`, Want: 2},
		{Name: "empty map", Prog: `(count {})`, Want: 0},
		{Name: "non-empty map", Prog: `(count {"a" 1 "b" 2})`, Want: 2},
	}, func(t *testing.T, c tut.SimpleCase[int]) {
		t.Helper()
		assert.Eval(t, envs.MalFun, c.Prog, tut.Num(int64(c.Want)))
	})

	t.Run("non-containers return Fatal", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun,
			`(count 1)`, `(count 0)`, `(count true)`, `(count false)`)
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(count)`, `(count [1] [2])`, `(count nil nil)`)
	})
}

func TestEmptyQ(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[bool]{
		{Name: "nil is empty", Prog: `(empty? nil)`, Want: true},
		{Name: "empty list", Prog: `(empty? (list))`, Want: true},
		{Name: "non-empty list", Prog: `(empty? (list 1))`, Want: false},
		{Name: "multi-element list", Prog: `(empty? (list 1 2 3))`, Want: false},
		{Name: "empty vector", Prog: `(empty? [])`, Want: true},
		{Name: "non-empty vector", Prog: `(empty? [1])`, Want: false},
		{Name: "empty map", Prog: `(empty? {})`, Want: true},
		{Name: "non-empty map", Prog: `(empty? {"a" 1})`, Want: false},
	}, func(t *testing.T, c tut.SimpleCase[bool]) {
		t.Helper()
		assert.Eval(t, envs.MalDefault, c.Prog, typ.Bool(c.Want))
	})

	t.Run("non-containers return Fatal", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun,
			`(empty? 1)`, `(empty? 0)`, `(empty? true)`, `(empty? false)`)
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(empty?)`, `(empty? [] [])`, `(empty? nil nil)`)
	})
}

func TestListConstruction(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		lst := tut.MustEvalAs[typ.List](t, envs.MalFun, `(list)`)
		assert.Equal(t, len(lst), 0)
		assert.Eval(t, envs.MalDefault, `(= (list) (list))`, typ.Bool(true))
	})

	t.Run("multi-element list", func(t *testing.T) {
		lst := tut.MustEvalAs[typ.List](t, envs.MalFun, `(list 1 2 3)`)
		assert.Equal(t, len(lst), 3)
		for i, src := range []string{`1`, `2`, `3`} {
			want := tut.MustEval(envs.MalFun, src)
			if !eq.Mal(lst[i], want) {
				t.Fatalf("element %d = %v, want %v", i, lst[i], want)
			}
		}
		assert.Eval(t, envs.MalDefault, `(= (list 1 2 3) (list 1 2 3))`, typ.Bool(true))
		assert.Eval(t, envs.MalDefault, `(= (list 1 2) (list 1 2 3))`, typ.Bool(false))
	})
}

func TestIsList(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[bool]{
		{Name: "empty list true", Prog: `(list? (list))`, Want: true},
		{Name: "multi-element list true", Prog: `(list? (list 1 2 3))`, Want: true},
		{Name: "vector false", Prog: `(list? [1 2])`, Want: false},
		{Name: "empty vector false", Prog: `(list? [])`, Want: false},
		{Name: "nil false", Prog: `(list? nil)`, Want: false},
		{Name: "scalar int false", Prog: `(list? 1)`, Want: false},
		{Name: "scalar string false", Prog: `(list? "hi")`, Want: false},
		{Name: "scalar bool false", Prog: `(list? true)`, Want: false},
		{Name: "map false", Prog: `(list? {})`, Want: false},
	}, func(t *testing.T, c tut.SimpleCase[bool]) {
		t.Helper()
		assert.Eval(t, envs.MalDefault, c.Prog, typ.Bool(c.Want))
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(list?)`, `(list? (list) (list))`)
	})
}

func TestReadString(t *testing.T) {
	t.Run("one form only", func(t *testing.T) {
		tut.RunSimple(t, []tut.SimpleCase[string]{
			{Name: "integer", Prog: `(read-string "1")`, Want: `1`},
			{Name: "string", Prog: `(read-string "\"hi\"")`, Want: `"hi"`},
			{Name: "list form", Prog: `(read-string "(1 2)")`, Want: `(list 1 2)`},
			{Name: "vector form", Prog: `(read-string "[1 2]")`, Want: `[1 2]`},
		}, func(t *testing.T, c tut.SimpleCase[string]) {
			t.Helper()
			assert.EvalSame(t, envs.MalFun, c.Prog, c.Want)
		})
	})

	t.Run("multi-form empty or invalid return Fatal", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun,
			`(read-string "1 2")`,
			`(read-string "")`,
			`(read-string "   ")`,
			`(read-string "(1 2")`,
			`(read-string "[")`,
		)
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(read-string)`, `(read-string "1" "2")`)
	})
}

func TestSlurp(t *testing.T) {
	t.Run("round-trip file content", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "hello.txt")
		content := "hello world\nsecond line"
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		input := fmt.Sprintf(`(slurp %q)`, path)
		got := tut.MustEvalAs[typ.String](t, envs.MalFun, input)
		assert.Equal(t, string(got), content)
	})

	t.Run("missing file returns Fatal without panic", func(t *testing.T) {
		dir := t.TempDir()
		missing := filepath.Join(dir, "does-not-exist.txt")
		input := fmt.Sprintf(`(slurp %q)`, missing)
		assert.EvalFatal(t, envs.MalFun, input)
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(slurp)`, fmt.Sprintf(`(slurp %q %q)`, "a", "b"))
	})
}

func TestAtomCreation(t *testing.T) {
	t.Run("creates atom holding value", func(t *testing.T) {
		tut.RunSimple(t, []tut.SimpleCase[string]{
			{Name: "int", Prog: `(atom 1)`, Want: ""},
			{Name: "string", Prog: `(atom "hi")`, Want: ""},
			{Name: "nil", Prog: `(atom nil)`, Want: ""},
			{Name: "list", Prog: `(atom (list 1 2))`, Want: ""},
		}, func(t *testing.T, c tut.SimpleCase[string]) {
			t.Helper()
			tut.MustEvalAs[*typ.Atom](t, envs.MalFun, c.Prog)
		})
	})

	t.Run("deref returns initial value", func(t *testing.T) {
		tut.RunSimple(t, []tut.SimpleCase[string]{
			{Name: "int", Prog: `(deref (atom 42))`, Want: `42`},
			{Name: "string", Prog: `(deref (atom "hi"))`, Want: `"hi"`},
			{Name: "nil", Prog: `(deref (atom nil))`, Want: `nil`},
			{Name: "list", Prog: `(deref (atom (list 1 2)))`, Want: `(list 1 2)`},
			{Name: "vector", Prog: `(deref (atom [1 2]))`, Want: `[1 2]`},
			{Name: "nested atom", Prog: `(atom? (deref (atom (atom 1))))`, Want: `true`},
		}, func(t *testing.T, c tut.SimpleCase[string]) {
			t.Helper()
			assert.EvalSame(t, envs.MalDefault, c.Prog, c.Want)
		})
	})

	t.Run("is atom", func(t *testing.T) {
		assert.Eval(t, envs.MalDefault, `(atom? (atom 1))`, typ.Bool(true))
		assert.Eval(t, envs.MalDefault, `(atom? (atom nil))`, typ.Bool(true))
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(atom)`, `(atom 1 2)`, `(atom nil nil)`)
	})
}

func TestAtomQ(t *testing.T) {
	tut.RunSimple(t, []tut.SimpleCase[bool]{
		{Name: "atom of int", Prog: `(atom? (atom 1))`, Want: true},
		{Name: "atom of string", Prog: `(atom? (atom "hi"))`, Want: true},
		{Name: "atom of nil", Prog: `(atom? (atom nil))`, Want: true},
		{Name: "atom of list", Prog: `(atom? (atom (list 1)))`, Want: true},
		{Name: "atom of vector", Prog: `(atom? (atom [1 2]))`, Want: true},
		{Name: "atom of atom", Prog: `(atom? (atom (atom 1)))`, Want: true},
		{Name: "int false", Prog: `(atom? 1)`, Want: false},
		{Name: "string false", Prog: `(atom? "hi")`, Want: false},
		{Name: "nil false", Prog: `(atom? nil)`, Want: false},
		{Name: "bool false", Prog: `(atom? true)`, Want: false},
		{Name: "list false", Prog: `(atom? (list 1 2))`, Want: false},
		{Name: "empty list false", Prog: `(atom? (list))`, Want: false},
		{Name: "vector false", Prog: `(atom? [1 2])`, Want: false},
		{Name: "empty vector false", Prog: `(atom? [])`, Want: false},
		{Name: "map false", Prog: `(atom? {})`, Want: false},
		{Name: "deref value false", Prog: `(atom? (deref (atom 1)))`, Want: false},
	}, func(t *testing.T, c tut.SimpleCase[bool]) {
		t.Helper()
		assert.Eval(t, envs.MalDefault, c.Prog, typ.Bool(c.Want))
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(atom?)`, `(atom? (atom 1) (atom 2))`)
	})
}

func TestDeref(t *testing.T) {
	t.Run("returns held value", func(t *testing.T) {
		tut.RunSimple(t, []tut.SimpleCase[string]{
			{Name: "int", Prog: `(deref (atom 1))`, Want: `1`},
			{Name: "negative int", Prog: `(deref (atom -7))`, Want: `-7`},
			{Name: "string", Prog: `(deref (atom "foo"))`, Want: `"foo"`},
			{Name: "empty string", Prog: `(deref (atom ""))`, Want: `""`},
			{Name: "bool", Prog: `(deref (atom true))`, Want: `true`},
			{Name: "nil", Prog: `(deref (atom nil))`, Want: `nil`},
			{Name: "list", Prog: `(deref (atom (list 1 2)))`, Want: `(list 1 2)`},
			{Name: "vector", Prog: `(deref (atom [1 2]))`, Want: `[1 2]`},
			{Name: "map", Prog: `(deref (atom {"a" 1}))`, Want: `{"a" 1}`},
		}, func(t *testing.T, c tut.SimpleCase[string]) {
			t.Helper()
			assert.EvalSame(t, envs.MalDefault, c.Prog, c.Want)
		})
	})

	t.Run("sees reset value via let", func(t *testing.T) {
		assert.Eval(t, envs.MalDefault,
			`(= (let* [a (atom 1)] (do (reset! a 2) (deref a))) 2)`,
			typ.Bool(true))
	})

	t.Run("sees swap value via let", func(t *testing.T) {
		assert.Eval(t, envs.MalDefault,
			`(= (let* [a (atom 1)] (do (swap! a list 2) (deref a))) (list 1 2))`,
			typ.Bool(true))
	})

	t.Run("non-atoms return Fatal", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun,
			`(deref 1)`, `(deref "hi")`, `(deref true)`, `(deref nil)`,
			`(deref (list 1))`, `(deref [1 2])`, `(deref {})`,
		)
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(deref)`, `(deref (atom 1) (atom 2))`)
	})
}

func TestResetOp(t *testing.T) {
	t.Run("returns new value", func(t *testing.T) {
		tut.RunSimple(t, []tut.SimpleCase[string]{
			{Name: "int to int", Prog: `(reset! (atom 1) 2)`, Want: `2`},
			{Name: "int to string", Prog: `(reset! (atom 1) "hi")`, Want: `"hi"`},
			{Name: "string to int", Prog: `(reset! (atom "hi") 42)`, Want: `42`},
			{Name: "to nil", Prog: `(reset! (atom 1) nil)`, Want: `nil`},
			{Name: "nil to int", Prog: `(reset! (atom nil) 1)`, Want: `1`},
			{Name: "to list", Prog: `(reset! (atom 1) (list 1 2))`, Want: `(list 1 2)`},
			{Name: "to vector", Prog: `(reset! (atom 1) [1 2])`, Want: `[1 2]`},
		}, func(t *testing.T, c tut.SimpleCase[string]) {
			t.Helper()
			assert.EvalSame(t, envs.MalFun, c.Prog, c.Want)
		})
	})

	t.Run("mutates atom", func(t *testing.T) {
		assert.Eval(t, envs.MalDefault,
			`(= (let* [a (atom 1)] (do (reset! a 2) (deref a))) 2)`,
			typ.Bool(true))
		assert.Eval(t, envs.MalDefault,
			`(= (let* [a (atom 1)] (do (reset! a "hi") (deref a))) "hi")`,
			typ.Bool(true))
		assert.Eval(t, envs.MalDefault,
			`(= (let* [a (atom 1)] (do (reset! a nil) (deref a))) nil)`,
			typ.Bool(true))
		assert.Eval(t, envs.MalDefault,
			`(= (let* [a (atom 1)] (do (reset! a (list 1 2)) (deref a))) (list 1 2))`,
			typ.Bool(true))
	})

	t.Run("non-atoms return Fatal", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun,
			`(reset! 1 2)`, `(reset! "hi" 2)`, `(reset! nil 2)`,
			`(reset! (list 1) 2)`, `(reset! [1 2] 2)`, `(reset! {} 2)`,
		)
	})

	t.Run("arity", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(reset!)`, `(reset! (atom 1))`, `(reset! (atom 1) 2 3)`)
	})
}

func TestSwapOp(t *testing.T) {
	t.Run("returns new value", func(t *testing.T) {
		tut.RunSimple(t, []tut.SimpleCase[string]{
			{Name: "list no extra args", Prog: `(swap! (atom 1) list)`, Want: `(list 1)`},
			{Name: "list with extra args", Prog: `(swap! (atom 1) list 2 3)`, Want: `(list 1 2 3)`},
			{Name: "str no extra args", Prog: `(swap! (atom "a") str)`, Want: `"a"`},
			{Name: "str with extra args", Prog: `(swap! (atom "a") str "b" "c")`, Want: `"abc"`},
			{Name: "str int atom", Prog: `(swap! (atom 1) str "a")`, Want: `"1a"`},
			{Name: "pr-str", Prog: `(swap! (atom 1) pr-str 2)`, Want: `"1 2"`},
		}, func(t *testing.T, c tut.SimpleCase[string]) {
			t.Helper()
			assert.EvalSame(t, envs.MalFun, c.Prog, c.Want)
		})
	})

	t.Run("mutates atom", func(t *testing.T) {
		assert.Eval(t, envs.MalDefault,
			`(= (let* [a (atom 1)] (do (swap! a list 2 3) (deref a))) (list 1 2 3))`,
			typ.Bool(true))
		assert.Eval(t, envs.MalDefault,
			`(= (let* [a (atom "a")] (do (swap! a str "b") (deref a))) "ab")`,
			typ.Bool(true))
	})

	t.Run("too few args return Fatal", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun, `(swap!)`, `(swap! (atom 1))`)
	})

	t.Run("non-atom first arg returns Fatal", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun,
			`(swap! 1 list)`, `(swap! "hi" list)`, `(swap! nil list)`,
			`(swap! (list 1) list)`, `(swap! [1 2] list)`,
		)
	})

	t.Run("non-function second arg returns Fatal", func(t *testing.T) {
		assert.EvalFatal(t, envs.MalFun,
			`(swap! (atom 1) 1)`, `(swap! (atom 1) "hi")`,
			`(swap! (atom 1) nil)`, `(swap! (atom 1) [1 2])`,
			`(swap! (atom 1) 1 2)`,
		)
	})
}
