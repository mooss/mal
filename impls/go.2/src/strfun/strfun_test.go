//revive:disable:cognitive-complexity

package strfun_test

import (
	"testing"

	"github.com/mooss/mal/src/strfun"
	"github.com/mooss/mal/src/tut"
	"github.com/mooss/mal/src/tut/assert"
)

///////////
// Utils //

func checkMalQuote(t *testing.T, unquoted, wantQuoted string) {
	t.Helper()
	got := strfun.MalQuote(unquoted)
	assert.Equal(t, got, wantQuoted)
	assert.StringStartsWith(t, got, `"`)
	assert.StringEndsWith(t, got, `"`)
}

func checkMalUnquote(t *testing.T, quoted, want string) {
	t.Helper()
	got, err := strfun.MalUnquote(quoted)
	assert.Equal(t, assert.Must(t, got, err), want)
}

func checkMalUnquoteError(t *testing.T, quoted, wantSub string) {
	t.Helper()
	_, err := strfun.MalUnquote(quoted)
	assert.ErrContains(t, err, wantSub)
}

//////////
// Data //

// Shared round-trip pairs: unquoted <-> quoted.
// Used by both TestMalQuote and TestMalUnquote so the two directions cannot drift apart.
// Input holds the unquoted value, Want holds the quoted value.
//
//revive:disable:line-length-limit
var sharedQuotePairs = []tut.Case[string]{
	{Name: "empty", Input: "", Want: `""`},
	{Name: "plain", Input: "hello", Want: `"hello"`},
	{Name: "spaces", Input: "hello world", Want: `"hello world"`},
	{Name: "digits", Input: "123", Want: `"123"`},
	{Name: "punctuation", Input: "a,b.c!?", Want: `"a,b.c!?"`},
	{Name: "single quote inside", Input: "it's", Want: `"it's"`},
	{Name: "tab preserved", Input: "a\tb", Want: "\"a\tb\""},
	{Name: "carriage return preserved", Input: "a\rb", Want: "\"a\rb\""},
	{Name: "unicode", Input: "héllo 🌍", Want: "\"héllo 🌍\""},

	{Name: "single newline", Input: "\n", Want: `"\n"`},
	{Name: "newline middle", Input: "a\nb", Want: `"a\nb"`},
	{Name: "newline start", Input: "\na", Want: `"\na"`},
	{Name: "newline end", Input: "a\n", Want: `"a\n"`},
	{Name: "multiple newlines", Input: "\n\n\n", Want: `"\n\n\n"`},
	{Name: "newline mixed text", Input: "line1\nline2\nline3", Want: `"line1\nline2\nline3"`},

	{Name: "single dquote", Input: `"`, Want: `"\""`},
	{Name: "multiple dquotes", Input: `"""`, Want: `"\"\"\""`},
	{Name: "dquote middle", Input: `say "hi"`, Want: `"say \"hi\""`},
	{Name: "dquote start end", Input: `"hi"`, Want: `"\"hi\""`},
	{Name: "dquote only ends", Input: `"a`, Want: `"\"a"`},

	{Name: "single backslash", Input: `\`, Want: `"\\"`},
	{Name: "double backslash", Input: `\\`, Want: `"\\\\"`},
	{Name: "triple backslash", Input: `\\\`, Want: `"\\\\\\"`},
	{Name: "backslash middle", Input: `a\b`, Want: `"a\\b"`},
	{Name: "backslash start", Input: `\a`, Want: `"\\a"`},
	{Name: "backslash end", Input: `a\`, Want: `"a\\"`},

	// backslash + 'n' (two chars) must not be confused with newline (one char)
	{Name: "backslash plus n", Input: "\\n", Want: `"\\n"`},
	{Name: "backslash plus n middle", Input: "a\\nb", Want: `"a\\nb"`},

	{Name: "newline plus backslash", Input: "a\n\\b", Want: `"a\n\\b"`},
	{Name: "backslash plus quote", Input: "\\\"", Want: `"\\\""`},
	{Name: "quote plus backslash", Input: `"\`, Want: `"\"\\"`},
	{Name: "quote plus two backslashes", Input: `"\\`, Want: `"\"\\\\"`},
	{Name: "all specials", Input: "\"\\\n", Want: `"\"\\\n"`},
	{Name: "all specials reversed", Input: "\n\\\"", Want: `"\n\\\""`},
	{Name: "sentence with specials", Input: "a\"b\\c\nd\"e\\f", Want: `"a\"b\\c\nd\"e\\f"`},
	{Name: "many specials", Input: "\\\"\\\n\"\\", Want: `"\\\"\\\n\"\\"`},
	{Name: "adjacent escapes", Input: "\n\n\"\"\\\\", Want: `"\n\n\"\"\\\\"`},
}

var unquoteErrors = []tut.Case[string]{
	{Name: "invalid q middle", Input: `"a\qb"`, Want: `invalid escape \q`, WantErr: true},
	{Name: "invalid t middle", Input: `"a\tb"`, Want: `invalid escape \t`, WantErr: true},
	{Name: "invalid r", Input: `"a\rb"`, Want: `invalid escape \r`, WantErr: true},
	{Name: "invalid b", Input: `"a\bb"`, Want: `invalid escape \b`, WantErr: true},
	{Name: "invalid f", Input: `"a\fb"`, Want: `invalid escape \f`, WantErr: true},
	{Name: "invalid v", Input: `"a\vb"`, Want: `invalid escape \v`, WantErr: true},
	{Name: "invalid 0", Input: `"a\0b"`, Want: `invalid escape \0`, WantErr: true},
	{Name: "invalid u", Input: `"a\ub"`, Want: `invalid escape \u`, WantErr: true},
	{Name: "invalid x", Input: `"a\xb"`, Want: `invalid escape \x`, WantErr: true},
	{Name: "invalid space", Input: `"a\ b"`, Want: `invalid escape \ `, WantErr: true},
	{Name: "invalid single quote", Input: `"a\'b"`, Want: `invalid escape \'`, WantErr: true},
	{Name: "invalid slash", Input: `"a\/b"`, Want: `invalid escape \/`, WantErr: true},
	{Name: "invalid a", Input: `"a\ab"`, Want: `invalid escape \a`, WantErr: true},
	{Name: "invalid e", Input: `"a\eb"`, Want: `invalid escape \e`, WantErr: true},
	{Name: "invalid N uppercase", Input: `"a\Nb"`, Want: `invalid escape \N`, WantErr: true},
	{Name: "invalid at start", Input: `"\qb"`, Want: `invalid escape \q`, WantErr: true},
	{Name: "invalid at end", Input: `"ab\q"`, Want: `invalid escape \q`, WantErr: true},
	{Name: "only invalid", Input: `"\q"`, Want: `invalid escape \q`, WantErr: true},
	{Name: "double invalid first wins", Input: `"\q\n"`, Want: `invalid escape \q`, WantErr: true},
	{Name: "invalid after valid", Input: `"\\n\q"`, Want: `invalid escape \q`, WantErr: true},
	{Name: "invalid digit", Input: `"a\1b"`, Want: `invalid escape \1`, WantErr: true},
	{Name: "invalid star", Input: `"a\*b"`, Want: `invalid escape \*`, WantErr: true},
	{Name: "empty", Input: "", Want: `missing surrounding quotes`, WantErr: true},
	{Name: "single char", Input: "a", Want: `missing surrounding quotes`, WantErr: true},
	{Name: "single dquote", Input: "\"", Want: `missing surrounding quotes`, WantErr: true},
	{Name: "missing leading quote", Input: `abc"`, Want: `missing surrounding quotes`, WantErr: true},
	{Name: "missing trailing quote", Input: `"abc`, Want: `missing surrounding quotes`, WantErr: true},
	{Name: "unquoted", Input: `abc`, Want: `missing surrounding quotes`, WantErr: true},
	{Name: "single quotes", Input: `'abc'`, Want: `missing surrounding quotes`, WantErr: true},
	{Name: "trailing lone backslash", Input: `"abc\"`, Want: `invalid escape`, WantErr: true},
	{Name: "only lone backslash", Input: `"\"`, Want: `invalid escape`, WantErr: true},
}

///////////
// Tests //

func TestMalQuote(t *testing.T) {
	t.Run("shared pairs", func(t *testing.T) {
		tut.Run(t, sharedQuotePairs, func(t *testing.T, c tut.Case[string]) {
			checkMalQuote(t, c.Input, c.Want)
		})
	})

	t.Run("extra cases", func(t *testing.T) {
		extra := []tut.Case[string]{
			{Name: "long", Input: "the quick brown fox jumps over the lazy dog", Want: `"the quick brown fox jumps over the lazy dog"`},
			{Name: "only spaces", Input: "   ", Want: `"   "`},
			{Name: "semicolon", Input: "(+ 1 2)", Want: `"(+ 1 2)"`},
		}
		tut.Run(t, extra, func(t *testing.T, c tut.Case[string]) {
			checkMalQuote(t, c.Input, c.Want)
		})
	})
}

func TestMalUnquote(t *testing.T) {
	t.Run("shared pairs", func(t *testing.T) {
		tut.Run(t, sharedQuotePairs, func(t *testing.T, c tut.Case[string]) {
			checkMalUnquote(t, c.Want, c.Input)
		})
	})

	t.Run("extra valid cases", func(t *testing.T) {
		//TODO: consider failing on those cases.
		extra := []tut.Case[string]{
			{Name: "escaped quote at end", Input: `a"`, Want: `"a\""`},
			{Name: "escaped backslash at end", Input: `a\`, Want: `"a\\"`},
			{Name: "escaped newline at end", Input: "a\n", Want: `"a\n"`},
			{Name: "only escaped quote", Input: `"`, Want: `"\""`},
			{Name: "only escaped backslash", Input: `\`, Want: `"\\"`},
			{Name: "only escaped newline", Input: "\n", Want: `"\n"`},
			{Name: "consecutive escapes", Input: "\\\"\n", Want: `"\\\"\n"`},
			{Name: "text around escapes", Input: "x\\y\"z\nw", Want: `"x\\y\"z\nw"`},
		}
		tut.Run(t, extra, func(t *testing.T, c tut.Case[string]) {
			checkMalUnquote(t, c.Want, c.Input)
		})
	})
}

func TestMalUnquoteErrors(t *testing.T) {
	tut.Run(t, unquoteErrors, func(t *testing.T, c tut.Case[string]) {
		checkMalUnquoteError(t, c.Input, c.Want)
	})
}

func TestMalQuoteUnquoteRoundTrip(t *testing.T) {
	t.Run("unquoted round trip", func(t *testing.T) {
		cases := []string{
			"", "hello", `"`, `\`, "\n", `a"b\c` + "\n" + `d`, "\\n", "\n\\\"\\\\", "héllo 🌍",
			"a\tb", "   ", "\n\n\"\"\\\\",
		}
		// include every shared unquoted value as well.
		for _, p := range sharedQuotePairs {
			cases = append(cases, p.Input)
		}
		for _, s := range cases {
			got, err := strfun.MalUnquote(strfun.MalQuote(s))
			assert.Equal(t, assert.Must(t, got, err), s)
		}
	})

	t.Run("quoted round trip", func(t *testing.T) {
		tut.Run(t, sharedQuotePairs, func(t *testing.T, c tut.Case[string]) {
			got, err := strfun.MalUnquote(c.Want)
			assert.Equal(t, strfun.MalQuote(assert.Must(t, got, err)), c.Want)
		})
	})
}
