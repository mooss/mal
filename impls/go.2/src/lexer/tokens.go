package lexer

type TokenType int

//revive:disable:var-naming
const (
	TOKEN_ERROR    TokenType = iota
	TOKEN_LPAREN             // Opening parenthesis (
	TOKEN_RPAREN             // Closing parenthesis )
	TOKEN_LBRACKET           // Opening square bracket [
	TOKEN_RBRACKET           // Closing square bracket ]
	TOKEN_LBRACE             // Opening curly brace {
	TOKEN_RBRACE             // Closing curly brace }
	TOKEN_DQS                // Double quoted string "..."
	TOKEN_SYMBOL             // Symbol.
	TOKEN_INT                // Integer.
)

type Token struct {
	Type  TokenType
	Value string
}

//revive:disable:cyclomatic
func (t TokenType) String() string {
	switch t {
	case TOKEN_ERROR:
		return "ERROR"
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RPAREN:
		return "RPAREN"
	case TOKEN_LBRACKET:
		return "LBRACKET"
	case TOKEN_RBRACKET:
		return "RBRACKET"
	case TOKEN_LBRACE:
		return "LBRACE"
	case TOKEN_RBRACE:
		return "RBRACE"
	case TOKEN_DQS:
		return "DQS"
	case TOKEN_SYMBOL:
		return "SYMBOL"
	case TOKEN_INT:
		return "INT"
	default:
		return "UNKNOWN"
	}
}
