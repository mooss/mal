//revive:disable:cognitive-complexity,line-length-limit,function-length
package lexer

import "testing"

func TestTokenizeStrings(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantTokens []string
		wantErr    bool
	}{
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
		{name: "raw newline with closing", input: "\"a\nb\"", wantTokens: []string{}, wantErr: true},
		{name: "CR with closing", input: "\"a\rb\"", wantTokens: []string{"\"a\rb\""}, wantErr: false},
		{name: "invalid utf8 with closing", input: "\"a\xffb\"", wantTokens: []string{"\"a�b\""}, wantErr: false},
		{name: "NUL no closing", input: "\"foo\x00", wantTokens: []string{}, wantErr: true},
		{name: "invalid no closing", input: "\"foo\xff", wantTokens: []string{}, wantErr: true},
		{name: "two adjacent", input: `"a""b"`, wantTokens: []string{`"a"`, `"b"`}, wantErr: false},
		{name: "trailing escaped backslash", input: `"foo\\"`, wantTokens: []string{`"foo\\"`}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, errFn := TokenIterator(tt.input)
			var gotTokens []string
			for tok := range tokens {
				if tok.Type != TOKEN_DQS {
					t.Errorf("input %q: got token type %v (%q), want %v", tt.input, tok.Type, tok.Value, TOKEN_DQS)
				}
				gotTokens = append(gotTokens, tok.Value)
			}
			gotErr := errFn() != nil
			if len(gotTokens) != len(tt.wantTokens) {
				t.Fatalf("input %q: got %d tokens %q (err %v), want %d tokens %q", tt.input, len(gotTokens), gotTokens, gotErr, len(tt.wantTokens), tt.wantTokens)
			}
			for i := range tt.wantTokens {
				if gotTokens[i] != tt.wantTokens[i] {
					t.Errorf("input %q token %d: got %q, want %q", tt.input, i, gotTokens[i], tt.wantTokens[i])
				}
			}
			if gotErr != tt.wantErr {
				if tt.wantErr {
					t.Errorf("input %q tokens %q: expected error, got nil", tt.input, gotTokens)
				} else {
					t.Errorf("input %q tokens %q: expected no error, got error", tt.input, gotTokens)
				}
			}
		})
	}
}

func testSingleCharToken(t *testing.T, char string, typ TokenType) {
	t.Helper()
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
		{name: "single", input: char, want: []Token{{Type: typ, Value: char}}, wantErr: false},
		{name: "doubled", input: char + char, want: []Token{{Type: typ, Value: char}, {Type: typ, Value: char}}, wantErr: false},
		{name: "spaced", input: char + " " + char, want: []Token{{Type: typ, Value: char}, {Type: typ, Value: char}}, wantErr: false},
		{name: "whitespace padded", input: "  " + char + "\n\t ", want: []Token{{Type: typ, Value: char}}, wantErr: false},
		{name: "symbol before", input: "foo" + char, want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: typ, Value: char}}, wantErr: false},
		{name: "symbol after", input: char + "foo", want: []Token{{Type: typ, Value: char}, {Type: TOKEN_SYMBOL, Value: "foo"}}, wantErr: false},
		{name: "int before", input: "123" + char, want: []Token{{Type: TOKEN_INT, Value: "123"}, {Type: typ, Value: char}}, wantErr: false},
		{name: "int after", input: char + "123", want: []Token{{Type: typ, Value: char}, {Type: TOKEN_INT, Value: "123"}}, wantErr: false},
		{name: "symbol spaced", input: "foo " + char + " bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: typ, Value: char}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, errFn := TokenIterator(tt.input)
			var got []Token
			for tok := range seq {
				got = append(got, tok)
			}
			gotErr := errFn() != nil
			if len(got) != len(tt.want) {
				t.Fatalf("input %q: got %d tokens %v (err %v), want %d tokens %v", tt.input, len(got), got, gotErr, len(tt.want), tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("input %q token %d: got %+v, want %+v", tt.input, i, got[i], tt.want[i])
				}
			}
			if gotErr != tt.wantErr {
				t.Errorf("input %q: gotErr %v, wantErr %v", tt.input, gotErr, tt.wantErr)
			}
		})
	}
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
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, errFn := TokenIterator(tt.input)
			var got []Token
			for tok := range seq {
				got = append(got, tok)
			}
			gotErr := errFn() != nil
			if len(got) != len(tt.want) {
				t.Fatalf("input %q: got %d tokens %v (err %v), want %d tokens %v", tt.input, len(got), got, gotErr, len(tt.want), tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("input %q token %d: got %+v, want %+v", tt.input, i, got[i], tt.want[i])
				}
			}
			if gotErr != tt.wantErr {
				t.Errorf("input %q: gotErr %v, wantErr %v", tt.input, gotErr, tt.wantErr)
			}
		})
	}
}

func TestTokenizeSymbols(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
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
		{name: "punct run", input: "*+-./", want: []Token{{Type: TOKEN_SYMBOL, Value: "*+-./"}}, wantErr: false},
		{name: "multiple", input: "foo bar baz", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_SYMBOL, Value: "bar"}, {Type: TOKEN_SYMBOL, Value: "baz"}}, wantErr: false},
		{name: "all plain punct", input: "a/b.c*d+e-f", want: []Token{{Type: TOKEN_SYMBOL, Value: "a/b.c*d+e-f"}}, wantErr: false},
		{name: "keywords", input: "nil true false defn", want: []Token{{Type: TOKEN_SYMBOL, Value: "nil"}, {Type: TOKEN_SYMBOL, Value: "true"}, {Type: TOKEN_SYMBOL, Value: "false"}, {Type: TOKEN_SYMBOL, Value: "defn"}}, wantErr: false},
		{name: "digit minus", input: "1-", want: []Token{}, wantErr: true},
		{name: "digits plus digits", input: "12+34", want: []Token{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, errFn := TokenIterator(tt.input)
			var got []Token
			for tok := range seq {
				got = append(got, tok)
			}
			gotErr := errFn() != nil
			if len(got) != len(tt.want) {
				t.Fatalf("input %q: got %d tokens %v (err %v), want %d tokens %v", tt.input, len(got), got, gotErr, len(tt.want), tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("input %q token %d: got %+v, want %+v", tt.input, i, got[i], tt.want[i])
				}
			}
			if gotErr != tt.wantErr {
				t.Errorf("input %q: gotErr %v, wantErr %v", tt.input, gotErr, tt.wantErr)
			}
		})
	}
}

func TestTokenizeCombinations(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, errFn := TokenIterator(tt.input)
			var got []Token
			for tok := range seq {
				got = append(got, tok)
			}
			gotErr := errFn() != nil
			if len(got) != len(tt.want) {
				t.Fatalf("input %q: got %d tokens %v (err %v), want %d tokens %v", tt.input, len(got), got, gotErr, len(tt.want), tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("input %q token %d: got %+v, want %+v", tt.input, i, got[i], tt.want[i])
				}
			}
			if gotErr != tt.wantErr {
				t.Errorf("input %q: gotErr %v, wantErr %v", tt.input, gotErr, tt.wantErr)
			}
		})
	}
}

func TestTokenizeInvalid(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
		{name: "backtick", input: "`", want: nil, wantErr: true},
		{name: "tilde", input: "~", want: nil, wantErr: true},
		{name: "bang", input: "!", want: nil, wantErr: true},
		{name: "hash", input: "#", want: nil, wantErr: true},
		{name: "dollar", input: "$", want: nil, wantErr: true},
		{name: "percent", input: "%", want: nil, wantErr: true},
		{name: "ampersand", input: "&", want: nil, wantErr: true},
		{name: "pipe", input: "|", want: nil, wantErr: true},
		{name: "backslash", input: `\`, want: nil, wantErr: true},
		{name: "caret", input: "^", want: nil, wantErr: true},
		{name: "underscore", input: "_", want: nil, wantErr: true},
		{name: "int then letters", input: "123abc", want: nil, wantErr: true},
		{name: "int then slash", input: "123/foo", want: nil, wantErr: true},
		{name: "symbol then bang", input: "foo!bar", want: nil, wantErr: true},
		{name: "symbol then quote", input: `foo"bar"`, want: nil, wantErr: true},
		{name: "int then quote", input: `123"foo"`, want: nil, wantErr: true},
		{name: "symbol boundary lparen", input: "foo(bar)", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "bar"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "symbol boundary rparen", input: "foo)bar", want: []Token{{Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_RPAREN, Value: ")"}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "int boundary lparen", input: "123(foo)", want: []Token{{Type: TOKEN_INT, Value: "123"}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "foo"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "string adjacent symbol", input: `"foo"bar`, want: []Token{{Type: TOKEN_DQS, Value: `"foo"`}, {Type: TOKEN_SYMBOL, Value: "bar"}}, wantErr: false},
		{name: "string adjacent int", input: `"foo"123`, want: []Token{{Type: TOKEN_DQS, Value: `"foo"`}, {Type: TOKEN_INT, Value: "123"}}, wantErr: false},
		{name: "string adjacent lparen", input: `"foo"(bar)`, want: []Token{{Type: TOKEN_DQS, Value: `"foo"`}, {Type: TOKEN_LPAREN, Value: "("}, {Type: TOKEN_SYMBOL, Value: "bar"}, {Type: TOKEN_RPAREN, Value: ")"}}, wantErr: false},
		{name: "escaped newline", input: "\"foo\\\nbar\"", want: nil, wantErr: true},
		{name: "trailing backslash", input: "\"foo\\", want: nil, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, errFn := TokenIterator(tt.input)
			var got []Token
			for tok := range seq {
				got = append(got, tok)
			}
			gotErr := errFn() != nil
			if len(got) != len(tt.want) {
				t.Fatalf("input %q: got %d tokens %v (err %v), want %d tokens %v", tt.input, len(got), got, gotErr, len(tt.want), tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("input %q token %d: got %+v, want %+v", tt.input, i, got[i], tt.want[i])
				}
			}
			if gotErr != tt.wantErr {
				t.Errorf("input %q: gotErr %v, wantErr %v", tt.input, gotErr, tt.wantErr)
			}
		})
	}
}

func TestTokenizeWhitespace(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, errFn := TokenIterator(tt.input)
			var got []Token
			for tok := range seq {
				got = append(got, tok)
			}
			gotErr := errFn() != nil
			if len(got) != len(tt.want) {
				t.Fatalf("input %q: got %d tokens %v (err %v), want %d tokens %v", tt.input, len(got), got, gotErr, len(tt.want), tt.want)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("input %q token %d: got %+v, want %+v", tt.input, i, got[i], tt.want[i])
				}
			}
			if gotErr != tt.wantErr {
				t.Errorf("input %q: gotErr %v, wantErr %v", tt.input, gotErr, tt.wantErr)
			}
		})
	}
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
