package lexer

type TokenType int

//revive:disable:var-naming
const (
	TOKEN_ERROR  TokenType = iota
	TOKEN_LPAREN           // Opening parenthesis (
	TOKEN_RPAREN           // Closing parenthesis )
	TOKEN_DQS              // Double quoted string "..."
	TOKEN_SYMBOL           // Symbol.
	TOKEN_INT              // Integer.
)

type Token struct {
	Type  TokenType
	Value string
}

func (t TokenType) String() string {
	switch t {
	case TOKEN_ERROR:
		return "ERROR"
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RPAREN:
		return "RPAREN"
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
