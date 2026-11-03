package lexer

import (
	"errors"
	"fmt"
	"iter"
	"unicode/utf8"
)

// TokenIterator turns input into an iterable sequence over its tokens and an error getter.
// The error should always be checked after iterating.
//
//revive:disable:cognitive-complexity
func TokenIterator(input string) (iter.Seq[Token], func() error) {
	lex := &Lexer{Input: input}
	var err error
	return func(yield func(Token) bool) {
		for {
			token, scanErr := lexMal(lex)
			if scanErr != nil {
				if errors.Is(scanErr, EndOfInput) {
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
	UnterminatedDQS LexicalFailure = "unterminated double quoted string"
	InvalidStart    LexicalFailure = "invalid token start"
)

const EOF = -1

type Lexer struct {
	Input string // The string being lexed.
	start int    // Start position of the current token.
	pos   int    // Current position in the input.
	width int    // Width of the last rune read.
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
	glyph, _ := utf8.DecodeRuneInString(lex.Input[lex.pos:])
	return glyph
}

// Next advances to the Next rune, returning the current one.
func (lex *Lexer) Next() rune {
	if lex.pos >= len(lex.Input) {
		return EOF
	}

	var glyph rune
	glyph, lex.width = utf8.DecodeRuneInString(lex.Input[lex.pos:])
	lex.pos += lex.width
	return glyph
}

// produce returns the current token and advances the input.
func (lex *Lexer) produce(typ TokenType) Token {
	token := lex.Input[lex.start:lex.pos]
	lex.start = lex.pos
	return Token{Type: typ, Value: token}
}

// Failure returns a lexical error wrapping the current token.
func (lex *Lexer) Failure(reason LexicalFailure) error {
	token := lex.produce(TOKEN_ERROR)
	return fmt.Errorf("%s: %#q", reason, token.Value)
}

// Success returns the current token without error and advances the input.
func (lex *Lexer) Success(typ TokenType) (Token, error) {
	return lex.produce(typ), nil
}

////////////
// Lexers //

// lexMal extracts the next Mal token.
// It returns the token on success, EndOfInput on end-of-input, or another error on lexical failure.
func lexMal(lex *Lexer) (Token, error) {
	switch lex.Peek() {
	case EOF:
		return Token{}, EndOfInput
	case '"':
		return lexDqs(lex)
	default:
		lex.Next() // Advance the input to show the rune in the error message.
		return Token{}, lex.Failure(InvalidStart)
	}
}

// lexDqs seeks the end of a double quoted string.
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
			return Token{}, lex.Failure(UnterminatedDQS)
		default:
		}
	}
}
