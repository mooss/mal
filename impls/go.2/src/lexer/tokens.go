package lexer

type TokenType int

const (
	TOKEN_ERROR  TokenType = iota
	TOKEN_LPAREN           // Opening parenthesis (
	TOKEN_RPAREN           // Closing parenthesis )
	TOKEN_DQS              // Double quoted string "..."
	TOKEN_SYMBOL           // Symbol.
)

type Token struct {
	Type  TokenType
	Value string
}
