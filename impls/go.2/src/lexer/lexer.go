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

// Peek returns the next rune without advancing the input.
func (lex *Lexer) Peek() rune {
	if lex.pos >= len(lex.Input) {
		return EOF
	}
	return lex.Input[lex.pos]
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
	peek := lex.Peek()
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
		switch lex.Peek() {
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
	peek := lex.Peek()
	switch peek {
	case EOF:
		return Token{}, EndOfInput
	case '(':
		lex.Next()
		return lex.Success(TOKEN_LPAREN)
	case ')':
		lex.Next()
		return lex.Success(TOKEN_RPAREN)
	case '[':
		lex.Next()
		return lex.Success(TOKEN_LBRACKET)
	case ']':
		lex.Next()
		return lex.Success(TOKEN_RBRACKET)
	case '{':
		lex.Next()
		return lex.Success(TOKEN_LBRACE)
	case '}':
		lex.Next()
		return lex.Success(TOKEN_RBRACE)
	case '\'':
		lex.Next()
		return lex.Success(TOKEN_SQUOTE)
	case '"':
		return lexDqs(lex)
	default:
	}

	if isNumASCII(peek) {
		return lexInt(lex)
	} else if peek == '-' && lex.pos+1 < len(lex.Input) && isNumASCII(lex.Input[lex.pos+1]) {
		return lexInt(lex)
	} else if isPlainTokenStart(peek) {
		return lexSymbol(lex)
	}

	lex.Next() // Advance the input to show the rune in the error message.
	return Token{}, lex.Failure(InvalidStart)
}

// lexSymbol lexes a symbol token.
// Assumes that the first symbol has already been tested as a valid symbol start.
func lexSymbol(lex *Lexer) (Token, error) {
	for isPlainTokenRest(lex.Peek()) {
		lex.Next()
	}
	return lex.SuccessOnBoundary(TOKEN_SYMBOL)
}

// lexInt lexes an integer token.
// A leading minus is allowed for negative integers; the caller ensures a digit follows it.
func lexInt(lex *Lexer) (Token, error) {
	if lex.Peek() == '-' {
		lex.Next()
	}
	for isNumASCII(lex.Peek()) {
		lex.Next()
	}
	return lex.SuccessOnBoundary(TOKEN_INT)
}

// lexDqs lexes a double quoted string.
// The lexer is supposed to be placed at the opening double quote.
func lexDqs(lex *Lexer) (Token, error) {
	lex.Next() // Consume the opening double quote.
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
