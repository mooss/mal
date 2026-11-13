//revive:disable:cognitive-complexity,line-length-limit,function-length
package lexer_test

import (
	"testing"

	"github.com/mooss/mal/src/lexer"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
)

const (
	DQS      = lexer.TOKEN_DQS
	INT      = lexer.TOKEN_INT
	LBRACE   = lexer.TOKEN_LBRACE
	LBRACKET = lexer.TOKEN_LBRACKET
	LPAREN   = lexer.TOKEN_LPAREN
	RBRACE   = lexer.TOKEN_RBRACE
	RBRACKET = lexer.TOKEN_RBRACKET
	RPAREN   = lexer.TOKEN_RPAREN
	SQUOTE   = lexer.TOKEN_SQUOTE
	SYMBOL   = lexer.TOKEN_SYMBOL
)

func tok(typ lexer.TokenType, val string) lexer.Token {
	return lexer.Token{Type: typ, Value: val}
}

func collectTokens(input string) ([]lexer.Token, bool) {
	seq, errFn := lexer.TokenIterator(input)
	var got []lexer.Token
	for t := range seq {
		got = append(got, t)
	}
	return got, errFn() != nil
}

func runTokenCases(t *testing.T, tests []tut.Case[[]lexer.Token]) {
	t.Helper()
	tut.Run(t, tests, func(t *testing.T, c tut.Case[[]lexer.Token]) {
		t.Helper()
		got, gotErr := collectTokens(c.Input)
		assert.Equal(t, got, c.Want)
		assert.Equal(t, gotErr, c.WantErr)
	})
}

func dqsTokens(vals []string) []lexer.Token {
	if vals == nil {
		return nil
	}
	got := make([]lexer.Token, len(vals))
	for i, v := range vals {
		got[i] = tok(DQS, v)
	}
	return got
}

func runStringCases(t *testing.T, tests []tut.Case[[]string]) {
	t.Helper()
	converted := make([]tut.Case[[]lexer.Token], len(tests))
	for i, tt := range tests {
		converted[i] = tut.Case[[]lexer.Token]{Name: tt.Name, Input: tt.Input, Want: dqsTokens(tt.Want), WantErr: tt.WantErr}
	}
	runTokenCases(t, converted)
}

func singleCharCases(char string, typ lexer.TokenType) []tut.Case[[]lexer.Token] {
	one := tok(typ, char)
	symFoo := tok(SYMBOL, "foo")
	symBar := tok(SYMBOL, "bar")
	int123 := tok(INT, "123")
	return []tut.Case[[]lexer.Token]{
		{Name: "single", Input: char, Want: []lexer.Token{one}},
		{Name: "doubled", Input: char + char, Want: []lexer.Token{one, one}},
		{Name: "spaced", Input: char + " " + char, Want: []lexer.Token{one, one}},
		{Name: "whitespace padded", Input: "  " + char + "\n\t ", Want: []lexer.Token{one}},
		{Name: "symbol before", Input: "foo" + char, Want: []lexer.Token{symFoo, one}},
		{Name: "symbol after", Input: char + "foo", Want: []lexer.Token{one, symFoo}},
		{Name: "int before", Input: "123" + char, Want: []lexer.Token{int123, one}},
		{Name: "int after", Input: char + "123", Want: []lexer.Token{one, int123}},
		{Name: "symbol spaced", Input: "foo " + char + " bar", Want: []lexer.Token{symFoo, one, symBar}},
	}
}

func TestTokenizeStrings(t *testing.T) {
	tests := []tut.Case[[]string]{
		{Name: "empty input", Input: ""},
		{Name: "empty string", Input: `""`, Want: []string{`""`}},
		{Name: "simple", Input: `"hello"`, Want: []string{`"hello"`}},
		{Name: "with space", Input: `"hello world"`, Want: []string{`"hello world"`}},
		{Name: "digits", Input: `"123"`, Want: []string{`"123"`}},
		{Name: "unicode", Input: `"héllo"`, Want: []string{`"héllo"`}},
		{Name: "emoji", Input: `"😀"`, Want: []string{`"😀"`}},
		{Name: "escaped quote", Input: `"a\"b"`, Want: []string{`"a\"b"`}},
		{Name: "escaped backslash", Input: `"a\\b"`, Want: []string{`"a\\b"`}},
		{Name: "single escaped backslash", Input: `"\\"`, Want: []string{`"\\"`}},
		{Name: "single escaped quote", Input: `"\""`, Want: []string{`"\""`}},
		{Name: "backslash n escape", Input: `"a\nb"`, Want: []string{`"a\nb"`}},
		{Name: "unterminated", Input: `"hello`, WantErr: true},
		{Name: "single quote", Input: `"`, WantErr: true},
		{Name: "ends with escaped quote", Input: `"hello\"`, WantErr: true},
		{Name: "raw newline no closing", Input: "\"foo\n", WantErr: true},
		{Name: "raw newline with closing", Input: "\"a\nb\"", Want: []string{"\"a\nb\""}},
		{Name: "CR with closing", Input: "\"a\rb\"", Want: []string{"\"a\rb\""}},
		{Name: "invalid utf8 with closing", Input: "\"a\xffb\"", Want: []string{"\"a�b\""}},
		{Name: "NUL no closing", Input: "\"foo\x00", WantErr: true},
		{Name: "invalid no closing", Input: "\"foo\xff", WantErr: true},
		{Name: "two adjacent", Input: `"a""b"`, Want: []string{`"a"`, `"b"`}},
		{Name: "trailing escaped backslash", Input: `"foo\\"`, Want: []string{`"foo\\"`}},
	}

	runStringCases(t, tests)
}

func testSingleCharToken(t *testing.T, char string, typ lexer.TokenType) {
	t.Helper()
	runTokenCases(t, singleCharCases(char, typ))
}

func TestTokenizeSingleCharTokens(t *testing.T) {
	tut.Run(t, []tut.Case[lexer.TokenType]{
		{Name: "lparen", Input: "(", Want: LPAREN},
		{Name: "rparen", Input: ")", Want: RPAREN},
		{Name: "lbracket", Input: "[", Want: LBRACKET},
		{Name: "rbracket", Input: "]", Want: RBRACKET},
		{Name: "lbrace", Input: "{", Want: LBRACE},
		{Name: "rbrace", Input: "}", Want: RBRACE},
		{Name: "squote", Input: "'", Want: SQUOTE},
	}, func(t *testing.T, c tut.Case[lexer.TokenType]) {
		t.Helper()
		testSingleCharToken(t, c.Input, c.Want)
	})
}

func TestTokenizeInts(t *testing.T) {
	tests := []tut.Case[[]lexer.Token]{
		{Name: "zero", Input: "0", Want: []lexer.Token{{Type: INT, Value: "0"}}},
		{Name: "single digit", Input: "1", Want: []lexer.Token{{Type: INT, Value: "1"}}},
		{Name: "multi digit", Input: "123", Want: []lexer.Token{{Type: INT, Value: "123"}}},
		{Name: "leading zeros", Input: "007", Want: []lexer.Token{{Type: INT, Value: "007"}}},
		{Name: "negative zero", Input: "-0", Want: []lexer.Token{{Type: INT, Value: "-0"}}},
		{Name: "negative single", Input: "-1", Want: []lexer.Token{{Type: INT, Value: "-1"}}},
		{Name: "negative multi", Input: "-42", Want: []lexer.Token{{Type: INT, Value: "-42"}}},
		{Name: "negative long", Input: "-9876543210", Want: []lexer.Token{{Type: INT, Value: "-9876543210"}}},
		{Name: "large", Input: "999999999999999999999999", Want: []lexer.Token{{Type: INT, Value: "999999999999999999999999"}}},
		{Name: "multiple", Input: "1 2 3", Want: []lexer.Token{{Type: INT, Value: "1"}, {Type: INT, Value: "2"}, {Type: INT, Value: "3"}}},
		{Name: "whitespace padded", Input: "  42\t\n-7  ", Want: []lexer.Token{{Type: INT, Value: "42"}, {Type: INT, Value: "-7"}}},
		{Name: "zeros", Input: "0 -0 00 -00", Want: []lexer.Token{{Type: INT, Value: "0"}, {Type: INT, Value: "-0"}, {Type: INT, Value: "00"}, {Type: INT, Value: "-00"}}},
		{Name: "plus digits", Input: "+123", Want: []lexer.Token{{Type: INT, Value: "+123"}}},
	}

	runTokenCases(t, tests)
}

func TestTokenizeSymbols(t *testing.T) {
	tests := []tut.Case[[]lexer.Token]{
		{Name: "single letter", Input: "a", Want: []lexer.Token{{Type: SYMBOL, Value: "a"}}},
		{Name: "simple", Input: "foo", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}}},
		{Name: "uppercase", Input: "ABC", Want: []lexer.Token{{Type: SYMBOL, Value: "ABC"}}},
		{Name: "mixed case digits", Input: "AbC123", Want: []lexer.Token{{Type: SYMBOL, Value: "AbC123"}}},
		{Name: "plus", Input: "+", Want: []lexer.Token{{Type: SYMBOL, Value: "+"}}},
		{Name: "star", Input: "*", Want: []lexer.Token{{Type: SYMBOL, Value: "*"}}},
		{Name: "single minus", Input: "-", Want: []lexer.Token{{Type: SYMBOL, Value: "-"}}},
		{Name: "double minus", Input: "--", Want: []lexer.Token{{Type: SYMBOL, Value: "--"}}},
		{Name: "double minus digit", Input: "--1", Want: []lexer.Token{{Type: SYMBOL, Value: "--1"}}},
		{Name: "minus letter", Input: "-a", Want: []lexer.Token{{Type: SYMBOL, Value: "-a"}}},
		{Name: "letter minus", Input: "a-", Want: []lexer.Token{{Type: SYMBOL, Value: "a-"}}},
		{Name: "dash", Input: "foo-bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo-bar"}}},
		{Name: "slash", Input: "foo/bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo/bar"}}},
		{Name: "dot", Input: "foo.bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo.bar"}}},
		{Name: "dots", Input: "...", Want: []lexer.Token{{Type: SYMBOL, Value: "..."}}},
		{Name: "stars", Input: "***", Want: []lexer.Token{{Type: SYMBOL, Value: "***"}}},
		{Name: "bang", Input: "!", Want: []lexer.Token{{Type: SYMBOL, Value: "!"}}},
		{Name: "double bang", Input: "!!", Want: []lexer.Token{{Type: SYMBOL, Value: "!!"}}},
		{Name: "leading bang", Input: "!foo", Want: []lexer.Token{{Type: SYMBOL, Value: "!foo"}}},
		{Name: "trailing bang", Input: "foo!", Want: []lexer.Token{{Type: SYMBOL, Value: "foo!"}}},
		{Name: "middle bang", Input: "foo!bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo!bar"}}},
		{Name: "bang with digits", Input: "!123", Want: []lexer.Token{{Type: SYMBOL, Value: "!123"}}},
		{Name: "punct run", Input: "*+-./", Want: []lexer.Token{{Type: SYMBOL, Value: "*+-./"}}},
		{Name: "punct run with bang", Input: "!*+-./!", Want: []lexer.Token{{Type: SYMBOL, Value: "!*+-./!"}}},
		{Name: "multiple", Input: "foo bar baz", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}, {Type: SYMBOL, Value: "baz"}}},
		{Name: "all plain punct", Input: "a/b.c*d+e-f", Want: []lexer.Token{{Type: SYMBOL, Value: "a/b.c*d+e-f"}}},
		{Name: "keywords", Input: "nil true false defn", Want: []lexer.Token{{Type: SYMBOL, Value: "nil"}, {Type: SYMBOL, Value: "true"}, {Type: SYMBOL, Value: "false"}, {Type: SYMBOL, Value: "defn"}}},
		{Name: "digit minus", Input: "1-", WantErr: true},
		{Name: "digits plus digits", Input: "12+34", WantErr: true},
		{Name: "ampersand", Input: "&", Want: []lexer.Token{{Type: SYMBOL, Value: "&"}}},
	}

	runTokenCases(t, tests)
}

func TestTokenizeCombinations(t *testing.T) {
	tests := []tut.Case[[]lexer.Token]{
		{Name: "empty input", Input: ""},
		{Name: "whitespace only", Input: " \t\n\r "},
		{Name: "empty list", Input: "()", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: RPAREN, Value: ")"}}},
		{Name: "symbol and int", Input: "(foo 123)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: INT, Value: "123"}, {Type: RPAREN, Value: ")"}}},
		{Name: "plus ints", Input: "(+ 1 2)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "+"}, {Type: INT, Value: "1"}, {Type: INT, Value: "2"}, {Type: RPAREN, Value: ")"}}},
		{Name: "def", Input: "(def foo 42)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "def"}, {Type: SYMBOL, Value: "foo"}, {Type: INT, Value: "42"}, {Type: RPAREN, Value: ")"}}},
		{Name: "int then symbol", Input: "123 abc", Want: []lexer.Token{{Type: INT, Value: "123"}, {Type: SYMBOL, Value: "abc"}}},
		{Name: "int vs symbol dash", Input: "-123 -abc +123", Want: []lexer.Token{{Type: INT, Value: "-123"}, {Type: SYMBOL, Value: "-abc"}, {Type: INT, Value: "+123"}}},
		{Name: "no spaces", Input: "(foo)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: RPAREN, Value: ")"}}},
		{Name: "int in parens", Input: "(1)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: INT, Value: "1"}, {Type: RPAREN, Value: ")"}}},
		{Name: "nested ints", Input: "((1 2) (3 4))", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: LPAREN, Value: "("}, {Type: INT, Value: "1"}, {Type: INT, Value: "2"}, {Type: RPAREN, Value: ")"}, {Type: LPAREN, Value: "("}, {Type: INT, Value: "3"}, {Type: INT, Value: "4"}, {Type: RPAREN, Value: ")"}, {Type: RPAREN, Value: ")"}}},
		{Name: "string and symbol", Input: `"hello" foo`, Want: []lexer.Token{{Type: DQS, Value: `"hello"`}, {Type: SYMBOL, Value: "foo"}}},
		{Name: "string in list", Input: `(foo "bar baz" -42)`, Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: DQS, Value: `"bar baz"`}, {Type: INT, Value: "-42"}, {Type: RPAREN, Value: ")"}}},
		{Name: "int adjacent parens", Input: "1(2)", Want: []lexer.Token{{Type: INT, Value: "1"}, {Type: LPAREN, Value: "("}, {Type: INT, Value: "2"}, {Type: RPAREN, Value: ")"}}},
		{Name: "all kinds", Input: `123 abc ( ) "hello"`, Want: []lexer.Token{{Type: INT, Value: "123"}, {Type: SYMBOL, Value: "abc"}, {Type: LPAREN, Value: "("}, {Type: RPAREN, Value: ")"}, {Type: DQS, Value: `"hello"`}}},
		{Name: "multiline", Input: "  ( \n foo\t 123 \r\n )  ", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: INT, Value: "123"}, {Type: RPAREN, Value: ")"}}},
		{Name: "unterminated string", Input: `"foo`, WantErr: true},
		{Name: "valid then invalid", Input: "foo : bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: ":"}, {Type: SYMBOL, Value: "bar"}}},

		{Name: "symbol boundary lparen", Input: "foo(bar)", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "bar"}, {Type: RPAREN, Value: ")"}}},
		{Name: "symbol boundary rparen", Input: "foo)bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: RPAREN, Value: ")"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "int boundary lparen", Input: "123(foo)", Want: []lexer.Token{{Type: INT, Value: "123"}, {Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: RPAREN, Value: ")"}}},
		{Name: "string adjacent symbol", Input: `"foo"bar`, Want: []lexer.Token{{Type: DQS, Value: `"foo"`}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "string adjacent int", Input: `"foo"123`, Want: []lexer.Token{{Type: DQS, Value: `"foo"`}, {Type: INT, Value: "123"}}},
		{Name: "string adjacent lparen", Input: `"foo"(bar)`, Want: []lexer.Token{{Type: DQS, Value: `"foo"`}, {Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "bar"}, {Type: RPAREN, Value: ")"}}},
		{Name: "escaped newline", Input: "\"foo\\\nbar\"", Want: []lexer.Token{{Type: DQS, Value: "\"foo\\\nbar\""}}},
	}

	runTokenCases(t, tests)
}

func TestTokenizeInvalid(t *testing.T) {
	tests := []tut.Case[[]lexer.Token]{
		{Name: "backtick", Input: "`", WantErr: true},
		{Name: "tilde", Input: "~", WantErr: true},
		{Name: "hash", Input: "#", WantErr: true},
		{Name: "dollar", Input: "$", WantErr: true},
		{Name: "percent", Input: "%", WantErr: true},
		{Name: "pipe", Input: "|", WantErr: true},
		{Name: "backslash", Input: `\`, WantErr: true},
		{Name: "caret", Input: "^", WantErr: true},
		{Name: "underscore", Input: "_", WantErr: true},
		{Name: "int then letters", Input: "123abc", WantErr: true},
		{Name: "int then slash", Input: "123/foo", WantErr: true},
		{Name: "symbol then quote", Input: `foo"bar"`, WantErr: true},
		{Name: "int then quote", Input: `123"foo"`, WantErr: true},
		{Name: "trailing backslash", Input: "\"foo\\", WantErr: true},
	}

	runTokenCases(t, tests)
}

func TestTokenizeWhitespace(t *testing.T) {
	tests := []tut.Case[[]lexer.Token]{
		{Name: "space only", Input: " "},
		{Name: "tab only", Input: "\t"},
		{Name: "newline only", Input: "\n"},
		{Name: "carriage return only", Input: "\r"},
		{Name: "single comma", Input: ","},
		{Name: "multiple commas", Input: ",,,"},
		{Name: "mixed whitespace only", Input: " \t\n\r"},
		{Name: "mixed whitespace and commas only", Input: " , ,\t\n\r, "},
		{Name: "commas with newlines only", Input: ",\n,\r\n,\t,"},
		{Name: "symbol padded spaces", Input: "  foo  ", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}}},
		{Name: "symbol padded tabs", Input: "\tfoo\t", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}}},
		{Name: "symbol padded newlines", Input: "\nfoo\n", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}}},
		{Name: "symbol padded carriage returns", Input: "\rfoo\r", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}}},
		{Name: "symbol padded commas", Input: ",foo,", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}}},
		{Name: "symbol padded mixed", Input: " \t\n\r,foo, \r\n\t ", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}}},
		{Name: "int padded mixed", Input: " ,\t\n\r42\r\n\t, ", Want: []lexer.Token{{Type: INT, Value: "42"}}},
		{Name: "symbols space separated", Input: "foo bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "symbols tab separated", Input: "foo\tbar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "symbols newline separated", Input: "foo\nbar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "symbols cr separated", Input: "foo\rbar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "symbols comma separated", Input: "foo,bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "symbols multi comma separated", Input: "foo,,,bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "symbols comma space mixed", Input: "foo , bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "symbols mixed whitespace separated", Input: "foo \t\n\r, bar", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "ints space separated", Input: "1 2", Want: []lexer.Token{{Type: INT, Value: "1"}, {Type: INT, Value: "2"}}},
		{Name: "ints tab separated", Input: "1\t2", Want: []lexer.Token{{Type: INT, Value: "1"}, {Type: INT, Value: "2"}}},
		{Name: "ints newline separated", Input: "1\n2", Want: []lexer.Token{{Type: INT, Value: "1"}, {Type: INT, Value: "2"}}},
		{Name: "ints cr separated", Input: "1\r2", Want: []lexer.Token{{Type: INT, Value: "1"}, {Type: INT, Value: "2"}}},
		{Name: "ints comma separated", Input: "1,2,3", Want: []lexer.Token{{Type: INT, Value: "1"}, {Type: INT, Value: "2"}, {Type: INT, Value: "3"}}},
		{Name: "ints mixed whitespace separated", Input: "1 \t\n\r, 2,\t3", Want: []lexer.Token{{Type: INT, Value: "1"}, {Type: INT, Value: "2"}, {Type: INT, Value: "3"}}},
		{Name: "int comma symbol", Input: "123,foo", Want: []lexer.Token{{Type: INT, Value: "123"}, {Type: SYMBOL, Value: "foo"}}},
		{Name: "symbol comma int", Input: "foo,123", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: INT, Value: "123"}}},
		{Name: "list comma separated", Input: "(foo,123)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: INT, Value: "123"}, {Type: RPAREN, Value: ")"}}},
		{Name: "list commas and spaces", Input: "(,foo,,123,,)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: INT, Value: "123"}, {Type: RPAREN, Value: ")"}}},
		{Name: "parens only commas", Input: "(,)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: RPAREN, Value: ")"}}},
		{Name: "empty parens with whitespace", Input: "( \t\n\r, )", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: RPAREN, Value: ")"}}},
		{Name: "leading commas before list", Input: ",,(foo)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: RPAREN, Value: ")"}}},
		{Name: "trailing commas after list", Input: "(foo),,", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: RPAREN, Value: ")"}}},
		{Name: "comma between parens", Input: "(1),(2)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: INT, Value: "1"}, {Type: RPAREN, Value: ")"}, {Type: LPAREN, Value: "("}, {Type: INT, Value: "2"}, {Type: RPAREN, Value: ")"}}},
		{Name: "crlf separated list", Input: "(foo\r\nbar)", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: SYMBOL, Value: "bar"}, {Type: RPAREN, Value: ")"}}},
		{Name: "symbol comma lparen", Input: "foo,(bar)", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}, {Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "bar"}, {Type: RPAREN, Value: ")"}}},
		{Name: "string padded whitespace", Input: "  \"hello\"  ", Want: []lexer.Token{{Type: DQS, Value: `"hello"`}}},
		{Name: "string padded commas", Input: ",\"hello\",", Want: []lexer.Token{{Type: DQS, Value: `"hello"`}}},
		{Name: "string padded mixed", Input: " \t\n\r,\"hello\", \r\n\t ", Want: []lexer.Token{{Type: DQS, Value: `"hello"`}}},
		{Name: "string with comma inside", Input: "\"a,b\"", Want: []lexer.Token{{Type: DQS, Value: `"a,b"`}}},
		{Name: "string comma symbol", Input: "\"foo\",bar", Want: []lexer.Token{{Type: DQS, Value: `"foo"`}, {Type: SYMBOL, Value: "bar"}}},
		{Name: "string space int", Input: "\"foo\" 123", Want: []lexer.Token{{Type: DQS, Value: `"foo"`}, {Type: INT, Value: "123"}}},
		{Name: "multiline mixed", Input: "  ( \n foo\t,123 \r\n )  ", Want: []lexer.Token{{Type: LPAREN, Value: "("}, {Type: SYMBOL, Value: "foo"}, {Type: INT, Value: "123"}, {Type: RPAREN, Value: ")"}}},
		{Name: "whitespace before invalid", Input: "foo , à", Want: []lexer.Token{{Type: SYMBOL, Value: "foo"}}, WantErr: true},
	}

	runTokenCases(t, tests)
}

func TestTokenIteratorEarlyBreak(t *testing.T) {
	seq, errFn := lexer.TokenIterator("(1 2 3)")
	var got []lexer.Token
	for tok := range seq {
		got = append(got, tok)
		break
	}
	assert.Equal(t, len(got), 1)
	assert.Equal(t, got[0], lexer.Token{Type: LPAREN, Value: "("})
	assert.NoErr(t, errFn())
}
