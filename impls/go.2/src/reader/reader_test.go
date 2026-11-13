//revive:disable:cognitive-complexity,line-length-limit,function-length,cyclomatic
package reader_test

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lexer"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
)

var (
	sym   = tut.Sym
	str   = tut.Str
	num   = tut.Num
	lst   = tut.Lst
	vec   = tut.Vec
	mapf  = tut.MapForm
	qot   = tut.Quote
	nul   = tut.Nil
	valst = func(args ...core.Value) []core.Value { return args }
)

func runSingleCases(t *testing.T, cases []tut.Case[core.Value]) {
	t.Helper()
	tut.Run(t, cases, func(t *testing.T, c tut.Case[core.Value]) {
		t.Helper()
		r := reader.ReadStr(c.Input)
		got, err := r.ReadForm()
		if c.WantErr {
			assert.Err(t, err)
			return
		}
		assert.NoErr(t, err)
		assert.MalEqual(t, got, c.Want)
	})
}

func runMultiCases(t *testing.T, cases []tut.Case[[]core.Value]) {
	t.Helper()
	tut.Run(t, cases, func(t *testing.T, c tut.Case[[]core.Value]) {
		t.Helper()
		r := reader.ReadStr(c.Input)
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
		assert.Equal(t, gotErr, c.WantErr)
		assert.MalEqualSlice(t, got, c.Want)
	})
}

func TestReadAtoms(t *testing.T) {
	runSingleCases(t, []tut.Case[core.Value]{
		{Name: "zero", Input: "0", Want: num(0)},
		{Name: "single digit", Input: "1", Want: num(1)},
		{Name: "multi digit", Input: "123", Want: num(123)},
		{Name: "leading zeros", Input: "007", Want: num(7)},
		{Name: "negative zero", Input: "-0", Want: num(0)},
		{Name: "negative single", Input: "-1", Want: num(-1)},
		{Name: "negative multi", Input: "-42", Want: num(-42)},
		{Name: "negative leading zeros", Input: "-007", Want: num(-7)},
		{Name: "max int64", Input: "9223372036854775807", Want: num(9223372036854775807)},
		{Name: "min int64", Input: "-9223372036854775808", Want: num(-9223372036854775808)},
		{Name: "plus digits", Input: "+123", Want: num(123)},
		{Name: "int overflow large", Input: "999999999999999999999999", WantErr: true},
		{Name: "int overflow max plus one", Input: "9223372036854775808", WantErr: true},
		{Name: "int overflow min minus one", Input: "-9223372036854775809", WantErr: true},
		{Name: "negative overflow large", Input: "-999999999999999999999999", WantErr: true},
		{Name: "single letter", Input: "a", Want: sym("a")},
		{Name: "simple", Input: "foo", Want: sym("foo")},
		{Name: "plus", Input: "+", Want: sym("+")},
		{Name: "star", Input: "*", Want: sym("*")},
		{Name: "single minus", Input: "-", Want: sym("-")},
		{Name: "double minus", Input: "--", Want: sym("--")},
		{Name: "minus letter", Input: "-a", Want: sym("-a")},
		{Name: "dash", Input: "foo-bar", Want: sym("foo-bar")},
		{Name: "slash", Input: "foo/bar", Want: sym("foo/bar")},
		{Name: "dot", Input: "foo.bar", Want: sym("foo.bar")},
		{Name: "nil", Input: "nil", Want: nul()},
		{Name: "true", Input: "true", Want: sym("true")},
		{Name: "empty string", Input: `""`, Want: str(``)},
		{Name: "simple string", Input: `"hello"`, Want: str(`hello`)},
		{Name: "string with space", Input: `"hello world"`, Want: str(`hello world`)},
		{Name: "string with comma", Input: `"a,b"`, Want: str(`a,b`)},
		{Name: "escaped quote", Input: `"a\"b"`, Want: str(`a"b`)},
		{Name: "padded symbol", Input: "  foo  ", Want: sym("foo")},
		{Name: "padded int", Input: " ,\t42\r\n, ", Want: num(42)},
	})
}

func TestReadNil(t *testing.T) {
	runSingleCases(t, []tut.Case[core.Value]{
		{Name: "nil padded", Input: "  nil  ", Want: nul()},
		{Name: "nil in list", Input: "(nil)", Want: lst(nul())},
		{Name: "nil in list mixed", Input: "(1 nil foo)", Want: lst(num(1), nul(), sym("foo"))},
		{Name: "nil in vector", Input: "[nil]", Want: vec(nul())},
		{Name: "nil in vector mixed", Input: `[nil 1 "hi"]`, Want: vec(nul(), num(1), str(`hi`))},
		{Name: "quoted nil", Input: "'nil", Want: qot(nul())},
		{Name: "nil quoted list", Input: "'(nil)", Want: qot(lst(nul()))},
		{Name: "nil prefix is symbol", Input: "nilx", Want: sym("nilx")},
		{Name: "nil suffix dash is symbol", Input: "nil-", Want: sym("nil-")},
		{Name: "uppercase NIL is symbol", Input: "NIL", Want: sym("NIL")},
		{Name: "nil vs symbol not equal", Input: "nil", Want: nul()},
	})
	runMultiCases(t, []tut.Case[[]core.Value]{
		//nolint:dupword
		{Name: "two nils", Input: "nil nil", Want: valst(nul(), nul())},
		{Name: "nil mixed", Input: `nil 1 foo "hi"`, Want: valst(nul(), num(1), sym("foo"), str(`hi`))},
		{Name: "nil in forms", Input: "(nil) [nil] nil", Want: valst(lst(nul()), vec(nul()), nul())},
	})
	t.Run("nil not equal other types", func(t *testing.T) {
		assert.MalDifferent(t, nul(), sym("nil"))
		assert.MalDifferent(t, nul(), lst())
		assert.MalDifferent(t, nul(), vec())
		assert.MalDifferent(t, nul(), num(0))
		assert.MalEqual(t, nul(), nul())
	})
}

func TestReadLists(t *testing.T) {
	runSingleCases(t, []tut.Case[core.Value]{
		{Name: "empty", Input: "()", Want: lst()},
		{Name: "empty padded", Input: "( )", Want: lst()},
		{Name: "empty commas", Input: "(,)", Want: lst()},
		{Name: "single symbol", Input: "(foo)", Want: lst(sym("foo"))},
		{Name: "single int", Input: "(1)", Want: lst(num(1))},
		{Name: "plus ints", Input: "(+ 1 2)", Want: lst(sym("+"), num(1), num(2))},
		{Name: "negative ints", Input: "(-1 -42)", Want: lst(num(-1), num(-42))},
		{Name: "mixed negative", Input: "(foo -42 bar)", Want: lst(sym("foo"), num(-42), sym("bar"))},
		{Name: "def", Input: "(def foo 42)", Want: lst(sym("def"), sym("foo"), num(42))},
		{Name: "nested", Input: "((1 2) (3 4))", Want: lst(lst(num(1), num(2)), lst(num(3), num(4)))},
		{Name: "deep nested", Input: "(1 (2 (3)))", Want: lst(num(1), lst(num(2), lst(num(3))))},
		{Name: "string in list", Input: `(foo "bar baz")`, Want: lst(sym("foo"), str(`bar baz`))},
		{Name: "commas", Input: "(foo,123)", Want: lst(sym("foo"), num(123))},
		{Name: "commas padded", Input: "(,foo,,123,,)", Want: lst(sym("foo"), num(123))},
		{Name: "no spaces", Input: "(foo)", Want: lst(sym("foo"))},
		{Name: "multiline", Input: "  ( \n foo\t 123 \r\n )  ", Want: lst(sym("foo"), num(123))},
		{Name: "quote inside", Input: "('foo 1)", Want: lst(qot(sym("foo")), num(1))},
		{Name: "vector inside", Input: "([1])", Want: lst(vec(num(1)))},
		{Name: "first form only", Input: "(1) (2)", Want: lst(num(1))},
	})
}

func TestReadQuotes(t *testing.T) {
	runSingleCases(t, []tut.Case[core.Value]{
		{Name: "symbol", Input: "'foo", Want: qot(sym("foo"))},
		{Name: "spaced", Input: "' foo", Want: qot(sym("foo"))},
		{Name: "int", Input: "'1", Want: qot(num(1))},
		{Name: "string", Input: `' "hello"`, Want: qot(str(`hello`))},
		{Name: "list", Input: "'(1 2)", Want: qot(lst(num(1), num(2)))},
		{Name: "empty list", Input: "'()", Want: qot(lst())},
		{Name: "double quote", Input: "''foo", Want: qot(qot(sym("foo")))},
		{Name: "triple quote", Input: "'''foo", Want: qot(qot(qot(sym("foo"))))},
		{Name: "quoted list with quote", Input: "'('foo)", Want: qot(lst(qot(sym("foo"))))},
		{Name: "quoted vector", Input: "'[1]", Want: qot(vec(num(1)))},
		{Name: "lone quote", Input: "'", WantErr: true},
		{Name: "lone quote padded", Input: "'   ", WantErr: true},
		{Name: "quote closer", Input: "' )", WantErr: true},
		{Name: "quote closer bracket", Input: "' ]", WantErr: true},
		{Name: "quote unbalanced", Input: "'(1", WantErr: true},
	})
}

func TestReadVectors(t *testing.T) {
	runSingleCases(t, []tut.Case[core.Value]{
		{Name: "empty brackets", Input: "[]", Want: vec()},
		{Name: "padded brackets", Input: "[ ]", Want: vec()},
		{Name: "commas", Input: "[,]", Want: vec()},
		{Name: "ints", Input: "[1 2 3]", Want: vec(num(1), num(2), num(3))},
		{Name: "list inside", Input: "[(1 2)]", Want: vec(lst(num(1), num(2)))},
		{Name: "nested brackets", Input: "[[]]", Want: vec(vec())},
		{Name: "quote inside", Input: "['foo]", Want: vec(qot(sym("foo")))},
		{Name: "string inside", Input: `["hello" foo]`, Want: vec(str(`hello`), sym("foo"))},
		{Name: "unclosed bracket", Input: "[", WantErr: true},
		{Name: "unclosed with content", Input: "[1 2", WantErr: true},
		{Name: "mismatched", Input: "[(1]", WantErr: true},
		{Name: "lone rbracket", Input: "]", WantErr: true},
	})
}

func TestReadMaps(t *testing.T) {
	runSingleCases(t, []tut.Case[core.Value]{
		{Name: "empty", Input: "{}", Want: mapf()},
		{Name: "padded", Input: "{ }", Want: mapf()},
		{Name: "commas", Input: "{,}", Want: mapf()},
		{Name: "commas padded", Input: "{,,,}", Want: mapf()},
		{Name: "single string key", Input: `{"a" 1}`, Want: mapf(str(`a`), num(1))},
		{Name: "doc example", Input: `{"a" a, "b" 3}`, Want: mapf(str(`a`), sym("a"), str(`b`), num(3))},
		{Name: "symbol keys", Input: "{foo bar, baz qux}", Want: mapf(sym("foo"), sym("bar"), sym("baz"), sym("qux"))},
		{Name: "int keys", Input: "{1 2, 3 4}", Want: mapf(num(1), num(2), num(3), num(4))},
		{Name: "no spaces", Input: `{"a"1}`, Want: mapf(str(`a`), num(1))},
		{Name: "commas between pairs", Input: `{,"a",1,,"b",2,,}`, Want: mapf(str(`a`), num(1), str(`b`), num(2))},
		{Name: "multiline", Input: "{\n\"a\" 1,\n\"b\" 2\n}", Want: mapf(str(`a`), num(1), str(`b`), num(2))},
		{Name: "nested values", Input: `{"a" [1 2], "b" (3 4)}`, Want: mapf(str(`a`), vec(num(1), num(2)), str(`b`), lst(num(3), num(4)))},
		{Name: "nested map", Input: `{"a" {"b" 2}}`, Want: mapf(str(`a`), mapf(str(`b`), num(2)))},
		{Name: "map in list", Input: `({"a" 1} {"b" 2})`, Want: lst(mapf(str(`a`), num(1)), mapf(str(`b`), num(2)))},
		{Name: "map in vector", Input: `[{"a" 1} {"b" 2}]`, Want: vec(mapf(str(`a`), num(1)), mapf(str(`b`), num(2)))},
		{Name: "list with map", Input: `(foo {"a" 1})`, Want: lst(sym("foo"), mapf(str(`a`), num(1)))},
		{Name: "vector with map", Input: `[foo {"a" 1}]`, Want: vec(sym("foo"), mapf(str(`a`), num(1)))},
		{Name: "quote inside value", Input: `{a 'foo}`, Want: mapf(sym("a"), qot(sym("foo")))},
		{Name: "quoted map", Input: `'{"a" 1}`, Want: qot(mapf(str(`a`), num(1)))},
		{Name: "quoted key", Input: `{'foo 1}`, Want: mapf(qot(sym("foo")), num(1))},
		{Name: "nil key and value", Input: `{nil 1, "a" nil}`, Want: mapf(nul(), num(1), str(`a`), nul())},
		{Name: "nil in map in forms", Input: `{nil nil}`, Want: mapf(nul(), nul())},
		{Name: "first form only", Input: `{"a" 1} {"b" 2}`, Want: mapf(str(`a`), num(1))},
		{Name: "unclosed", Input: "{", WantErr: true},
		{Name: "unclosed with content", Input: `{"a" 1`, WantErr: true},
		{Name: "unclosed nested", Input: `{"a" {"b" 2}`, WantErr: true},
		{Name: "mismatched paren", Input: "{(1}", WantErr: true},
		{Name: "mismatched bracket", Input: "{[1}", WantErr: true},
		{Name: "close paren", Input: `{"a" 1)`, WantErr: true},
		{Name: "close bracket", Input: `{"a" 1]`, WantErr: true},
		{Name: "lone rbrace", Input: "}", WantErr: true},
		{Name: "odd single", Input: "{1}", WantErr: true},
		{Name: "odd triple", Input: "{1 2 3}", WantErr: true},
		{Name: "odd single string", Input: `{"a"}`, WantErr: true},
		{Name: "odd five", Input: `{"a" 1 "b" 2 "c"}`, WantErr: true},
		{Name: "unterminated string key", Input: `{"a"`, WantErr: true},
	})
	runMultiCases(t, []tut.Case[[]core.Value]{
		{Name: "two maps", Input: `{"a" 1} {"b" 2}`, Want: valst(mapf(str(`a`), num(1)), mapf(str(`b`), num(2)))},
		{Name: "map then list", Input: `{"a" 1} (1)`, Want: valst(mapf(str(`a`), num(1)), lst(num(1)))},
		{Name: "second unbalanced", Input: `{"a" 1} {`, Want: valst(mapf(str(`a`), num(1))), WantErr: true},
		{Name: "second odd", Input: `{} {1}`, Want: valst(mapf()), WantErr: true},
	})
}

func TestReadErrors(t *testing.T) {
	runSingleCases(t, []tut.Case[core.Value]{
		{Name: "empty", Input: "", WantErr: true},
		{Name: "whitespace only", Input: " \t\n\r ", WantErr: true},
		{Name: "commas only", Input: ",,,", WantErr: true},
		{Name: "lone rparen", Input: ")", WantErr: true},
		{Name: "lone rbracket", Input: "]", WantErr: true},
		{Name: "lone rbrace", Input: "}", WantErr: true},
		{Name: "unclosed list", Input: "(", WantErr: true},
		{Name: "unclosed with symbol", Input: "(foo", WantErr: true},
		{Name: "unclosed nested", Input: "(foo (bar)", WantErr: true},
		{Name: "bracket in list", Input: "(]", WantErr: true},
		{Name: "brace in list", Input: "(}", WantErr: true},
		{Name: "backtick", Input: "`", WantErr: true},
		{Name: "int then letters", Input: "123abc", WantErr: true},
		{Name: "int then dash", Input: "1-", WantErr: true},
		{Name: "unterminated string", Input: `"foo`, WantErr: true},
		{Name: "escaped newline", Input: "\"foo\\\nbar\"", WantErr: true},
		{Name: "trailing backslash", Input: "\"foo\\", WantErr: true},
		{Name: "int overflow", Input: "999999999999999999999999", WantErr: true},
		{Name: "lone quote", Input: "'", WantErr: true},
	})
}

func TestReadMultipleForms(t *testing.T) {
	runMultiCases(t, []tut.Case[[]core.Value]{
		{Name: "empty", Input: ""},
		{Name: "whitespace only", Input: " \t\n\r,, "},
		{Name: "single int", Input: "1", Want: valst(num(1))},
		{Name: "two ints", Input: "1 2", Want: valst(num(1), num(2))},
		{Name: "mixed", Input: `1 foo "hi"`, Want: valst(num(1), sym("foo"), str(`hi`))},
		{Name: "two lists", Input: "(1) (2)", Want: valst(lst(num(1)), lst(num(2)))},
		{Name: "quote then list", Input: "'a (1 2)", Want: valst(qot(sym("a")), lst(num(1), num(2)))},
		{Name: "commas", Input: "1, 2 , (3)", Want: valst(num(1), num(2), lst(num(3)))},
		{Name: "vector then list", Input: "[] (1)", Want: valst(vec(), lst(num(1)))},
		{Name: "trailing rparen", Input: "(1))", Want: valst(lst(num(1))), WantErr: true},
		{Name: "trailing rbracket", Input: "(1) ]", Want: valst(lst(num(1))), WantErr: true},
		{Name: "second unbalanced", Input: "(1) (2", Want: valst(lst(num(1))), WantErr: true},
		{Name: "token error", Input: "1 `", WantErr: true},
		{Name: "unterminated", Input: "(1) \"foo", WantErr: true},
	})
}

func TestReaderCursor(t *testing.T) {
	t.Run("empty done rest", func(t *testing.T) {
		r := reader.ReadStr("")
		assert.Equal(t, r.Done(), true)
		assert.Equal(t, len(r.Rest()), 0)
		_, err := r.ReadForm()
		assert.Err(t, err)
	})

	t.Run("peek next rest done", func(t *testing.T) {
		r := reader.ReadStr("1 foo")
		assert.Equal(t, r.Done(), false)
		peeked := r.Peek()
		assert.Equal(t, peeked.Type, lexer.TOKEN_INT)
		assert.Equal(t, peeked.Value, "1")
		assert.Equal(t, len(r.Rest()), 2)
		next := r.Next()
		assert.Equal(t, next.Type, lexer.TOKEN_INT)
		assert.Equal(t, next.Value, "1")
		peeked = r.Peek()
		assert.Equal(t, peeked.Type, lexer.TOKEN_SYMBOL)
		assert.Equal(t, peeked.Value, "foo")
		rest := r.Rest()
		assert.Equal(t, len(rest), 1)
		assert.Equal(t, rest[0].Value, "foo")
		next = r.Next()
		assert.Equal(t, next.Value, "foo")
		assert.Equal(t, r.Done(), true)
		if r.Rest() != nil {
			t.Fatal("r.Rest is not nil")
		}
	})

	t.Run("readform advances", func(t *testing.T) {
		r := reader.ReadStr("(1) 2")
		first, err := r.ReadForm()
		assert.NoErr(t, err)
		assert.MalEqual(t, first, lst(num(1)))
		assert.Equal(t, r.Done(), false)
		rest := r.Rest()
		assert.Equal(t, len(rest), 1)
		assert.Equal(t, rest[0].Value, "2")
		second, err := r.ReadForm()
		assert.NoErr(t, err)
		assert.MalEqual(t, second, num(2))
		assert.Equal(t, r.Done(), true)
	})
}
