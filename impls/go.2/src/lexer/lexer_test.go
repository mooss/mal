//revive:disable:cognitive-complexity,line-length-limit
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
		{name: "invalid utf8 with closing", input: "\"a\xffb\"", wantTokens: []string{"\"a\xffb\""}, wantErr: false},
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
