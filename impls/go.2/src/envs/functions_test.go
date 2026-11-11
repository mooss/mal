package envs_test

import (
	"bytes"
	"io"
	"os"
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
