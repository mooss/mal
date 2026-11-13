//revive:disable:cognitive-complexity

package envs_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/typ"
)

type tc[T any] struct {
	name  string
	input string
	want  T
}

func runTable[T any](t *testing.T, cases []tc[T], check func(t *testing.T, c tc[T])) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			check(t, c)
		})
	}
}

func evalFun(src string) core.Value {
	return tut.MustEval(envs.MalFun, src)
}

func evalDefault(src string) core.Value {
	return tut.MustEval(envs.MalDefault, src)
}

func assertString(t *testing.T, c tc[string]) {
	t.Helper()
	got := evalFun(c.input)
	s, ok := got.(typ.String)
	if !ok {
		t.Fatalf("MustEval(%q) type = %T, want %T", c.input, got, typ.String(""))
	}
	if string(s) != c.want {
		t.Fatalf("MustEval(%q) = %q, want %q", c.input, string(s), c.want)
	}
}

func assertPrint(t *testing.T, c tc[string]) {
	t.Helper()
	var got core.Value
	output := captureStdout(t, func() {
		got = evalFun(c.input)
	})
	if _, ok := got.(typ.Nil); !ok {
		t.Fatalf("MustEval(%q) type = %T, want %T", c.input, got, typ.Nil{})
	}
	if output != c.want {
		t.Fatalf("MustEval(%q) printed %q, want %q", c.input, output, c.want)
	}
}

func assertBoolCheck(t *testing.T, c tc[bool]) {
	t.Helper()
	assertBool(t, c.input, c.want)
}

func assertCountCheck(t *testing.T, c tc[int]) {
	t.Helper()
	assertCountValue(t, c.input, c.want)
}

func assertFatalAll(t *testing.T, inputs []string) {
	t.Helper()
	for _, input := range inputs {
		assertFatal(t, input)
	}
}

func equalWith(eval func(string) core.Value) func(t *testing.T, c tc[string]) {
	return func(t *testing.T, c tc[string]) {
		t.Helper()
		got := eval(c.input)
		if isFatalValue(got) {
			t.Fatalf("MustEval(%q) = Fatal %v, want %q", c.input, got, c.want)
		}
		want := eval(c.want)
		if !got.Equal(want) {
			t.Fatalf("MustEval(%q) = %v, want %v", c.input, got, want)
		}
	}
}

func assertAtom(t *testing.T, c tc[string]) {
	t.Helper()
	got := evalFun(c.input)
	if isFatalValue(got) {
		t.Fatalf("MustEval(%q) = Fatal %v, want atom", c.input, got)
	}
	if _, ok := got.(*typ.Atom); !ok {
		t.Fatalf("MustEval(%q) type = %T, want *typ.Atom", c.input, got)
	}
}

func assertListLen(t *testing.T, input string, wantLen int) typ.List {
	t.Helper()
	got := evalFun(input)
	lst, ok := got.(typ.List)
	if !ok {
		t.Fatalf("MustEval(%q) type = %T, want %T", input, got, typ.List{})
	}
	if len(lst) != wantLen {
		t.Fatalf("MustEval(%q) length = %d, want %d", input, len(lst), wantLen)
	}
	return lst
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("pipe writer Close() error = %v", err)
	}
	os.Stdout = old

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy() error = %v", err)
	}
	if err := r.Close(); err != nil {
		t.Fatalf("pipe reader Close() error = %v", err)
	}
	return buf.String()
}

func TestStr(t *testing.T) {
	runTable(t, []tc[string]{
		{"no args", `(str)`, ""},
		{"single string", `(str "foo")`, "foo"},
		{"two strings concatenated without separator", `(str "foo" "bar")`, "foobar"},
		{"strings with spaces", `(str "hello" " " "world")`, "hello world"},
		{"empty string", `(str "")`, ""},
		{"empty plus string", `(str "" "a")`, "a"},
		{"single integer", `(str 42)`, "42"},
		{"integers concatenated without separator", `(str 1 2 3)`, "123"},
		{"string and integer", `(str "a" 1)`, "a1"},
		{"integer and string", `(str 1 "a")`, "1a"},
		{"multiple mixed", `(str "a" 1 "b" 2)`, "a1b2"},
		{"negative integer", `(str -7)`, "-7"},
		{"zero", `(str 0)`, "0"},
	}, assertString)
}

func TestPrStr(t *testing.T) {
	runTable(t, []tc[string]{
		{"no args", `(pr-str)`, ""},
		{"single integer", `(pr-str 42)`, "42"},
		{"integers joined with space", `(pr-str 1 2 3)`, "1 2 3"},
		{"single string is readable", `(pr-str "hi")`, `"hi"`},
		{"strings joined with space and readable", `(pr-str "a" "b")`, `"a" "b"`},
		{"string and integer", `(pr-str "a" 1)`, `"a" 1`},
		{"integer and string", `(pr-str 1 "a")`, `1 "a"`},
		{"mixed", `(pr-str "a" 1 "b" 2)`, `"a" 1 "b" 2`},
		{"empty string is quoted", `(pr-str "")`, `""`},
		{"negative integer", `(pr-str -7)`, "-7"},
	}, assertString)
}

func TestPrn(t *testing.T) {
	runTable(t, []tc[string]{
		{"no args prints newline", `(prn)`, "\n"},
		{"single integer", `(prn 1)`, "1\n"},
		{"integers joined with space", `(prn 1 2 3)`, "1 2 3\n"},
		{"single string is readable", `(prn "hi")`, "\"hi\"\n"},
		{"strings are readable", `(prn "a" "b")`, "\"a\" \"b\"\n"},
		{"string and integer", `(prn "a" 1)`, "\"a\" 1\n"},
		{"integer and string", `(prn 1 "a")`, "1 \"a\"\n"},
	}, assertPrint)
}

func TestPrintln(t *testing.T) {
	runTable(t, []tc[string]{
		{"no args prints newline", `(println)`, "\n"},
		{"single integer", `(println 1)`, "1\n"},
		{"integers joined with space", `(println 1 2 3)`, "1 2 3\n"},
		{"single string is raw", `(println "hi")`, "hi\n"},
		{"strings are raw and joined with space", `(println "a" "b")`, "a b\n"},
		{"string and integer raw", `(println "a" 1)`, "a 1\n"},
		{"integer and string raw", `(println 1 "a")`, "1 a\n"},
		{"sentence", `(println "hello" "world")`, "hello world\n"},
	}, assertPrint)
}

func isFatalValue(v core.Value) bool {
	if _, ok := v.(typ.Fatal); ok {
		return true
	}
	if _, ok := v.(*typ.Fatal); ok {
		return true
	}
	return false
}

func assertBool(t *testing.T, input string, want bool) {
	t.Helper()
	got := tut.MustEval(envs.MalDefault, input)
	b, ok := got.(typ.Bool)
	if !ok {
		t.Fatalf("MustEval(%q) type = %T (%v), want %T", input, got, got, typ.Bool(false))
	}
	if bool(b) != want {
		t.Fatalf("MustEval(%q) = %v, want %v", input, bool(b), want)
	}
}

func assertFatal(t *testing.T, input string) {
	t.Helper()
	got := tut.MustEval(envs.MalFun, input)
	if !isFatalValue(got) {
		t.Fatalf("MustEval(%q) = %v (%T), want Fatal error", input, got, got)
	}
}

func assertCountValue(t *testing.T, input string, want int) {
	t.Helper()
	got := tut.MustEval(envs.MalFun, input)
	if isFatalValue(got) {
		t.Fatalf("MustEval(%q) = Fatal %v, want %d", input, got, want)
	}
	if fmt.Sprintf("%v", got) != fmt.Sprintf("%d", want) {
		t.Fatalf("MustEval(%q) = %v (%T), want %d", input, got, got, want)
	}
}

func TestEqualOp(t *testing.T) {
	runTable(t, []tc[bool]{
		{"int same value", `(= 1 1)`, true},
		{"int different value", `(= 1 2)`, false},
		{"negative int same", `(= -7 -7)`, true},
		{"string same value", `(= "foo" "foo")`, true},
		{"string different value", `(= "foo" "bar")`, false},
		{"empty string equal", `(= "" "")`, true},
		{"bool true equal", `(= true true)`, true},
		{"bool false equal", `(= false false)`, true},
		{"bool true vs false", `(= true false)`, false},
		{"nil equal", `(= nil nil)`, true},
		{"empty list equal", `(= (list) (list))`, true},
		{"list same values", `(= (list 1 2) (list 1 2))`, true},
		{"list different values", `(= (list 1 2) (list 1 3))`, false},
		{"list different length", `(= (list 1) (list 1 2))`, false},
		{"empty vector equal", `(= [] [])`, true},
		{"vector same values", `(= [1 2] [1 2])`, true},
		{"vector different values", `(= [1 2] [1 3])`, false},
		{"empty map equal", `(= {} {})`, true},
		{"map same values", `(= {"a" 1} {"a" 1})`, true},
		{"map different values", `(= {"a" 1} {"a" 2})`, false},
		{"cross int vs string", `(= 1 "1")`, false},
		{"cross int vs bool", `(= 1 true)`, false},
		{"cross int vs nil", `(= 0 nil)`, false},
		{"cross string vs bool", `(= "true" true)`, false},
		{"cross string vs nil", `(= "" nil)`, false},
		{"cross bool vs nil true", `(= true nil)`, false},
		{"cross bool vs nil false", `(= false nil)`, false},
		{"cross list vs vector", `(= (list 1 2) [1 2])`, true},
		{"cross list vs nil", `(= (list 1) nil)`, false},
		{"cross vector vs map", `(= [1] {"a" 1})`, false},
		{"cross vector vs nil", `(= [] nil)`, false},
		{"cross map vs nil", `(= {} nil)`, false},
		{"cross int vs list", `(= 1 (list 1))`, false},
		{"cross string vs vector", `(= "a" ["a"])`, false},
	}, assertBoolCheck)

	t.Run("arity", func(t *testing.T) {
		//nolint:dupword
		assertFatalAll(t, []string{`(=)`, `(= 1)`, `(= 1 2 3)`, `(= nil nil nil)`})
	})
}

func TestCount(t *testing.T) {
	runTable(t, []tc[int]{
		{"nil as zero", `(count nil)`, 0},
		{"empty list", `(count (list))`, 0},
		{"multi-element list", `(count (list 1 2 3))`, 3},
		{"empty vector", `(count [])`, 0},
		{"non-empty vector", `(count [1 2])`, 2},
		{"empty map", `(count {})`, 0},
		{"non-empty map", `(count {"a" 1 "b" 2})`, 2},
	}, assertCountCheck)

	t.Run("non-containers return Fatal", func(t *testing.T) {
		assertFatalAll(t, []string{`(count 1)`, `(count 0)`, `(count true)`, `(count false)`})
	})

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{`(count)`, `(count [1] [2])`, `(count nil nil)`})
	})
}

func TestEmptyQ(t *testing.T) {
	runTable(t, []tc[bool]{
		{"nil is empty", `(empty? nil)`, true},
		{"empty list", `(empty? (list))`, true},
		{"non-empty list", `(empty? (list 1))`, false},
		{"multi-element list", `(empty? (list 1 2 3))`, false},
		{"empty vector", `(empty? [])`, true},
		{"non-empty vector", `(empty? [1])`, false},
		{"empty map", `(empty? {})`, true},
		{"non-empty map", `(empty? {"a" 1})`, false},
	}, assertBoolCheck)

	t.Run("non-containers return Fatal", func(t *testing.T) {
		assertFatalAll(t, []string{
			`(empty? 1)`, `(empty? 0)`, `(empty? true)`, `(empty? false)`,
		})
	})

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{`(empty?)`, `(empty? [] [])`, `(empty? nil nil)`})
	})
}

func TestListConstruction(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		assertListLen(t, `(list)`, 0)
		assertBool(t, `(= (list) (list))`, true)
	})

	t.Run("multi-element list", func(t *testing.T) {
		lst := assertListLen(t, `(list 1 2 3)`, 3)
		for i, src := range []string{`1`, `2`, `3`} {
			want := evalFun(src)
			if !lst[i].Equal(want) {
				t.Fatalf("element %d = %v, want %v", i, lst[i], want)
			}
		}
		assertBool(t, `(= (list 1 2 3) (list 1 2 3))`, true)
		assertBool(t, `(= (list 1 2) (list 1 2 3))`, false)
	})
}

func TestIsList(t *testing.T) {
	runTable(t, []tc[bool]{
		{"empty list true", `(list? (list))`, true},
		{"multi-element list true", `(list? (list 1 2 3))`, true},
		{"vector false", `(list? [1 2])`, false},
		{"empty vector false", `(list? [])`, false},
		{"nil false", `(list? nil)`, false},
		{"scalar int false", `(list? 1)`, false},
		{"scalar string false", `(list? "hi")`, false},
		{"scalar bool false", `(list? true)`, false},
		{"map false", `(list? {})`, false},
	}, assertBoolCheck)

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{`(list?)`, `(list? (list) (list))`})
	})
}

func TestReadString(t *testing.T) {
	t.Run("one form only", func(t *testing.T) {
		runTable(t, []tc[string]{
			{"integer", `(read-string "1")`, `1`},
			{"string", `(read-string "\"hi\"")`, `"hi"`},
			{"list form", `(read-string "(1 2)")`, `(list 1 2)`},
			{"vector form", `(read-string "[1 2]")`, `[1 2]`},
		}, equalWith(evalFun))
	})

	t.Run("multi-form empty or invalid return Fatal", func(t *testing.T) {
		assertFatalAll(t, []string{
			`(read-string "1 2")`,
			`(read-string "")`,
			`(read-string "   ")`,
			`(read-string "(1 2")`,
			`(read-string "[")`,
		})
	})

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{`(read-string)`, `(read-string "1" "2")`})
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
		got := tut.MustEval(envs.MalFun, input)
		s, ok := got.(typ.String)
		if !ok {
			t.Fatalf("MustEval(%q) type = %T, want %T", input, got, typ.String(""))
		}
		if string(s) != content {
			t.Fatalf("MustEval(%q) = %q, want %q", input, string(s), content)
		}
	})

	t.Run("missing file returns Fatal without panic", func(t *testing.T) {
		dir := t.TempDir()
		missing := filepath.Join(dir, "does-not-exist.txt")
		input := fmt.Sprintf(`(slurp %q)`, missing)
		got := tut.MustEval(envs.MalFun, input)
		if !isFatalValue(got) {
			t.Fatalf("MustEval(%q) = %v (%T), want Fatal error", input, got, got)
		}
	})

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{`(slurp)`, fmt.Sprintf(`(slurp %q %q)`, "a", "b")})
	})
}

func TestAtomCreation(t *testing.T) {
	t.Run("creates atom holding value", func(t *testing.T) {
		runTable(t, []tc[string]{
			{"int", `(atom 1)`, ""},
			{"string", `(atom "hi")`, ""},
			{"nil", `(atom nil)`, ""},
			{"list", `(atom (list 1 2))`, ""},
		}, assertAtom)
	})

	t.Run("deref returns initial value", func(t *testing.T) {
		runTable(t, []tc[string]{
			{"int", `(deref (atom 42))`, `42`},
			{"string", `(deref (atom "hi"))`, `"hi"`},
			{"nil", `(deref (atom nil))`, `nil`},
			{"list", `(deref (atom (list 1 2)))`, `(list 1 2)`},
			{"vector", `(deref (atom [1 2]))`, `[1 2]`},
			{"nested atom", `(atom? (deref (atom (atom 1))))`, `true`},
		}, equalWith(evalDefault))
	})

	t.Run("is atom", func(t *testing.T) {
		assertBool(t, `(atom? (atom 1))`, true)
		assertBool(t, `(atom? (atom nil))`, true)
	})

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{`(atom)`, `(atom 1 2)`, `(atom nil nil)`})
	})
}

func TestAtomQ(t *testing.T) {
	runTable(t, []tc[bool]{
		{"atom of int", `(atom? (atom 1))`, true},
		{"atom of string", `(atom? (atom "hi"))`, true},
		{"atom of nil", `(atom? (atom nil))`, true},
		{"atom of list", `(atom? (atom (list 1)))`, true},
		{"atom of vector", `(atom? (atom [1 2]))`, true},
		{"atom of atom", `(atom? (atom (atom 1)))`, true},
		{"int false", `(atom? 1)`, false},
		{"string false", `(atom? "hi")`, false},
		{"nil false", `(atom? nil)`, false},
		{"bool false", `(atom? true)`, false},
		{"list false", `(atom? (list 1 2))`, false},
		{"empty list false", `(atom? (list))`, false},
		{"vector false", `(atom? [1 2])`, false},
		{"empty vector false", `(atom? [])`, false},
		{"map false", `(atom? {})`, false},
		{"deref value false", `(atom? (deref (atom 1)))`, false},
	}, assertBoolCheck)

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{`(atom?)`, `(atom? (atom 1) (atom 2))`})
	})
}

func TestDeref(t *testing.T) {
	t.Run("returns held value", func(t *testing.T) {
		runTable(t, []tc[string]{
			{"int", `(deref (atom 1))`, `1`},
			{"negative int", `(deref (atom -7))`, `-7`},
			{"string", `(deref (atom "foo"))`, `"foo"`},
			{"empty string", `(deref (atom ""))`, `""`},
			{"bool", `(deref (atom true))`, `true`},
			{"nil", `(deref (atom nil))`, `nil`},
			{"list", `(deref (atom (list 1 2)))`, `(list 1 2)`},
			{"vector", `(deref (atom [1 2]))`, `[1 2]`},
			{"map", `(deref (atom {"a" 1}))`, `{"a" 1}`},
		}, equalWith(evalDefault))
	})

	t.Run("sees reset value via let", func(t *testing.T) {
		assertBool(t, `(= (let* [a (atom 1)] (do (reset! a 2) (deref a))) 2)`, true)
	})

	t.Run("sees swap value via let", func(t *testing.T) {
		assertBool(t, `(= (let* [a (atom 1)] (do (swap! a list 2) (deref a))) (list 1 2))`, true)
	})

	t.Run("non-atoms return Fatal", func(t *testing.T) {
		assertFatalAll(t, []string{
			`(deref 1)`, `(deref "hi")`, `(deref true)`, `(deref nil)`,
			`(deref (list 1))`, `(deref [1 2])`, `(deref {})`,
		})
	})

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{`(deref)`, `(deref (atom 1) (atom 2))`})
	})
}

func TestResetOp(t *testing.T) {
	t.Run("returns new value", func(t *testing.T) {
		runTable(t, []tc[string]{
			{"int to int", `(reset! (atom 1) 2)`, `2`},
			{"int to string", `(reset! (atom 1) "hi")`, `"hi"`},
			{"string to int", `(reset! (atom "hi") 42)`, `42`},
			{"to nil", `(reset! (atom 1) nil)`, `nil`},
			{"nil to int", `(reset! (atom nil) 1)`, `1`},
			{"to list", `(reset! (atom 1) (list 1 2))`, `(list 1 2)`},
			{"to vector", `(reset! (atom 1) [1 2])`, `[1 2]`},
		}, equalWith(evalFun))
	})

	t.Run("mutates atom", func(t *testing.T) {
		assertBool(t, `(= (let* [a (atom 1)] (do (reset! a 2) (deref a))) 2)`, true)
		assertBool(t, `(= (let* [a (atom 1)] (do (reset! a "hi") (deref a))) "hi")`, true)
		assertBool(t, `(= (let* [a (atom 1)] (do (reset! a nil) (deref a))) nil)`, true)
		assertBool(t, `(= (let* [a (atom 1)] (do (reset! a (list 1 2)) (deref a))) (list 1 2))`,
			true)
	})

	t.Run("non-atoms return Fatal", func(t *testing.T) {
		assertFatalAll(t, []string{
			`(reset! 1 2)`, `(reset! "hi" 2)`, `(reset! nil 2)`,
			`(reset! (list 1) 2)`, `(reset! [1 2] 2)`, `(reset! {} 2)`,
		})
	})

	t.Run("arity", func(t *testing.T) {
		assertFatalAll(t, []string{
			`(reset!)`, `(reset! (atom 1))`, `(reset! (atom 1) 2 3)`,
		})
	})
}

func TestSwapOp(t *testing.T) {
	t.Run("returns new value", func(t *testing.T) {
		runTable(t, []tc[string]{
			{"list no extra args", `(swap! (atom 1) list)`, `(list 1)`},
			{"list with extra args", `(swap! (atom 1) list 2 3)`, `(list 1 2 3)`},
			{"str no extra args", `(swap! (atom "a") str)`, `"a"`},
			{"str with extra args", `(swap! (atom "a") str "b" "c")`, `"abc"`},
			{"str int atom", `(swap! (atom 1) str "a")`, `"1a"`},
			{"pr-str", `(swap! (atom 1) pr-str 2)`, `"1 2"`},
		}, equalWith(evalFun))
	})

	t.Run("mutates atom", func(t *testing.T) {
		assertBool(t, `(= (let* [a (atom 1)] (do (swap! a list 2 3) (deref a))) (list 1 2 3))`,
			true)
		assertBool(t, `(= (let* [a (atom "a")] (do (swap! a str "b") (deref a))) "ab")`, true)
	})

	t.Run("too few args return Fatal", func(t *testing.T) {
		assertFatalAll(t, []string{`(swap!)`, `(swap! (atom 1))`})
	})

	t.Run("non-atom first arg returns Fatal", func(t *testing.T) {
		assertFatalAll(t, []string{
			`(swap! 1 list)`, `(swap! "hi" list)`, `(swap! nil list)`,
			`(swap! (list 1) list)`, `(swap! [1 2] list)`,
		})
	})

	t.Run("non-function second arg returns Fatal", func(t *testing.T) {
		assertFatalAll(t, []string{
			`(swap! (atom 1) 1)`, `(swap! (atom 1) "hi")`,
			`(swap! (atom 1) nil)`, `(swap! (atom 1) [1 2])`,
			`(swap! (atom 1) 1 2)`,
		})
	})
}
