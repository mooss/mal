//revive:disable:cognitive-complexity,line-length-limit,function-length,cyclomatic
package reader

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lexer"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

// Constructors keep tables concise and make new types easy to add.
func sym(s string) typ.Symbol { return typ.Symbol(s) }
func str(s string) typ.String { return typ.String(s) }
func num(n int64) typ.Integer { return typ.Integer(n) }
func _nil() typ.Nil           { return typ.Nil{} }

func vec(elems ...core.Value) typ.VectorForm {
	return typ.VectorForm(elems)
}
func quoted(f core.Value) typ.List {
	return typ.List{typ.Symbol("quote"), f}
}

func lst(elems ...core.Value) typ.List {
	return typ.List(elems)
}

func mapl(pairs ...core.Value) typ.MapForm {
	form, err := typ.NewMapForm(pairs...)
	if err != nil {
		panic(err)
	}
	return form
}

func eqMapForm(left typ.MapForm, right core.Value) bool {
	rightVal, ok := right.(typ.MapForm)
	if !ok || len(left) != len(rightVal) {
		return false
	}
	for i := range left {
		if !equalType(left[i].Key, rightVal[i].Key) || !equalType(left[i].Value, rightVal[i].Value) {
			return false
		}
	}
	return true
}

func eqslc[T ~[]core.Value](leftVal T, right core.Value) bool {
	rightVal, ok := right.(T)
	if !ok || len(leftVal) != len(rightVal) {
		return false
	}
	for i := range leftVal {
		if !equalType(leftVal[i], rightVal[i]) {
			return false
		}
	}
	return true
}

func eqmap(leftVal typ.Map, right any) bool {
	rightVal, ok := right.(typ.Map)
	if !ok || len(leftVal) != len(rightVal) {
		return false
	}
	type pair struct {
		k core.Value
		v core.Value
	}
	bPairs := make([]pair, 0, len(rightVal))
	for k, v := range rightVal {
		bPairs = append(bPairs, pair{k, v})
	}
	used := make([]bool, len(bPairs))
	for ak, avv := range leftVal {
		found := false
		for i, p := range bPairs {
			if !used[i] && equalType(ak, p.k) && equalType(avv, p.v) {
				used[i] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func eqscal[T comparable](left T, right any) bool {
	rightVal, ok := right.(T)
	return ok && left == rightVal
}

// equalType compares Mal values structurally.
// Extend it when new core types are added (for example strings).
func equalType(left, right core.Value) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}

	switch leftVal := left.(type) {
	case typ.Integer:
		return eqscal(leftVal, right)
	case typ.Symbol:
		return eqscal(leftVal, right)
	case typ.String:
		return eqscal(leftVal, right)
	case typ.List:
		return eqslc(leftVal, right)
	case typ.VectorForm:
		return eqslc(leftVal, right)
	case typ.MapForm:
		return eqMapForm(leftVal, right)
	case typ.Map:
		return eqmap(leftVal, right)
	case typ.Nil:
		return ut.Is[typ.Nil](right)
	default:
		return false
	}
}

func equalForms(a, b []core.Value) bool {
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
	want    core.Value
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
	want    []core.Value
	wantErr bool
}

func runMultiCases(t *testing.T, cases []multiCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := ReadStr(tc.input)
			var got []core.Value
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
		{name: "nil", input: "nil", want: _nil()},
		{name: "true", input: "true", want: sym("true")},
		{name: "empty string", input: `""`, want: str(``)},
		{name: "simple string", input: `"hello"`, want: str(`hello`)},
		{name: "string with space", input: `"hello world"`, want: str(`hello world`)},
		{name: "string with comma", input: `"a,b"`, want: str(`a,b`)},
		{name: "escaped quote", input: `"a\"b"`, want: str(`a"b`)},
		{name: "padded symbol", input: "  foo  ", want: sym("foo")},
		{name: "padded int", input: " ,\t42\r\n, ", want: num(42)},
	})
}

func TestReadNil(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "nil padded", input: "  nil  ", want: _nil()},
		{name: "nil in list", input: "(nil)", want: lst(_nil())},
		{name: "nil in list mixed", input: "(1 nil foo)", want: lst(num(1), _nil(), sym("foo"))},
		{name: "nil in vector", input: "[nil]", want: vec(_nil())},
		{name: "nil in vector mixed", input: `[nil 1 "hi"]`, want: vec(_nil(), num(1), str(`hi`))},
		{name: "quoted nil", input: "'nil", want: quoted(_nil())},
		{name: "nil quoted list", input: "'(nil)", want: quoted(lst(_nil()))},
		{name: "nil prefix is symbol", input: "nilx", want: sym("nilx")},
		{name: "nil suffix dash is symbol", input: "nil-", want: sym("nil-")},
		{name: "uppercase NIL is symbol", input: "NIL", want: sym("NIL")},
		{name: "nil vs symbol not equal", input: "nil", want: _nil()},
	})
	runMultiCases(t, []multiCase{
		//nolint:dupword
		{name: "two nils", input: "nil nil", want: []core.Value{_nil(), _nil()}, wantErr: false},
		{name: "nil mixed", input: `nil 1 foo "hi"`, want: []core.Value{_nil(), num(1), sym("foo"), str(`hi`)}, wantErr: false},
		{name: "nil in forms", input: "(nil) [nil] nil", want: []core.Value{lst(_nil()), vec(_nil()), _nil()}, wantErr: false},
	})
	t.Run("nil not equal other types", func(t *testing.T) {
		if equalType(_nil(), sym("nil")) {
			t.Error("nil should not equal symbol nil")
		}
		if equalType(_nil(), lst()) {
			t.Error("nil should not equal empty list")
		}
		if equalType(_nil(), vec()) {
			t.Error("nil should not equal empty vector")
		}
		if equalType(_nil(), num(0)) {
			t.Error("nil should not equal zero")
		}
		if !equalType(_nil(), _nil()) {
			t.Error("nil should equal nil")
		}
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
		{name: "string in list", input: `(foo "bar baz")`, want: lst(sym("foo"), str(`bar baz`))},
		{name: "commas", input: "(foo,123)", want: lst(sym("foo"), num(123))},
		{name: "commas padded", input: "(,foo,,123,,)", want: lst(sym("foo"), num(123))},
		{name: "no spaces", input: "(foo)", want: lst(sym("foo"))},
		{name: "multiline", input: "  ( \n foo\t 123 \r\n )  ", want: lst(sym("foo"), num(123))},
		{name: "quote inside", input: "('foo 1)", want: lst(quoted(sym("foo")), num(1))},
		{name: "vector inside", input: "([1])", want: lst(vec(num(1)))},
		{name: "first form only", input: "(1) (2)", want: lst(num(1))},
	})
}

func TestReadQuotes(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "symbol", input: "'foo", want: quoted(sym("foo"))},
		{name: "spaced", input: "' foo", want: quoted(sym("foo"))},
		{name: "int", input: "'1", want: quoted(num(1))},
		{name: "string", input: `' "hello"`, want: quoted(str(`hello`))},
		{name: "list", input: "'(1 2)", want: quoted(lst(num(1), num(2)))},
		{name: "empty list", input: "'()", want: quoted(lst())},
		{name: "double quote", input: "''foo", want: quoted(quoted(sym("foo")))},
		{name: "triple quote", input: "'''foo", want: quoted(quoted(quoted(sym("foo"))))},
		{name: "quoted list with quote", input: "'('foo)", want: quoted(lst(quoted(sym("foo"))))},
		{name: "quoted vector", input: "'[1]", want: quoted(vec(num(1)))},
		{name: "lone quote", input: "'", wantErr: true},
		{name: "lone quote padded", input: "'   ", wantErr: true},
		{name: "quote closer", input: "' )", wantErr: true},
		{name: "quote closer bracket", input: "' ]", wantErr: true},
		{name: "quote unbalanced", input: "'(1", wantErr: true},
	})
}

func TestReadVectors(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "empty brackets", input: "[]", want: vec()},
		{name: "padded brackets", input: "[ ]", want: vec()},
		{name: "commas", input: "[,]", want: vec()},
		{name: "ints", input: "[1 2 3]", want: vec(num(1), num(2), num(3))},
		{name: "list inside", input: "[(1 2)]", want: vec(lst(num(1), num(2)))},
		{name: "nested brackets", input: "[[]]", want: vec(vec())},
		{name: "quote inside", input: "['foo]", want: vec(quoted(sym("foo")))},
		{name: "string inside", input: `["hello" foo]`, want: vec(str(`hello`), sym("foo"))},
		{name: "unclosed bracket", input: "[", wantErr: true},
		{name: "unclosed with content", input: "[1 2", wantErr: true},
		{name: "mismatched", input: "[(1]", wantErr: true},
		{name: "lone rbracket", input: "]", wantErr: true},
	})
}

func TestReadMaps(t *testing.T) {
	runSingleCases(t, []singleCase{
		{name: "empty", input: "{}", want: mapl()},
		{name: "padded", input: "{ }", want: mapl()},
		{name: "commas", input: "{,}", want: mapl()},
		{name: "commas padded", input: "{,,,}", want: mapl()},
		{name: "single string key", input: `{"a" 1}`, want: mapl(str(`a`), num(1))},
		{name: "doc example", input: `{"a" a, "b" 3}`, want: mapl(str(`a`), sym("a"), str(`b`), num(3))},
		{name: "symbol keys", input: "{foo bar, baz qux}", want: mapl(sym("foo"), sym("bar"), sym("baz"), sym("qux"))},
		{name: "int keys", input: "{1 2, 3 4}", want: mapl(num(1), num(2), num(3), num(4))},
		{name: "no spaces", input: `{"a"1}`, want: mapl(str(`a`), num(1))},
		{name: "commas between pairs", input: `{,"a",1,,"b",2,,}`, want: mapl(str(`a`), num(1), str(`b`), num(2))},
		{name: "multiline", input: "{\n\"a\" 1,\n\"b\" 2\n}", want: mapl(str(`a`), num(1), str(`b`), num(2))},
		{name: "nested values", input: `{"a" [1 2], "b" (3 4)}`, want: mapl(str(`a`), vec(num(1), num(2)), str(`b`), lst(num(3), num(4)))},
		{name: "nested map", input: `{"a" {"b" 2}}`, want: mapl(str(`a`), mapl(str(`b`), num(2)))},
		{name: "map in list", input: `({"a" 1} {"b" 2})`, want: lst(mapl(str(`a`), num(1)), mapl(str(`b`), num(2)))},
		{name: "map in vector", input: `[{"a" 1} {"b" 2}]`, want: vec(mapl(str(`a`), num(1)), mapl(str(`b`), num(2)))},
		{name: "list with map", input: `(foo {"a" 1})`, want: lst(sym("foo"), mapl(str(`a`), num(1)))},
		{name: "vector with map", input: `[foo {"a" 1}]`, want: vec(sym("foo"), mapl(str(`a`), num(1)))},
		{name: "quote inside value", input: `{a 'foo}`, want: mapl(sym("a"), quoted(sym("foo")))},
		{name: "quoted map", input: `'{"a" 1}`, want: quoted(mapl(str(`a`), num(1)))},
		{name: "quoted key", input: `{'foo 1}`, want: mapl(quoted(sym("foo")), num(1))},
		{name: "nil key and value", input: `{nil 1, "a" nil}`, want: mapl(_nil(), num(1), str(`a`), _nil())},
		{name: "nil in map in forms", input: `{nil nil}`, want: mapl(_nil(), _nil())},
		{name: "first form only", input: `{"a" 1} {"b" 2}`, want: mapl(str(`a`), num(1))},
		{name: "unclosed", input: "{", wantErr: true},
		{name: "unclosed with content", input: `{"a" 1`, wantErr: true},
		{name: "unclosed nested", input: `{"a" {"b" 2}`, wantErr: true},
		{name: "mismatched paren", input: "{(1}", wantErr: true},
		{name: "mismatched bracket", input: "{[1}", wantErr: true},
		{name: "close paren", input: `{"a" 1)`, wantErr: true},
		{name: "close bracket", input: `{"a" 1]`, wantErr: true},
		{name: "lone rbrace", input: "}", wantErr: true},
		{name: "odd single", input: "{1}", wantErr: true},
		{name: "odd triple", input: "{1 2 3}", wantErr: true},
		{name: "odd single string", input: `{"a"}`, wantErr: true},
		{name: "odd five", input: `{"a" 1 "b" 2 "c"}`, wantErr: true},
		{name: "unterminated string key", input: `{"a"`, wantErr: true},
	})
	runMultiCases(t, []multiCase{
		{name: "two maps", input: `{"a" 1} {"b" 2}`, want: []core.Value{mapl(str(`a`), num(1)), mapl(str(`b`), num(2))}, wantErr: false},
		{name: "map then list", input: `{"a" 1} (1)`, want: []core.Value{mapl(str(`a`), num(1)), lst(num(1))}, wantErr: false},
		{name: "second unbalanced", input: `{"a" 1} {`, want: []core.Value{mapl(str(`a`), num(1))}, wantErr: true},
		{name: "second odd", input: `{} {1}`, want: []core.Value{mapl()}, wantErr: true},
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
		{name: "single int", input: "1", want: []core.Value{num(1)}, wantErr: false},
		{name: "two ints", input: "1 2", want: []core.Value{num(1), num(2)}, wantErr: false},
		{name: "mixed", input: `1 foo "hi"`, want: []core.Value{num(1), sym("foo"), str(`hi`)}, wantErr: false},
		{name: "two lists", input: "(1) (2)", want: []core.Value{lst(num(1)), lst(num(2))}, wantErr: false},
		{name: "quote then list", input: "'a (1 2)", want: []core.Value{quoted(sym("a")), lst(num(1), num(2))}, wantErr: false},
		{name: "commas", input: "1, 2 , (3)", want: []core.Value{num(1), num(2), lst(num(3))}, wantErr: false},
		{name: "vector then list", input: "[] (1)", want: []core.Value{vec(), lst(num(1))}, wantErr: false},
		{name: "trailing rparen", input: "(1))", want: []core.Value{lst(num(1))}, wantErr: true},
		{name: "trailing rbracket", input: "(1) ]", want: []core.Value{lst(num(1))}, wantErr: true},
		{name: "second unbalanced", input: "(1) (2", want: []core.Value{lst(num(1))}, wantErr: true},
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
