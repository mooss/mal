package envs_test

import (
	"testing"

	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
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
