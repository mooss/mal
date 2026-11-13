//revive:disable:cognitive-complexity,line-length-limit,function-length
package lexer

import "testing"

type tokenCase struct {
	name    string
	input   string
	want    []Token
	wantErr bool
}

type stringCase struct {
	name       string
	input      string
	wantTokens []string
	wantErr    bool
}

func tok(typ TokenType, val string) Token {
	return Token{Type: typ, Value: val}
}

func collectTokens(input string) ([]Token, bool) {
	seq, errFn := TokenIterator(input)
	var got []Token
	for t := range seq {
		got = append(got, t)
	}
	return got, errFn() != nil
}

func checkTokens(t *testing.T, input string, got, want []Token, gotErr, wantErr bool) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("input %q: got %d tokens %v (err %v), want %d tokens %v", input, len(got), got, gotErr, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("input %q token %d: got %+v, want %+v", input, i, got[i], want[i])
		}
	}
	if gotErr != wantErr {
		t.Errorf("input %q: gotErr %v, wantErr %v", input, gotErr, wantErr)
	}
}

func runTokenCases(t *testing.T, tests []tokenCase) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := collectTokens(tt.input)
			checkTokens(t, tt.input, got, tt.want, gotErr, tt.wantErr)
		})
	}
}

func dqsTokens(vals []string) []Token {
	if vals == nil {
		return nil
	}
	got := make([]Token, len(vals))
	for i, v := range vals {
		got[i] = tok(TOKEN_DQS, v)
	}
	return got
}

func runStringCases(t *testing.T, tests []stringCase) {
	t.Helper()
	converted := make([]tokenCase, len(tests))
	for i, tt := range tests {
		converted[i] = tokenCase{name: tt.name, input: tt.input, want: dqsTokens(tt.wantTokens), wantErr: tt.wantErr}
	}
	runTokenCases(t, converted)
}

func singleCharCases(char string, typ TokenType) []tokenCase {
	one := tok(typ, char)
	symFoo := tok(TOKEN_SYMBOL, "foo")
	symBar := tok(TOKEN_SYMBOL, "bar")
	int123 := tok(TOKEN_INT, "123")
	return []tokenCase{
		{name: "single", input: char, want: []Token{one}, wantErr: false},
		{name: "doubled", input: char + char, want: []Token{one, one}, wantErr: false},
		{name: "spaced", input: char + " " + char, want: []Token{one, one}, wantErr: false},
		{name: "whitespace padded", input: "  " + char + "\n\t ", want: []Token{one}, wantErr: false},
		{name: "symbol before", input: "foo" + char, want: []Token{symFoo, one}, wantErr: false},
		{name: "symbol after", input: char + "foo", want: []Token{one, symFoo}, wantErr: false},
		{name: "int before", input: "123" + char, want: []Token{int123, one}, wantErr: false},
		{name: "int after", input: char + "123", want: []Token{one, int123}, wantErr: false},
		{name: "symbol spaced", input: "foo " + char + " bar", want: []Token{symFoo, one, symBar}, wantErr: false},
	}
}

func TestTokenizeStrings(t *testing.T) {
	tests := []stringCase{
		{name: "empty input", input: "", wantTokens: nil, wantErr: false},
		{name: "empty string", input: `""`, wantTokens: []string{`""`}, wantErr: false},
		{name: "simple", input: `"hello"`, wantTokens: []string{`"hello"`}, wantErr: false},
		{name: "with space", input: `"hello world"`, wantTokens: []string{`"hello world"`}, wantErr: false},
		{name: "digits", input: `"123"`, wantTokens: []string{`"123"`}, wantErr: false},
		{name: "unicode", input: `"héllo"`, wantTokens: []string{`"héllo"`}, wantErr: false},
		{name: "emoji", input: `"😀"`, wantTokens: []string{`"😀"`}, wantErr: false},
		{name: "escaped quote", input: `"a\"b"`, wantTokens: []string{`"a\"b"`}, wantErr: false},
		{name: "escaped backslash", input: `"a\\b"`, wantTokens: []string{`"a\\b"`}, wantErr: false},
		{name: "single escaped backslash", input: `"\\"`, wantTokens: []string{`"\\"`}, wantErr: false},
		{name: "single escaped quote", input: `"\""`, wantTokens: []string{`"\""`}, wantErr: false},
		{name: "backslash n escape", input: `"a\nb"`, wantTokens: []string{`"a\nb"`}, wantErr: false},
		{name: "unterminated", input: `"hello`, wantTokens: []string{}, wantErr: true},
		{name: "single quote", input: `"`, wantTokens: []string{}, wantErr: true},
		{name: "ends with escaped quote", input: `"hello\"`, wantTokens: []string{}, wantErr: true},
		{name: "raw newline no closing", input: "\"foo\n", wantTokens: []string{}, wantErr: true},
		{name: "raw newline with closing", input: "\"a\nb\"", wantTokens: []string{"\"a\nb\""}, wantErr: false},
		{name: "CR with closing", input: "\"a\rb\"", wantTokens: []string{"\"a\rb\""}, wantErr: false},
		{name: "invalid utf8 with closing", input: "\"a\xffb\"", wantTokens: []string{"\"a�b\""}, wantErr: false},
		{name: "NUL no closing", input: "\"foo\x00", wantTokens: []string{}, wantErr: true},
		{name: "invalid no closing", input: "\"foo\xff", wantTokens: []string{}, wantErr: true},
		{name: "two adjacent", input: `"a""b"`, wantTokens: []string{`"a"`, `"b"`}, wantErr: false},
		{name: "trailing escaped backslash", input: `"foo\\"`, wantTokens: []string{`"foo\\"`}, wantErr: false},
	}

	runStringCases(t, tests)
}

func testSingleCharToken(t *testing.T, char string, typ TokenType) {
	t.Helper()
	runTokenCases(t, singleCharCases(char, typ))
}

func TestTokenizeSingleCharTokens(t *testing.T) {
	cases := []struct {
		name string
		char string
		typ  TokenType
	}{
		{name: "lparen", char: "(", typ: TOKEN_LPAREN},
		{name: "rparen", char: ")", typ: TOKEN_RPAREN},
		{name: "lbracket", char: "[", typ: TOKEN_LBRACKET},
		{name: "rbracket", char: "]", typ: TOKEN_RBRACKET},
		{name: "lbrace", char: "{", typ: TOKEN_LBRACE},
		{name: "rbrace", char: "}", typ: TOKEN_RBRACE},
		{name: "squote", char: "'", typ: TOKEN_SQUOTE},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testSingleCharToken(t, tc.char, tc.typ)
		})
	}
}

func TestTokenizeInts(t *testing.T) {
	tests := []tokenCase{
		{name: "zero", input: "0", want: []Token{{Type: TOKEN_INT, Value: "0"}}, wantErr: false},
		{name: "single digit", input: "1", want: []Token{{Type: TOKEN_INT, Value: "1"}}, wantErr: false},
		{name: "multi digit", input: "123", want: []Token{{Type: TOKEN_INT, Value: "123"}}, wantErr: false},
		{name: "leading zeros", input: "007", want: []Token{{Type: TOKEN_INT, Value: "007"}}, wantErr: false},
		{name: "negative zero", input: "-0", want: []Token{{Type: TOKEN_INT, Value: "-0"}}, wantErr: false},
		{name: "negative single", input: "-1", want: []Token{{Type: TOKEN_INT, Value: "-1"}}, wantErr: false},
		{name: "negative multi", input: "-42", want: []Token{{Type: TOKEN_INT, Value: "-42"}}, wantErr: false},
		{name: "negative long", input: "-9876543210", want: []Token{{Type: TOKEN_INT, Value: "-9876543210"}}, wantErr: false},
		{name: "large", input: "999999999999999999999999", want: []Token{{Type: TOKEN_INT, Value: "999999999999999999999999"}}, wantErr: false},
		{name: "multiple", input: "1 2 3", want: []Token{{Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}, {Type: TOKEN_INT, Value: "3"}}, wantErr: false},
		{name: "whitespace padded", input: "  42\t\n-7  ", want: []Token{{Type: TOKEN_INT, Value: "42"}, {Type: TOKEN_INT, Value: "-7"}}, wantErr: false},
		{name: "zeros", input: "0 -0 00 -00", want: []Token{{Type: TOKEN_INT, Value: "0"}, {Type: TOKEN_INT, Value: "-0"}, {Type: TOKEN_INT, Value: "00"}, {Type: TOKEN_INT, Value: "-00"}}, wantErr: false},
		{name: "plus digits", input: "+123", want: []Token{{Type: TOKEN_INT, Value: "+123"}}, wantErr: false},
	}

	runTokenCases(t, tests)
}

func TestTokenizeSymbols(t *testing.T) {
	tests := []tokenCase{
		{name: "single letter", input: "a", want: []Token{{Type: TOKEN_SYMBOL, Value: "a"}}, wantErr: false},
		{name: "simple", input: "foo", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "uppercase", input: "ABC", want: []Token{{Type: TOKEN_SYMBOL, Value: "ABC"}}, wantErr: false},
		{name: "mixed case digits", input: "AbC123", want: []Token{{Type: TOKEN_SYMBOL, Value: "AbC123"}}, wantErr: false},
		{name: "plus", input: "+", want: []Token{{Type: TOKEN_SYMBOL, Value: "+"}}, wantErr: false},
		{name: "star", input: "*", want: []Token{{Type: TOKEN_SYMBOL, Value: "*"}}, wantErr: false},
		{name: "single minus", input: "-", want: []Token{{Type: TOKEN_SYMBOL, Value: "-"}}, wantErr: false},
		{name: "double minus", input: "--", want: []Token{{Type: TOKEN_SYMBOL, Value: "--"}}, wantErr: false},
		{name: "double minus digit", input: "--1", want: []Token{{Type: TOKEN_SYMBOL, Value: "--1"}}, wantErr: false},
		{name: "minus letter", input: "-a", want: []Token{{Type: TOKEN_SYMBOL, Value: "-a"}}, wantErr: false},
		{name: "letter minus", input: "a-", want: []Token{{Type: TOKEN_SYMBOL, Value: "a-"}}, wantErr: false},
		{name: "dash", input: "foo-bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo-bar"}}, wantErr: false},
		{name: "slash", input: "foo/bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo/bar"}}, wantErr: false},
		{name: "dot", input: "foo.bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo.bar"}}, wantErr: false},
		{name: "dots", input: "...", want: []Token{{Type: TOKEN_SYMBOL, Value: "..."}}, wantErr: false},
		{name: "stars", input: "***", want: []Token{{Type: TOKEN_SYMBOL, Value: "***"}}, wantErr: false},
		{name: "bang", input: "!", want: []Token{{Type: TOKEN_SYMBOL, Value: "!"}}, wantErr: false},
		{name: "double bang", input: "!!", want: []Token{{Type: TOKEN_SYMBOL, Value: "!!"}}, wantErr: false},
		{name: "leading bang", input: "!foo", want: []Token{{Type: TOKEN_SYMBOL, Value: "!foo"}}, wantErr: false},
		{name: "trailing bang", input: "foo!", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo!"}}, wantErr: false},
		{name: "middle bang", input: "foo!bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo!bar"}}, wantErr: false},
		{name: "bang with digits", input: "!123", want: []Token{{Type: TOKEN_SYMBOL, Value: "!123"}}, wantErr: false},
		{name: "punct run", input: "*+-./", want: []Token{{Type: TOKEN_SYMBOL, Value: "*+-./"}}, wantErr: false},
		{name: "punct run with bang", input: "!*+-./!", want: []Token{{Type: TOKEN_SYMBOL, Value: "!*+-./!"}}, wantErr: false},
		{name: "multiple", input: "foo bar baz", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}, {Type: TOKEN_SYMBOL, Value: "baz"}}, wantErr: false},
		{name: "all plain punct", input: "a/b.c*d+e-f", want: []Token{{Type: TOKEN_SYMBOL, Value: "a/b.c*d+e-f"}}, wantErr: false},
		{name: "keywords", input: "nil true false defn", want: []Token{{Type: TOKEN_SYMBOL, Value: "nil"}, {Type: TOKEN_SYMBOL, Value: "true"}, {Type: TOKEN_SYMBOL, Value: "false"}, {Type: TOKEN_SYMBOL, Value: "defn"}}, wantErr: false},
		{name: "digit minus", input: "1-", want: []Token{}, wantErr: true},
		{name: "digits plus digits", input: "12+34", want: []Token{}, wantErr: true},
		{name: "ampersand", input: "&", want: []Token{{Type: TOKEN_SYMBOL, Value: "&"}}, wantErr: false},
	}

	runTokenCases(t, tests)
}

func TestTokenizeCombinations(t *testing.T) {
	tests := []tokenCase{
		{name: "empty input", input: "", want: nil, wantErr: false},
		{name: "whitespace only", input: " \t\n\r ", want: nil, wantErr: false},
		{name: "empty list", input: "()", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "symbol and int", input: "(foo 123)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "plus ints", input: "(+ 1 2)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "+"}, {Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "def", input: "(def foo 42)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "def"}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_INT, Value: "42"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "int then symbol", input: "123 abc", want: []Token{{Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_SYMBOL, Value: "abc"}}, wantErr: false},
		{name: "int vs symbol dash", input: "-123 -abc +123", want: []Token{{Type: TOKEN_INT, Value: "-123"}, {Type: TOKEN_SYMBOL, Value: "-abc"}, {Type: TOKEN_INT, Value: "+123"}}, wantErr: false},
		{name: "no spaces", input: "(foo)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "int in parens", input: "(1)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "nested ints", input: "((1 2) (3 4))", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}, {Type: TOKEN_RPAREN, Value: ")"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_INT, Value: "3"}, {Type: TOKEN_INT, Value: "4"}, {Type: TOKEN_RPAREN, Value: ")"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "string and symbol", input: `"hello" foo`, want: []Token{{Type: TOKEN_DQS, Value: `"hello"`}, {Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "string in list", input: `(foo "bar baz" -42)`, want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_DQS, Value: `"bar baz"`}, {Type: TOKEN_INT, Value: "-42"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "int adjacent parens", input: "1(2)", want: []Token{{Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_INT, Value: "2"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "all kinds", input: `123 abc ( ) "hello"`, want: []Token{{Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_SYMBOL, Value: "abc"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_RPAREN, Value: ")"}, {Type: TOKEN_DQS, Value: `"hello"`}}, wantErr: false},
		{name: "multiline", input: "  ( \n foo\t 123 \r\n )  ", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "unterminated string", input: `"foo`, want: nil, wantErr: true},
		{name: "valid then invalid", input: "foo : bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: ":"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},

		{name: "symbol boundary lparen", input: "foo(bar)", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "bar"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "symbol boundary rparen", input: "foo)bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_RPAREN, Value: ")"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "int boundary lparen", input: "123(foo)", want: []Token{{Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "string adjacent symbol", input: `"foo"bar`, want: []Token{{Type: TOKEN_DQS, Value: `"foo"`}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "string adjacent int", input: `"foo"123`, want: []Token{{Type: TOKEN_DQS, Value: `"foo"`}, {Type: TOKEN_INT, Value: "123"}}, wantErr: false},
		{name: "string adjacent lparen", input: `"foo"(bar)`, want: []Token{{Type: TOKEN_DQS, Value: `"foo"`}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "bar"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "escaped newline", input: "\"foo\\\nbar\"", want: []Token{{Type: TOKEN_DQS, Value: "\"foo\\\nbar\""}}, wantErr: false},
	}

	runTokenCases(t, tests)
}

func TestTokenizeInvalid(t *testing.T) {
	tests := []tokenCase{
		{name: "backtick", input: "`", want: nil, wantErr: true},
		{name: "tilde", input: "~", want: nil, wantErr: true},
		{name: "hash", input: "#", want: nil, wantErr: true},
		{name: "dollar", input: "$", want: nil, wantErr: true},
		{name: "percent", input: "%", want: nil, wantErr: true},
		{name: "pipe", input: "|", want: nil, wantErr: true},
		{name: "backslash", input: `\`, want: nil, wantErr: true},
		{name: "caret", input: "^", want: nil, wantErr: true},
		{name: "underscore", input: "_", want: nil, wantErr: true},
		{name: "int then letters", input: "123abc", want: nil, wantErr: true},
		{name: "int then slash", input: "123/foo", want: nil, wantErr: true},
		{name: "symbol then quote", input: `foo"bar"`, want: nil, wantErr: true},
		{name: "int then quote", input: `123"foo"`, want: nil, wantErr: true},
		{name: "trailing backslash", input: "\"foo\\", want: nil, wantErr: true},
	}

	runTokenCases(t, tests)
}

func TestTokenizeWhitespace(t *testing.T) {
	tests := []tokenCase{
		{name: "space only", input: " ", want: nil, wantErr: false},
		{name: "tab only", input: "\t", want: nil, wantErr: false},
		{name: "newline only", input: "\n", want: nil, wantErr: false},
		{name: "carriage return only", input: "\r", want: nil, wantErr: false},
		{name: "single comma", input: ",", want: nil, wantErr: false},
		{name: "multiple commas", input: ",,,", want: nil, wantErr: false},
		{name: "mixed whitespace only", input: " \t\n\r", want: nil, wantErr: false},
		{name: "mixed whitespace and commas only", input: " , ,\t\n\r, ", want: nil, wantErr: false},
		{name: "commas with newlines only", input: ",\n,\r\n,\t,", want: nil, wantErr: false},
		{name: "symbol padded spaces", input: "  foo  ", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "symbol padded tabs", input: "\tfoo\t", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "symbol padded newlines", input: "\nfoo\n", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "symbol padded carriage returns", input: "\rfoo\r", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "symbol padded commas", input: ",foo,", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "symbol padded mixed", input: " \t\n\r,foo, \r\n\t ", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "int padded mixed", input: " ,\t\n\r42\r\n\t, ", want: []Token{{Type: TOKEN_INT, Value: "42"}}, wantErr: false},
		{name: "symbols space separated", input: "foo bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "symbols tab separated", input: "foo\tbar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "symbols newline separated", input: "foo\nbar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "symbols cr separated", input: "foo\rbar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "symbols comma separated", input: "foo,bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "symbols multi comma separated", input: "foo,,,bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "symbols comma space mixed", input: "foo , bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "symbols mixed whitespace separated", input: "foo \t\n\r, bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "ints space separated", input: "1 2", want: []Token{{Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}}, wantErr: false},
		{name: "ints tab separated", input: "1\t2", want: []Token{{Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}}, wantErr: false},
		{name: "ints newline separated", input: "1\n2", want: []Token{{Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}}, wantErr: false},
		{name: "ints cr separated", input: "1\r2", want: []Token{{Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}}, wantErr: false},
		{name: "ints comma separated", input: "1,2,3", want: []Token{{Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}, {Type: TOKEN_INT, Value: "3"}}, wantErr: false},
		{name: "ints mixed whitespace separated", input: "1 \t\n\r, 2,\t3", want: []Token{{Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_INT, Value: "2"}, {Type: TOKEN_INT, Value: "3"}}, wantErr: false},
		{name: "int comma symbol", input: "123,foo", want: []Token{{Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "symbol comma int", input: "foo,123", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_INT, Value: "123"}}, wantErr: false},
		{name: "list comma separated", input: "(foo,123)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "list commas and spaces", input: "(,foo,,123,,)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "parens only commas", input: "(,)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "empty parens with whitespace", input: "( \t\n\r, )", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "leading commas before list", input: ",,(foo)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "trailing commas after list", input: "(foo),,", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "comma between parens", input: "(1),(2)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_INT, Value: "1"}, {Type: TOKEN_RPAREN, Value: ")"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_INT, Value: "2"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "crlf separated list", input: "(foo\r\nbar)", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "symbol comma lparen", input: "foo,(bar)", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "bar"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "string padded whitespace", input: "  \"hello\"  ", want: []Token{{Type: TOKEN_DQS, Value: `"hello"`}}, wantErr: false},
		{name: "string padded commas", input: ",\"hello\",", want: []Token{{Type: TOKEN_DQS, Value: `"hello"`}}, wantErr: false},
		{name: "string padded mixed", input: " \t\n\r,\"hello\", \r\n\t ", want: []Token{{Type: TOKEN_DQS, Value: `"hello"`}}, wantErr: false},
		{name: "string with comma inside", input: "\"a,b\"", want: []Token{{Type: TOKEN_DQS, Value: `"a,b"`}}, wantErr: false},
		{name: "string comma symbol", input: "\"foo\",bar", want: []Token{{Type: TOKEN_DQS, Value: `"foo"`}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "string space int", input: "\"foo\" 123", want: []Token{{Type: TOKEN_DQS, Value: `"foo"`}, {Type: TOKEN_INT, Value: "123"}}, wantErr: false},
		{name: "multiline mixed", input: "  ( \n foo\t,123 \r\n )  ", want: []Token{{Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "whitespace before invalid", input: "foo , à", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: true},
	}

	runTokenCases(t, tests)
}

func TestTokenIteratorEarlyBreak(t *testing.T) {
	seq, errFn := TokenIterator("(1 2 3)")
	var got []Token
	for tok := range seq {
		got = append(got, tok)
		break
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 token after early break, got %d tokens %v", len(got), got)
	}
	if want := (Token{Type: TOKEN_LPAREN, Value: "("}); got[0] != want {
		t.Errorf("first token: got %+v, want %+v", got[0], want)
	}
	if err := errFn(); err != nil {
		t.Errorf("expected nil error after early break, got %v", err)
	}
}
