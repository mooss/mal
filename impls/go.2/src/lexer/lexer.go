package lexer

import (
	"errors"
	"fmt"
	"iter"
	"strings"
)

// TokenIterator turns input into an iterable sequence over its tokens and an error getter.
// The error should always be checked after iterating.
//
//revive:disable:cognitive-complexity
func TokenIterator(input string) (iter.Seq[Token], func() error) {
	lex := &Lexer{Input: []rune(input)}
	var err error
	return func(yield func(Token) bool) {
		for {
			token, scanErr := lexMal(lex)
			if scanErr != nil {
				if scanErr == EndOfInput { //nolint:errorlint
					return
				}
				err = scanErr
				return
			}
			if !yield(token) {
				return
			}
		}
	}, func() error { return err }
}

var EndOfInput = errors.New("end of input") //nolint:revive,staticcheck

// LexicalFailure describes what caused the lexer to fail.
// It can be followed by additional information specified after `: `.
type LexicalFailure string

const (
	UnbalancedDoubleQuote LexicalFailure = "unbalanced double quote"
	InvalidStart          LexicalFailure = "invalid token start"
	InvalidTokenBoundary  LexicalFailure = "invalid token boundary"
)

const EOF = -1

type Lexer struct {
	Input []rune // The runes being lexed.
	start int    // Start position of the current token.
	pos   int    // Current position in the input.
}

///////////////////////
// Iteration methods //

// Over returns true when all the input has been lexed.
func (lex *Lexer) Over() bool {
	return lex.pos >= len(lex.Input)
}

// Peek returns the rune at index+lex.pos without advancing the input.
// Always returns EOF when out of bounds.
// Index must not be negative.
func (lex *Lexer) Peek(index int) rune {
	index += lex.pos
	if index >= len(lex.Input) {
		return EOF
	}
	return lex.Input[index]
}

// Next advances to the next rune, returning the current one.
func (lex *Lexer) Next() rune {
	if lex.pos >= len(lex.Input) {
		return EOF
	}

	glyph := lex.Input[lex.pos]
	lex.pos++
	return glyph
}

// produce returns the current token and advances the input.
func (lex *Lexer) produce(typ TokenType) Token {
	token := string(lex.Input[lex.start:lex.pos])
	lex.start = lex.pos
	return Token{Type: typ, Value: token}
}

// Failuref formats a lexical error wrapping the current token.
func (lex *Lexer) Failuref(reason LexicalFailure, msg string, args ...any) error {
	token := lex.produce(TOKEN_ERROR)
	args = append([]any{reason, token.Value}, args...)
	return fmt.Errorf("%s at %#q: "+msg, args...)
}

// Failure returns a lexical error wrapping the current token.
func (lex *Lexer) Failure(reason LexicalFailure) error {
	token := lex.produce(TOKEN_ERROR)
	return fmt.Errorf("%s at %#q", reason, token.Value)
}

// Success returns the current token without error and advances the input.
func (lex *Lexer) Success(typ TokenType) (Token, error) {
	return lex.produce(typ), nil
}

// SuccessOnBoundary returns the current token without error if the next rune is a valid token
// boundary.
// Otherwise returns a special error identifying the invalid boundary.
func (lex *Lexer) SuccessOnBoundary(typ TokenType) (Token, error) {
	peek := lex.Peek(0)
	switch peek {
	case '(', ')', '[', ']', '{', '}', '\'', EOF, ' ', '\t', '\n', '\r', ',':
	default:
		lex.Next() // Make the faulty rune a part of the token included in the error message.
		return Token{}, lex.Failuref(InvalidTokenBoundary, "%q is not a boundary character", peek)
	}
	return lex.Success(typ)
}

// SkipWhitespace advances the current position and the start position to the next non-whitespace
// character.
func (lex *Lexer) SkipWhitespace() {
	for {
		switch lex.Peek(0) {
		// mal treats comma as whitespace.
		case ' ', '\t', '\n', '\r', ',':
			lex.Next()
			lex.start = lex.pos
		default:
			return
		}
	}
}

////////////
// Lexers //

// lexMal extracts the next Mal token.
// It returns the token on success, EndOfInput on end-of-input, or another error on lexical failure.
//
//revive:disable:cyclomatic
func lexMal(lex *Lexer) (Token, error) {
	lex.SkipWhitespace()
	next := lex.Next()
	switch next {
	case EOF:
		return Token{}, EndOfInput

	// Balanced expressions.
	case '(':
		return lex.Success(TOKEN_LPAREN)
	case ')':
		return lex.Success(TOKEN_RPAREN)
	case '[':
		return lex.Success(TOKEN_LBRACKET)
	case ']':
		return lex.Success(TOKEN_RBRACKET)
	case '{':
		return lex.Success(TOKEN_LBRACE)
	case '}':
		return lex.Success(TOKEN_RBRACE)
	case '"':
		return lexDqs(lex)

	// Atoms.
	case '-', '+':
		if isNumASCII(lex.Peek(0)) {
			return lexInt(lex, next)
		}
	case '\'':
		return lex.Success(TOKEN_SQUOTE)
	default:
	}

	if isNumASCII(next) {
		return lexInt(lex, next)
	} else if isPlainTokenStart(next) {
		return lexSymbol(lex)
	}

	return Token{}, lex.Failure(InvalidStart)
}

// lexSymbol lexes a symbol token.
// Assumes that the first symbol has already been tested as a valid symbol start.
func lexSymbol(lex *Lexer) (Token, error) {
	for isPlainTokenRest(lex.Peek(0)) {
		lex.Next()
	}
	return lex.SuccessOnBoundary(TOKEN_SYMBOL)
}

// lexInt lexes an integer token.
func lexInt(lex *Lexer, first rune) (Token, error) {
	if first == '-' || first == '+' {
		lex.Next()
	}
	for isNumASCII(lex.Peek(0)) {
		lex.Next()
	}
	return lex.SuccessOnBoundary(TOKEN_INT)
}

// lexDqs lexes a double quoted string.
// Assumes that the opening double quote has already been consumed.
func lexDqs(lex *Lexer) (Token, error) {
	for {
		switch lex.Next() {
		case '"':
			return lex.Success(TOKEN_DQS)
		case '\\':
			if glyph := lex.Next(); glyph != EOF && glyph != '\n' {
				break // The escaped character is "normal".
			}
			fallthrough
		case EOF, '\n':
			return Token{}, lex.Failure(UnbalancedDoubleQuote)
		default:
		}
	}
}

///////////////////////
// Utility functions //
// Might not be optimal, if needed could be turned into lookup tables.

// isPlainTokenRest returns whether r can go in a plain token (outside of the first character).
//
//revive:disable:line-length-limit
func isPlainTokenRest(r rune) bool {
	return strings.ContainsRune("*+-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", r)
}

// isPlainTokenStart returns whether r can be the first character of a plain token.
//
//revive:disable:line-length-limit
func isPlainTokenStart(r rune) bool {
	return strings.ContainsRune("*+-./:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", r)
}

// isNumASCII returns true if r is a numerical ASCII character.
func isNumASCII(r rune) bool {
	return ('0' <= r && r <= '9')
}
