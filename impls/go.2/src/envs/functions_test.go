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
	tests := []struct {
		name  string
		input string
		want  string
	}{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tut.MustEval(envs.MalFun, tt.input)
			s, ok := got.(typ.String)
			if !ok {
				t.Fatalf("MustEval(%q) type = %T, want %T", tt.input, got, typ.String(""))
			}
			if string(s) != tt.want {
				t.Fatalf("MustEval(%q) = %q, want %q", tt.input, string(s), tt.want)
			}
		})
	}
}

func TestPrStr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tut.MustEval(envs.MalFun, tt.input)
			s, ok := got.(typ.String)
			if !ok {
				t.Fatalf("MustEval(%q) type = %T, want %T", tt.input, got, typ.String(""))
			}
			if string(s) != tt.want {
				t.Fatalf("MustEval(%q) = %q, want %q", tt.input, string(s), tt.want)
			}
		})
	}
}

func TestPrn(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPrint string
	}{
		{"no args prints newline", `(prn)`, "\n"},
		{"single integer", `(prn 1)`, "1\n"},
		{"integers joined with space", `(prn 1 2 3)`, "1 2 3\n"},
		{"single string is readable", `(prn "hi")`, "\"hi\"\n"},
		{"strings are readable", `(prn "a" "b")`, "\"a\" \"b\"\n"},
		{"string and integer", `(prn "a" 1)`, "\"a\" 1\n"},
		{"integer and string", `(prn 1 "a")`, "1 \"a\"\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got core.Value
			output := captureStdout(t, func() {
				got = tut.MustEval(envs.MalFun, tt.input)
			})
			if _, ok := got.(typ.Nil); !ok {
				t.Fatalf("MustEval(%q) type = %T, want %T", tt.input, got, typ.Nil{})
			}
			if output != tt.wantPrint {
				t.Fatalf("MustEval(%q) printed %q, want %q", tt.input, output, tt.wantPrint)
			}
		})
	}
}

func TestPrintln(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPrint string
	}{
		{"no args prints newline", `(println)`, "\n"},
		{"single integer", `(println 1)`, "1\n"},
		{"integers joined with space", `(println 1 2 3)`, "1 2 3\n"},
		{"single string is raw", `(println "hi")`, "hi\n"},
		{"strings are raw and joined with space", `(println "a" "b")`, "a b\n"},
		{"string and integer raw", `(println "a" 1)`, "a 1\n"},
		{"integer and string raw", `(println 1 "a")`, "1 a\n"},
		{"sentence", `(println "hello" "world")`, "hello world\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got core.Value
			output := captureStdout(t, func() {
				got = tut.MustEval(envs.MalFun, tt.input)
			})
			if _, ok := got.(typ.Nil); !ok {
				t.Fatalf("MustEval(%q) type = %T, want %T", tt.input, got, typ.Nil{})
			}
			if output != tt.wantPrint {
				t.Fatalf("MustEval(%q) printed %q, want %q", tt.input, output, tt.wantPrint)
			}
		})
	}
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
	tests := []struct {
		name  string
		input string
		want  bool
	}{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBool(t, tt.input, tt.want)
		})
	}

	t.Run("arity", func(t *testing.T) {
		//nolint:dupword
		for _, input := range []string{`(=)`, `(= 1)`, `(= 1 2 3)`, `(= nil nil nil)`} {
			assertFatal(t, input)
		}
	})
}

func TestCount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"nil as zero", `(count nil)`, 0},
		{"empty list", `(count (list))`, 0},
		{"multi-element list", `(count (list 1 2 3))`, 3},
		{"empty vector", `(count [])`, 0},
		{"non-empty vector", `(count [1 2])`, 2},
		{"empty map", `(count {})`, 0},
		{"non-empty map", `(count {"a" 1 "b" 2})`, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCountValue(t, tt.input, tt.want)
		})
	}

	t.Run("non-containers return Fatal", func(t *testing.T) {
		for _, input := range []string{`(count 1)`, `(count 0)`, `(count true)`, `(count false)`} {
			assertFatal(t, input)
		}
	})

	t.Run("arity", func(t *testing.T) {
		for _, input := range []string{`(count)`, `(count [1] [2])`, `(count nil nil)`} {
			assertFatal(t, input)
		}
	})
}

func TestEmptyQ(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"nil is empty", `(empty? nil)`, true},
		{"empty list", `(empty? (list))`, true},
		{"non-empty list", `(empty? (list 1))`, false},
		{"multi-element list", `(empty? (list 1 2 3))`, false},
		{"empty vector", `(empty? [])`, true},
		{"non-empty vector", `(empty? [1])`, false},
		{"empty map", `(empty? {})`, true},
		{"non-empty map", `(empty? {"a" 1})`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBool(t, tt.input, tt.want)
		})
	}

	t.Run("non-containers return Fatal", func(t *testing.T) {
		for _, input := range []string{
			`(empty? 1)`, `(empty? 0)`, `(empty? true)`, `(empty? false)`,
		} {
			assertFatal(t, input)
		}
	})

	t.Run("arity", func(t *testing.T) {
		for _, input := range []string{`(empty?)`, `(empty? [] [])`, `(empty? nil nil)`} {
			assertFatal(t, input)
		}
	})
}

func TestListConstruction(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		got := tut.MustEval(envs.MalFun, `(list)`)
		lst, ok := got.(typ.List)
		if !ok {
			t.Fatalf("MustEval(%q) type = %T, want %T", `(list)`, got, typ.List{})
		}
		if len(lst) != 0 {
			t.Fatalf("MustEval(%q) length = %d, want 0", `(list)`, len(lst))
		}
		assertBool(t, `(= (list) (list))`, true)
	})

	t.Run("multi-element list", func(t *testing.T) {
		got := tut.MustEval(envs.MalFun, `(list 1 2 3)`)
		lst, ok := got.(typ.List)
		if !ok {
			t.Fatalf("MustEval(%q) type = %T, want %T", `(list 1 2 3)`, got, typ.List{})
		}
		if len(lst) != 3 {
			t.Fatalf("MustEval(%q) length = %d, want 3", `(list 1 2 3)`, len(lst))
		}
		for i, src := range []string{`1`, `2`, `3`} {
			want := tut.MustEval(envs.MalFun, src)
			if !lst[i].Equal(want) {
				t.Fatalf("element %d = %v, want %v", i, lst[i], want)
			}
		}
		assertBool(t, `(= (list 1 2 3) (list 1 2 3))`, true)
		assertBool(t, `(= (list 1 2) (list 1 2 3))`, false)
	})
}

func TestIsList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty list true", `(list? (list))`, true},
		{"multi-element list true", `(list? (list 1 2 3))`, true},
		{"vector false", `(list? [1 2])`, false},
		{"empty vector false", `(list? [])`, false},
		{"nil false", `(list? nil)`, false},
		{"scalar int false", `(list? 1)`, false},
		{"scalar string false", `(list? "hi")`, false},
		{"scalar bool false", `(list? true)`, false},
		{"map false", `(list? {})`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBool(t, tt.input, tt.want)
		})
	}

	t.Run("arity", func(t *testing.T) {
		for _, input := range []string{`(list?)`, `(list? (list) (list))`} {
			assertFatal(t, input)
		}
	})
}

func TestReadString(t *testing.T) {
	t.Run("one form only", func(t *testing.T) {
		cases := []struct {
			name string
			src  string
			want string
		}{
			{"integer", `(read-string "1")`, `1`},
			{"string", `(read-string "\"hi\"")`, `"hi"`},
			{"list form", `(read-string "(1 2)")`, `(list 1 2)`},
			{"vector form", `(read-string "[1 2]")`, `[1 2]`},
		}
		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				got := tut.MustEval(envs.MalFun, tt.src)
				if isFatalValue(got) {
					t.Fatalf("MustEval(%q) = Fatal %v, want %q", tt.src, got, tt.want)
				}
				want := tut.MustEval(envs.MalFun, tt.want)
				if !got.Equal(want) {
					t.Fatalf("MustEval(%q) = %v, want %v", tt.src, got, want)
				}
			})
		}
	})

	t.Run("multi-form empty or invalid return Fatal", func(t *testing.T) {
		for _, input := range []string{
			`(read-string "1 2")`,
			`(read-string "")`,
			`(read-string "   ")`,
			`(read-string "(1 2")`,
			`(read-string "[")`,
		} {
			assertFatal(t, input)
		}
	})

	t.Run("arity", func(t *testing.T) {
		for _, input := range []string{`(read-string)`, `(read-string "1" "2")`} {
			assertFatal(t, input)
		}
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
		for _, input := range []string{`(slurp)`, fmt.Sprintf(`(slurp %q %q)`, "a", "b")} {
			assertFatal(t, input)
		}
	})
}
