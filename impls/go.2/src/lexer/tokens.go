package lexer

type TokenType int

//revive:disable:var-naming
const (
	// Lexer error, must stay to top so that the default value is invalid.
	TOKEN_ERROR    TokenType = iota
	TOKEN_DQS                // Double quoted string "..."
	TOKEN_INT                // Integer.
	TOKEN_LBRACE             // Opening brace {
	TOKEN_LBRACKET           // Opening bracket [
	TOKEN_LPAREN             // Opening paren (
	TOKEN_RBRACE             // Closing brace }
	TOKEN_RBRACKET           // Closing bracket ]
	TOKEN_RPAREN             // Closing paren )
	TOKEN_SIGIL              // Sigil prefix like ' ` ~@
	TOKEN_SYMBOL             // Symbol.
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
	case TOKEN_DQS:
		return "DQS"
	case TOKEN_INT:
		return "INT"
	case TOKEN_LBRACE:
		return "LBRACE"
	case TOKEN_LBRACKET:
		return "LBRACKET"
	case TOKEN_LPAREN:
		return "LPAREN"
	case TOKEN_RBRACE:
		return "RBRACE"
	case TOKEN_RBRACKET:
		return "RBRACKET"
	case TOKEN_RPAREN:
		return "RPAREN"
	case TOKEN_SIGIL:
		return "SIGIL"
	case TOKEN_SYMBOL:
		return "SYMBOL"
	default:
		return "UNKNOWN"
	}
}

// What returns a human readable string describing what type of token it is.
func (t TokenType) What() string {
	switch t {
	case TOKEN_ERROR:
		return "error"
	case TOKEN_DQS:
		return "double quoted string"
	case TOKEN_INT:
		return "integer"
	case TOKEN_LBRACE:
		return "opening brace"
	case TOKEN_LBRACKET:
		return "opening bracket"
	case TOKEN_LPAREN:
		return "opening paren"
	case TOKEN_RBRACE:
		return "closing brace"
	case TOKEN_RBRACKET:
		return "closing bracket"
	case TOKEN_RPAREN:
		return "closing paren"
	case TOKEN_SIGIL:
		return "sigil"
	case TOKEN_SYMBOL:
		return "symbol"
	default:
		return "unknown"
	}
}
