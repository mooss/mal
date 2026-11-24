package lexer

import (
	"iter"
	"strings"

	"github.com/mooss/harper/go/lex"
)

type Token = lex.Token[Kind]

var mal = lex.Dialect[Kind]{
	Whitespace:  lex.InSet(" \t\n\r,"),
	LineComment: ';',
	Rules: []lex.Rule[Kind]{ //TODO: reorder?
		lex.AnyOf("()[]{}", KindDelimiter),
		lex.Quoted('"', '\\', lex.LexicalFailure("unbalanced double quote"), KindDQS),
		lex.LongestPrefix(KindSigil, "'", "`", "~@", "~", "@"),
		lex.SignedInt(lex.IsDigitASCII, malBoundary, KindInt),
		lex.Symbol(isMalStart, isMalRest, malBoundary, KindSymbol),
	},
}

func TokenIterator(input string) (iter.Seq[Token], func() error) {
	return lex.TokenIterator(input, mal.Scanner())
}

//revive:disable:line-length-limit
func isMalStart(r rune) bool {
	return strings.ContainsRune("!&*+-./:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", r)
}

//revive:disable:line-length-limit
func isMalRest(r rune) bool {
	return strings.ContainsRune("!&*+-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", r)
}

func malBoundary(r rune) bool {
	return strings.ContainsRune("()[]{} \t\n\r,;", r)
}

////////////////
// Token kind //

type Kind int

const (
	// Lexer error, must stay to top so that the default value is invalid.
	KindError Kind = iota
	KindDQS
	KindInt
	KindDelimiter
	KindSigil
	KindSymbol
)

// String returns the short name of a token kind.
func (k Kind) String() string {
	switch k {
	case KindError:
		return "ERROR"
	case KindDQS:
		return "DQS"
	case KindInt:
		return "INT"
	case KindDelimiter:
		return "DELIMITER"
	case KindSigil:
		return "SIGIL"
	case KindSymbol:
		return "SYMBOL"
	default:
		return "UNKNOWN"
	}
}

// What returns a human-readable description of a token kind.
func (k Kind) What() string {
	switch k {
	case KindError:
		return "error"
	case KindDQS:
		return "double quoted string"
	case KindInt:
		return "integer"
	case KindDelimiter:
		return "delimiter"
	case KindSigil:
		return "sigil"
	case KindSymbol:
		return "symbol"
	default:
		return "unknown"
	}
}
