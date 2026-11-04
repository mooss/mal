//revive:disable:cognitive-complexity,line-length-limit,function-length,cyclomatic
package reader

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lexer"
)

// Constructors keep tables concise and make new types easy to add.
func sym(s string) core.Symbol { return core.Symbol(s) }
func num(n int64) core.Integer { return core.Integer(n) }
func dummy() core.Dummy        { return core.Dummy{} }
func quoted(f core.Type) core.List {
	return core.List{core.Symbol("quote"), f}
}

func lst(elems ...core.Type) core.List {
	if len(elems) == 0 {
		return core.List{}
	}
	return core.List(elems)
}

// equalType compares Mal values structurally.
// Extend it when new core types are added (for example strings).
// Nil vs empty list counts as equal.
func equalType(a, b core.Type) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if a.Type() != b.Type() {
		return false
	}
	switch av := a.(type) {
	case core.Integer:
		bv, ok := b.(core.Integer)
		return ok && av == bv
	case core.Symbol:
		bv, ok := b.(core.Symbol)
		return ok && av == bv
	case core.List:
		bv, ok := b.(core.List)
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !equalType(av[i], bv[i]) {
				return false
			}
		}
		return true
	case core.Dummy:
		_, ok := b.(core.Dummy)
		return ok
	default:
		return false
	}
}

func equalForms(a, b []core.Type) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !equalType(a[i], b[i]) {
			return false
		}
	}
	return true
}

// singleCase tests a single ReadForm call (only the first form is checked,
// trailing tokens are ignored).
type singleCase struct {
	name    string
	input   string
	want    core.Type
	wantErr bool
}

func runSingleCases(t *testing.T, cases []singleCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := ReadStr(tc.input)
			got, err := r.ReadForm()
			gotErr := err != nil
			if gotErr != tc.wantErr {
				if tc.wantErr {
					t.Fatalf("input %q: expected error, got value %#v", tc.input, got)
				} else {
					t.Fatalf("input %q: expected no error, got error %v", tc.input, err)
				}
			}
			if !tc.wantErr && !equalType(got, tc.want) {
				t.Errorf("input %q: got %#v, want %#v", tc.input, got, tc.want)
			}
		})
	}
}

// multiCase tests reading until Done or error.
type multiCase struct {
	name    string
	input   string
	want    []core.Type
	wantErr bool
}

func runMultiCases(t *testing.T, cases []multiCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := ReadStr(tc.input)
			var got []core.Type
			gotErr := false
			for !r.Done() {
				form, err := r.ReadForm()
				if err != nil {
					gotErr = true
					break
				}
				got = append(got, form)
			}
			if gotErr != tc.wantErr {
				if tc.wantErr {
					t.Fatalf("input %q: expected error, got forms %#v", tc.input, got)
				} else {
					t.Fatalf("input %q: expected no error, got error while reading %#v", tc.input, got)
				}
			}
			if !equalForms(got, tc.want) {
				t.Errorf("input %q: got %#v, want %#v", tc.input, got, tc.want)
			}
		})
	}
}

func TestReadAtoms(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "zero", input: "0", want: num(0)},
		{name: "single digit", input: "1", want: num(1)},
		{name: "multi digit", input: "123", want: num(123)},
		{name: "leading zeros", input: "007", want: num(7)},
		{name: "negative zero", input: "-0", want: num(0)},
		{name: "negative single", input: "-1", want: num(-1)},
		{name: "negative multi", input: "-42", want: num(-42)},
		{name: "negative leading zeros", input: "-007", want: num(-7)},
		{name: "max int64", input: "9223372036854775807", want: num(9223372036854775807)},
		{name: "min int64", input: "-9223372036854775808", want: num(-9223372036854775808)},
		{name: "plus digits", input: "+123", want: num(123)},
		{name: "int overflow large", input: "999999999999999999999999", wantErr: true},
		{name: "int overflow max plus one", input: "9223372036854775808", wantErr: true},
		{name: "int overflow min minus one", input: "-9223372036854775809", wantErr: true},
		{name: "negative overflow large", input: "-999999999999999999999999", wantErr: true},
		{name: "single letter", input: "a", want: sym("a")},
		{name: "simple", input: "foo", want: sym("foo")},
		{name: "plus", input: "+", want: sym("+")},
		{name: "star", input: "*", want: sym("*")},
		{name: "single minus", input: "-", want: sym("-")},
		{name: "double minus", input: "--", want: sym("--")},
		{name: "minus letter", input: "-a", want: sym("-a")},
		{name: "dash", input: "foo-bar", want: sym("foo-bar")},
		{name: "slash", input: "foo/bar", want: sym("foo/bar")},
		{name: "dot", input: "foo.bar", want: sym("foo.bar")},
		{name: "nil", input: "nil", want: sym("nil")},
		{name: "true", input: "true", want: sym("true")},
		{name: "empty string", input: `""`, want: sym(`""`)},
		{name: "simple string", input: `"hello"`, want: sym(`"hello"`)},
		{name: "string with space", input: `"hello world"`, want: sym(`"hello world"`)},
		{name: "string with comma", input: `"a,b"`, want: sym(`"a,b"`)},
		{name: "escaped quote", input: `"a\"b"`, want: sym(`"a\"b"`)},
		{name: "padded symbol", input: "  foo  ", want: sym("foo")},
		{name: "padded int", input: " ,\t42\r\n, ", want: num(42)},
	})
}

func TestReadLists(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "empty", input: "()", want: lst()},
		{name: "empty padded", input: "( )", want: lst()},
		{name: "empty commas", input: "(,)", want: lst()},
		{name: "single symbol", input: "(foo)", want: lst(sym("foo"))},
		{name: "single int", input: "(1)", want: lst(num(1))},
		{name: "plus ints", input: "(+ 1 2)", want: lst(sym("+"), num(1), num(2))},
		{name: "negative ints", input: "(-1 -42)", want: lst(num(-1), num(-42))},
		{name: "mixed negative", input: "(foo -42 bar)", want: lst(sym("foo"), num(-42), sym("bar"))},
		{name: "def", input: "(def foo 42)", want: lst(sym("def"), sym("foo"), num(42))},
		{name: "nested", input: "((1 2) (3 4))", want: lst(lst(num(1), num(2)), lst(num(3), num(4)))},
		{name: "deep nested", input: "(1 (2 (3)))", want: lst(num(1), lst(num(2), lst(num(3))))},
		{name: "string in list", input: `(foo "bar baz")`, want: lst(sym("foo"), sym(`"bar baz"`))},
		{name: "commas", input: "(foo,123)", want: lst(sym("foo"), num(123))},
		{name: "commas padded", input: "(,foo,,123,,)", want: lst(sym("foo"), num(123))},
		{name: "no spaces", input: "(foo)", want: lst(sym("foo"))},
		{name: "multiline", input: "  ( \n foo\t 123 \r\n )  ", want: lst(sym("foo"), num(123))},
		{name: "quote inside", input: "('foo 1)", want: lst(quoted(sym("foo")), num(1))},
		{name: "dummy inside", input: "([1])", want: lst(dummy())},
		{name: "braces inside", input: "({a})", want: lst(dummy())},
		{name: "first form only", input: "(1) (2)", want: lst(num(1))},
	})
}

func TestReadQuotes(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "symbol", input: "'foo", want: quoted(sym("foo"))},
		{name: "spaced", input: "' foo", want: quoted(sym("foo"))},
		{name: "int", input: "'1", want: quoted(num(1))},
		{name: "string", input: `' "hello"`, want: quoted(sym(`"hello"`))},
		{name: "list", input: "'(1 2)", want: quoted(lst(num(1), num(2)))},
		{name: "empty list", input: "'()", want: quoted(lst())},
		{name: "double quote", input: "''foo", want: quoted(quoted(sym("foo")))},
		{name: "triple quote", input: "'''foo", want: quoted(quoted(quoted(sym("foo"))))},
		{name: "quoted list with quote", input: "'('foo)", want: quoted(lst(quoted(sym("foo"))))},
		{name: "quoted dummy", input: "'[1]", want: quoted(dummy())},
		{name: "quoted dummy braces", input: "'{a}", want: quoted(dummy())},
		{name: "lone quote", input: "'", wantErr: true},
		{name: "lone quote padded", input: "'   ", wantErr: true},
		{name: "quote closer", input: "' )", wantErr: true},
		{name: "quote closer bracket", input: "' ]", wantErr: true},
		{name: "quote unbalanced", input: "'(1", wantErr: true},
		{name: "quote unbalanced dummy", input: "'[1", wantErr: true},
	})
}

func TestReadDummies(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "empty brackets", input: "[]", want: dummy()},
		{name: "empty braces", input: "{}", want: dummy()},
		{name: "padded brackets", input: "[ ]", want: dummy()},
		{name: "padded braces", input: "{ }", want: dummy()},
		{name: "commas", input: "[,]", want: dummy()},
		{name: "ints", input: "[1 2 3]", want: dummy()},
		{name: "symbols", input: "{foo bar}", want: dummy()},
		{name: "list inside", input: "[(1 2)]", want: dummy()},
		{name: "nested brackets", input: "[[]]", want: dummy()},
		{name: "nested braces", input: "{{}}", want: dummy()},
		{name: "mixed nested", input: "[{}]", want: dummy()},
		{name: "quote inside", input: "['foo]", want: dummy()},
		{name: "string inside", input: `["hello" foo]`, want: dummy()},
		{name: "unclosed bracket", input: "[", wantErr: true},
		{name: "unclosed brace", input: "{", wantErr: true},
		{name: "unclosed with content", input: "[1 2", wantErr: true},
		{name: "mismatched", input: "[(1]", wantErr: true},
		{name: "paren in dummy", input: "[)", wantErr: true},
		{name: "lone rbracket", input: "]", wantErr: true},
		{name: "lone rbrace", input: "}", wantErr: true},
	})
}

func TestReadErrors(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "empty", input: "", wantErr: true},
		{name: "whitespace only", input: " \t\n\r ", wantErr: true},
		{name: "commas only", input: ",,,", wantErr: true},
		{name: "lone rparen", input: ")", wantErr: true},
		{name: "lone rbracket", input: "]", wantErr: true},
		{name: "lone rbrace", input: "}", wantErr: true},
		{name: "unclosed list", input: "(", wantErr: true},
		{name: "unclosed with symbol", input: "(foo", wantErr: true},
		{name: "unclosed nested", input: "(foo (bar)", wantErr: true},
		{name: "bracket in list", input: "(]", wantErr: true},
		{name: "brace in list", input: "(}", wantErr: true},
		{name: "paren in dummy", input: "{)", wantErr: true},
		{name: "backtick", input: "`", wantErr: true},
		{name: "int then letters", input: "123abc", wantErr: true},
		{name: "int then dash", input: "1-", wantErr: true},
		{name: "unterminated string", input: `"foo`, wantErr: true},
		{name: "escaped newline", input: "\"foo\\\nbar\"", wantErr: true},
		{name: "trailing backslash", input: "\"foo\\", wantErr: true},
		{name: "int overflow", input: "999999999999999999999999", wantErr: true},
		{name: "lone quote", input: "'", wantErr: true},
	})
}

func TestReadMultipleForms(t *testing.T) {
	runMultiCases(t, []multiCase{
		{name: "empty", input: "", want: nil, wantErr: false},
		{name: "whitespace only", input: " \t\n\r,, ", want: nil, wantErr: false},
		{name: "single int", input: "1", want: []core.Type{num(1)}, wantErr: false},
		{name: "two ints", input: "1 2", want: []core.Type{num(1), num(2)}, wantErr: false},
		{name: "mixed", input: `1 foo "hi"`, want: []core.Type{num(1), sym("foo"), sym(`"hi"`)}, wantErr: false},
		{name: "two lists", input: "(1) (2)", want: []core.Type{lst(num(1)), lst(num(2))}, wantErr: false},
		{name: "quote then list", input: "'a (1 2)", want: []core.Type{quoted(sym("a")), lst(num(1), num(2))}, wantErr: false},
		{name: "commas", input: "1, 2 , (3)", want: []core.Type{num(1), num(2), lst(num(3))}, wantErr: false},
		{name: "dummy then list", input: "[] (1)", want: []core.Type{dummy(), lst(num(1))}, wantErr: false},
		{name: "trailing rparen", input: "(1))", want: []core.Type{lst(num(1))}, wantErr: true},
		{name: "trailing rbracket", input: "(1) ]", want: []core.Type{lst(num(1))}, wantErr: true},
		{name: "second unbalanced", input: "(1) (2", want: []core.Type{lst(num(1))}, wantErr: true},
		{name: "token error", input: "1 `", want: nil, wantErr: true},
		{name: "unterminated", input: "(1) \"foo", want: nil, wantErr: true},
	})
}

func TestReaderCursor(t *testing.T) {
	t.Run("empty done rest", func(t *testing.T) {
		r := ReadStr("")
		if !r.Done() {
			t.Error("empty input: Done() = false, want true")
		}
		if r.Rest() != nil {
			t.Errorf("empty input: Rest() = %v, want nil", r.Rest())
		}
		if _, err := r.ReadForm(); err == nil {
			t.Error("empty input: ReadForm() expected error, got nil")
		}
	})

	t.Run("peek next rest done", func(t *testing.T) {
		r := ReadStr("1 foo")
		if r.Done() {
			t.Fatal("Done() = true before reading, want false")
		}
		if got := r.Peek(); got.Type != lexer.TOKEN_INT || got.Value != "1" {
			t.Fatalf("Peek() = %+v, want INT 1", got)
		}
		if got := len(r.Rest()); got != 2 {
			t.Fatalf("len(Rest()) = %d, want 2", got)
		}
		if got := r.Next(); got.Type != lexer.TOKEN_INT || got.Value != "1" {
			t.Fatalf("Next() = %+v, want INT 1", got)
		}
		if got := r.Peek(); got.Type != lexer.TOKEN_SYMBOL || got.Value != "foo" {
			t.Fatalf("Peek() after Next = %+v, want SYMBOL foo", got)
		}
		if got := r.Rest(); len(got) != 1 || got[0].Value != "foo" {
			t.Fatalf("Rest() after Next = %v, want [foo]", got)
		}
		if got := r.Next(); got.Value != "foo" {
			t.Fatalf("Next() second = %+v, want foo", got)
		}
		if !r.Done() {
			t.Error("Done() = false after consuming all, want true")
		}
		if r.Rest() != nil {
			t.Errorf("Rest() after consuming all = %v, want nil", r.Rest())
		}
	})

	t.Run("readform advances", func(t *testing.T) {
		r := ReadStr("(1) 2")
		first, err := r.ReadForm()
		if err != nil {
			t.Fatalf("first ReadForm() error: %v", err)
		}
		if !equalType(first, lst(num(1))) {
			t.Fatalf("first ReadForm() = %#v, want %#v", first, lst(num(1)))
		}
		if r.Done() {
			t.Fatal("Done() = true after first form, want false")
		}
		if got := r.Rest(); len(got) != 1 || got[0].Value != "2" {
			t.Fatalf("Rest() after first form = %v, want [2]", got)
		}
		second, err := r.ReadForm()
		if err != nil {
			t.Fatalf("second ReadForm() error: %v", err)
		}
		if !equalType(second, num(2)) {
			t.Fatalf("second ReadForm() = %#v, want 2", second)
		}
		if !r.Done() {
			t.Error("Done() = false after all forms, want true")
		}
	})
}
