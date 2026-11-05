package lexer

type TokenType int

//revive:disable:var-naming
const (
	TOKEN_ERROR    TokenType = iota
	TOKEN_LPAREN             // Opening paren (
	TOKEN_RPAREN             // Closing paren )
	TOKEN_LBRACKET           // Opening bracket [
	TOKEN_RBRACKET           // Closing bracket ]
	TOKEN_LBRACE             // Opening brace {
	TOKEN_RBRACE             // Closing brace }
	TOKEN_SQUOTE             // Single quote '
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
	case TOKEN_SQUOTE:
		return "SQUOTE"
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

// What returns a human readable string describing what type of token it is.
func (t TokenType) What() string {
	switch t {
	case TOKEN_ERROR:
		return "error"
	case TOKEN_LPAREN:
		return "opening paren"
	case TOKEN_RPAREN:
		return "closing paren"
	case TOKEN_LBRACKET:
		return "opening bracket"
	case TOKEN_RBRACKET:
		return "closing bracket"
	case TOKEN_LBRACE:
		return "opening brace"
	case TOKEN_RBRACE:
		return "closing brace"
	case TOKEN_SQUOTE:
		return "single quote"
	case TOKEN_DQS:
		return "double quoted string"
	case TOKEN_SYMBOL:
		return "symbol"
	case TOKEN_INT:
		return "integer"
	default:
		return "unknown"
	}
}
