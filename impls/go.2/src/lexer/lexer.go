package lexer

import (
	"fmt"
	"iter"
	"unicode/utf8"
)

// LexicalFailure describes what caused the lexer to fail.
// It can be followed by additional information specified after `: `.
type LexicalFailure string

const (
	UnterminatedDQS LexicalFailure = "unterminated double quoted string"
	InvalidStart    LexicalFailure = "invalid token start"
)

const eof = -1

type lexer struct {
	input string // The string being lexed.
	start int    // Start position of the current token.
	pos   int    // Current position in the input.
	width int    // Width of the last rune read.

	consumer func(Token) bool // Token consumer.
	err      error            // Error that happened during lexing.
}

///////////////////////
// Iteration methods //

// over returns true when all the input has been lexed.
func (lex *lexer) over() bool {
	return lex.pos >= len(lex.input)
}

// next advances to the next rune, returning the current one.
func (lex *lexer) next() rune {
	if lex.pos >= len(lex.input) {
		return eof
	}

	var glyph rune
	glyph, lex.width = utf8.DecodeRuneInString(lex.input[lex.pos:])
	lex.pos += lex.width
	return glyph
}

// produce returns the current token and advances the input.
func (lex *lexer) produce(typ TokenType) Token {
	token := lex.input[lex.start:lex.pos]
	lex.start = lex.pos
	return Token{Type: typ, Value: token}
}

/////////////////////////////////
// Communication with consumer //

// feed yields the token to the consumer and advances the input to the end if the consumer signals
// to stop the production.
func (lex *lexer) feed(token Token) {
	if !lex.consumer(token) {
		lex.pos = len(lex.input)
	}
}

// failure ends the iteration and sets the error.
func (lex *lexer) failure(reason LexicalFailure) stateFun {
	token := lex.produce(TOKEN_ERROR)
	lex.err = fmt.Errorf("%s: %#q", reason, token.Value)
	lex.pos = len(lex.input) // Iteration cannot continue after error.
	return nil
}

// success yields the current token without error and advances the input.
func (lex *lexer) success(typ TokenType) {
	lex.feed(lex.produce(typ))
}

//////////////////
// Lexer states //

// stateFun represents the state of the lexer as a function returning the next state.
type stateFun func(*lexer) stateFun

// Tokenize turns input into an iterable sequence over its tokens and an error getter.
// The error should always be checked after iterating.
func Tokenize(input string) (iter.Seq[Token], func() error) {
	lex := &lexer{input: input}
	return func(yield func(Token) bool) {
		lex.consumer = yield
		for state := lexMain; state != nil; {
			state = state(lex)
		}
	}, func() error { return lex.err }
}

// lexDqs seeks the end of a double quoted string.
// The lexer is supposed to be placed at the opening double quote.
func lexDqs(lex *lexer) stateFun {
	lex.pos++
	for {
		switch lex.next() {
		case '"':
			lex.success(TOKEN_DQS)
			return lexMain
		case '\\':
			if glyph := lex.next(); glyph != eof && glyph != '\n' {
				break // The escaped character is "normal".
			}
			fallthrough
		case eof, '\n':
			return lex.failure(UnterminatedDQS)
		}
	}
}

// lexMain is the entry point of the lexer.
func lexMain(lex *lexer) stateFun {
	for !lex.over() {
		switch lex.input[lex.pos] {
		case '"':
			return lexDqs(lex)
		default:
			lex.next() // Advance the input to show the rune in the error message.
			lex.failure(InvalidStart)
			return nil
		}
	}
	return nil
}
