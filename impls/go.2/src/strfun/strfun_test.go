//revive:disable:cognitive-complexity

package strfun

import (
	"strings"
	"testing"
)

///////////
// Utils //

func checkMalQuote(t *testing.T, unquoted, wantQuoted string) {
	t.Helper()
	got := MalQuote(unquoted)
	if got != wantQuoted {
		t.Errorf("MalQuote(%q) = %q, want %q", unquoted, got, wantQuoted)
	}
	// structural invariants shared by all quoted strings
	if !strings.HasPrefix(got, `"`) || !strings.HasSuffix(got, `"`) || len(got) < 2 {
		t.Errorf("MalQuote(%q) = %q: must be wrapped in double quotes", unquoted, got)
	}
}

func checkMalUnquote(t *testing.T, quoted, want string) {
	t.Helper()
	got, err := MalUnquote(quoted)
	if err != nil {
		t.Fatalf("MalUnquote(%q) unexpected error: %v", quoted, err)
	}
	if got != want {
		t.Errorf("MalUnquote(%q) = %q, want %q", quoted, got, want)
	}
}

func checkMalUnquoteError(t *testing.T, quoted, wantSub string) {
	t.Helper()
	got, err := MalUnquote(quoted)
	if err == nil {
		t.Fatalf("MalUnquote(%q) = %q, want error containing %q", quoted, got, wantSub)
	}
	if wantSub != "" && !strings.Contains(err.Error(), wantSub) {
		t.Errorf("MalUnquote(%q) error = %q, want substring %q", quoted, err.Error(), wantSub)
	}
}

func mustUnquote(t *testing.T, quoted string) string {
	t.Helper()
	s, err := MalUnquote(quoted)
	if err != nil {
		t.Fatalf("mustUnquote(%q) unexpected error: %v", quoted, err)
	}
	return s
}

//////////
// Data //

// Shared round-trip pairs: unquoted <-> quoted.
// Used by both TestMalQuote and TestMalUnquote so the two directions cannot drift apart.
type quotePair struct {
	name     string
	unquoted string
	quoted   string
}

//revive:disable:line-length-limit
var sharedQuotePairs = []quotePair{
	{"empty", "", `""`},
	{"plain", "hello", `"hello"`},
	{"spaces", "hello world", `"hello world"`},
	{"digits", "123", `"123"`},
	{"punctuation", "a,b.c!?", `"a,b.c!?"`},
	{"single quote inside", "it's", `"it's"`},
	{"tab preserved", "a\tb", "\"a\tb\""},
	{"carriage return preserved", "a\rb", "\"a\rb\""},
	{"unicode", "héllo 🌍", "\"héllo 🌍\""},

	{"single newline", "\n", `"\n"`},
	{"newline middle", "a\nb", `"a\nb"`},
	{"newline start", "\na", `"\na"`},
	{"newline end", "a\n", `"a\n"`},
	{"multiple newlines", "\n\n\n", `"\n\n\n"`},
	{"newline mixed text", "line1\nline2\nline3", `"line1\nline2\nline3"`},

	{"single dquote", `"`, `"\""`},
	{"multiple dquotes", `"""`, `"\"\"\""`},
	{"dquote middle", `say "hi"`, `"say \"hi\""`},
	{"dquote start end", `"hi"`, `"\"hi\""`},
	{"dquote only ends", `"a`, `"\"a"`},

	{"single backslash", `\`, `"\\"`},
	{"double backslash", `\\`, `"\\\\"`},
	{"triple backslash", `\\\`, `"\\\\\\"`},
	{"backslash middle", `a\b`, `"a\\b"`},
	{"backslash start", `\a`, `"\\a"`},
	{"backslash end", `a\`, `"a\\"`},

	// backslash + 'n' (two chars) must not be confused with newline (one char)
	{"backslash plus n", "\\n", `"\\n"`},
	{"backslash plus n middle", "a\\nb", `"a\\nb"`},

	{"newline plus backslash", "a\n\\b", `"a\n\\b"`},
	{"backslash plus quote", "\\\"", `"\\\""`},
	{"quote plus backslash", `"\`, `"\"\\"`},
	{"quote plus two backslashes", `"\\`, `"\"\\\\"`},
	{"all specials", "\"\\\n", `"\"\\\n"`},
	{"all specials reversed", "\n\\\"", `"\n\\\""`},
	{"sentence with specials", "a\"b\\c\nd\"e\\f", `"a\"b\\c\nd\"e\\f"`},
	{"many specials", "\\\"\\\n\"\\", `"\\\"\\\n\"\\"`},
	{"adjacent escapes", "\n\n\"\"\\\\", `"\n\n\"\"\\\\"`},
}

type unquoteError struct {
	name    string
	quoted  string
	wantSub string
}

var unquoteErrors = []unquoteError{
	{"invalid q middle", `"a\qb"`, `invalid escape \q`},
	{"invalid t middle", `"a\tb"`, `invalid escape \t`},
	{"invalid r", `"a\rb"`, `invalid escape \r`},
	{"invalid b", `"a\bb"`, `invalid escape \b`},
	{"invalid f", `"a\fb"`, `invalid escape \f`},
	{"invalid v", `"a\vb"`, `invalid escape \v`},
	{"invalid 0", `"a\0b"`, `invalid escape \0`},
	{"invalid u", `"a\ub"`, `invalid escape \u`},
	{"invalid x", `"a\xb"`, `invalid escape \x`},
	{"invalid space", `"a\ b"`, `invalid escape \ `},
	{"invalid single quote", `"a\'b"`, `invalid escape \'`},
	{"invalid slash", `"a\/b"`, `invalid escape \/`},
	{"invalid a", `"a\ab"`, `invalid escape \a`},
	{"invalid e", `"a\eb"`, `invalid escape \e`},
	{"invalid N uppercase", `"a\Nb"`, `invalid escape \N`},
	{"invalid at start", `"\qb"`, `invalid escape \q`},
	{"invalid at end", `"ab\q"`, `invalid escape \q`},
	{"only invalid", `"\q"`, `invalid escape \q`},
	{"double invalid first wins", `"\q\n"`, `invalid escape \q`},
	{"invalid after valid", `"\\n\q"`, `invalid escape \q`},
	{"invalid digit", `"a\1b"`, `invalid escape \1`},
	{"invalid star", `"a\*b"`, `invalid escape \*`},
	{"empty", "", `missing surrounding quotes`},
	{"single char", "a", `missing surrounding quotes`},
	{"single dquote", "\"", `missing surrounding quotes`},
	{"missing leading quote", `abc"`, `missing surrounding quotes`},
	{"missing trailing quote", `"abc`, `missing surrounding quotes`},
	{"unquoted", `abc`, `missing surrounding quotes`},
	{"single quotes", `'abc'`, `missing surrounding quotes`},
	{"trailing lone backslash", `"abc\"`, `invalid escape`},
	{"only lone backslash", `"\"`, `invalid escape`},
}

///////////
// Tests //

func TestMalQuote(t *testing.T) {
	t.Run("shared pairs", func(t *testing.T) {
		for _, tc := range sharedQuotePairs {
			t.Run(tc.name, func(t *testing.T) {
				checkMalQuote(t, tc.unquoted, tc.quoted)
			})
		}
	})

	t.Run("extra cases", func(t *testing.T) {
		extra := []quotePair{
			{"long", "the quick brown fox jumps over the lazy dog", `"the quick brown fox jumps over the lazy dog"`},
			{"only spaces", "   ", `"   "`},
			{"semicolon", "(+ 1 2)", `"(+ 1 2)"`},
		}
		for _, tc := range extra {
			t.Run(tc.name, func(t *testing.T) {
				checkMalQuote(t, tc.unquoted, tc.quoted)
			})
		}
	})
}

func TestMalUnquote(t *testing.T) {
	t.Run("shared pairs", func(t *testing.T) {
		for _, tc := range sharedQuotePairs {
			t.Run(tc.name, func(t *testing.T) {
				checkMalUnquote(t, tc.quoted, tc.unquoted)
			})
		}
	})

	t.Run("extra valid cases", func(t *testing.T) {
		//TODO: consider failing on those cases.
		extra := []quotePair{
			{"escaped quote at end", `a"`, `"a\""`},
			{"escaped backslash at end", `a\`, `"a\\"`},
			{"escaped newline at end", "a\n", `"a\n"`},
			{"only escaped quote", `"`, `"\""`},
			{"only escaped backslash", `\`, `"\\"`},
			{"only escaped newline", "\n", `"\n"`},
			{"consecutive escapes", "\\\"\n", `"\\\"\n"`},
			{"text around escapes", "x\\y\"z\nw", `"x\\y\"z\nw"`},
		}
		for _, tc := range extra {
			t.Run(tc.name, func(t *testing.T) {
				checkMalUnquote(t, tc.quoted, tc.unquoted)
			})
		}
	})
}

func TestMalUnquoteErrors(t *testing.T) {
	for _, tc := range unquoteErrors {
		t.Run(tc.name, func(t *testing.T) {
			checkMalUnquoteError(t, tc.quoted, tc.wantSub)
		})
	}
}

func TestMalQuoteUnquoteRoundTrip(t *testing.T) {
	t.Run("unquoted round trip", func(t *testing.T) {
		cases := []string{
			"", "hello", `"`, `\`, "\n", `a"b\c` + "\n" + `d`, "\\n", "\n\\\"\\\\", "héllo 🌍",
			"a\tb", "   ", "\n\n\"\"\\\\",
		}
		// include every shared unquoted value as well.
		for _, p := range sharedQuotePairs {
			cases = append(cases, p.unquoted)
		}
		for _, s := range cases {
			quoted := MalQuote(s)
			got, err := MalUnquote(quoted)
			if err != nil {
				t.Errorf("MalUnquote(MalQuote(%q)=%q) unexpected error: %v", s, quoted, err)
				continue
			}
			if got != s {
				t.Errorf("MalUnquote(MalQuote(%q)=%q) = %q, want %q", s, quoted, got, s)
			}
		}
	})

	t.Run("quoted round trip", func(t *testing.T) {
		for _, p := range sharedQuotePairs {
			got := MalQuote(mustUnquote(t, p.quoted))
			if got != p.quoted {
				t.Errorf("MalQuote(MalUnquote(%q)) = %q, want %q", p.quoted, got, p.quoted)
			}
		}
	})
}
