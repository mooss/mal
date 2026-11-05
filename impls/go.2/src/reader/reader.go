package reader

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lexer"
)

// Reader constructs Mal types from a string.
type Reader struct {
	tokens   []lexer.Token
	position int
	tokerr   func() error
}

// ReadStr returns a Reader parsing the given string.
func ReadStr(input string) Reader {
	tokenIt, errFn := lexer.TokenIterator(input)
	return Reader{
		tokens:   slices.Collect(tokenIt),
		position: 0,
		tokerr:   errFn,
	}
}

// Next returns the next token and advances the input.
func (r *Reader) Next() lexer.Token {
	res := r.tokens[r.position]
	r.position++
	return res
}

// Peek returns the next token without advancing the input.
func (r *Reader) Peek() lexer.Token {
	return r.tokens[r.position]
}

// Done returns true when all the input has been read.
func (r *Reader) Done() bool {
	return r.position >= len(r.tokens)
}

// Rest returns all the tokens that are yet to be read.
func (r *Reader) Rest() []lexer.Token {
	if r.Done() {
		return nil
	}

	return r.tokens[r.position:]
}

// wrap wraps err into a reader error.
func (*Reader) wrap(err error) error {
	return fmt.Errorf("reader: %w", err)
}

// wrapf wraps err into a reader error, formatted with msg and args.
func (*Reader) wrapf(err error, msg string, args ...any) error {
	return fmt.Errorf("reader: "+msg+": %w", append(args, err)...)
}

// err returns a reader error with the given error message.
func (*Reader) err(msg string) error {
	return errors.New("reader: " + msg)
}

// errf returns a reader error formatted with msg and args.
func (*Reader) errf(msg string, args ...any) error {
	return fmt.Errorf("reader: "+msg, args...)
}

//////////////////
// Read methods //

// ReadForm constructs the next type from the input tokens, after ensuring that there was no
// tokenization error and that there are tokens left to read.
func (r *Reader) ReadForm() (core.Type, error) {
	if err := r.tokerr(); err != nil {
		return nil, r.wrap(err)
	}
	if r.Done() {
		return nil, r.err("no tokens left to read")
	}
	return r.readForm()
}

// readForm constructs the next type from the input tokens.
//
//revive:disable:cyclomatic
func (r *Reader) readForm() (core.Type, error) {
	token := r.Next()
	switch token.Type {
	// Balanced expressions.
	case lexer.TOKEN_LPAREN:
		return r.readSequence[core.List](lexer.TOKEN_RPAREN)
	case lexer.TOKEN_RPAREN:
		return nil, r.err("unexpected closing paren")
	case lexer.TOKEN_LBRACKET:
		return r.readSequence[core.Vector](lexer.TOKEN_RBRACKET)
	case lexer.TOKEN_RBRACKET:
		return nil, r.err("unexpected closing bracket")
	case lexer.TOKEN_LBRACE:
		return r.readDummy(lexer.TOKEN_RBRACE)
	case lexer.TOKEN_RBRACE:
		return nil, r.err("unexpected closing dummy")

	// Quotes, atoms and errors.
	case lexer.TOKEN_SQUOTE:
		return r.readQuote()
	case lexer.TOKEN_DQS:
		return r.readString(token)
	case lexer.TOKEN_SYMBOL:
		return r.readSymbol(token)
	case lexer.TOKEN_INT:
		return r.readInt(token)
	case lexer.TOKEN_ERROR:
		return nil, r.err("unexpected error token")
	default:
		return nil, r.err("unexpected token type")
	}
}

// readSequence constructs a sequence type from the remaining input tokens (assumes that the opening
// delimiter has already been consumed).
func (r *Reader) readSequence[T ~[]core.Type](closing lexer.TokenType) (T, error) {
	var res T
	for {
		if r.Done() {
			return nil, r.errf("unbalanced expression (missing %s)", closing.What())
		}

		if r.Peek().Type == closing {
			r.position++
			return res, nil
		}

		next, err := r.readForm()
		if err != nil {
			return nil, err
		}
		res = append(res, next)
	}
}

// readInt constructs an int from the given input token.
func (r *Reader) readInt(token lexer.Token) (core.Type, error) {
	integer, err := strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		return nil, r.wrapf(err, "cannot read %q as int", token.Value)
	}
	return core.Integer(integer), nil
}

// readSymbol constructs a symbol from the given input token.
func (*Reader) readSymbol(token lexer.Token) (core.Type, error) {
	return core.Symbol(token.Value), nil
}

// readString constructs a string from the given input token.
func (r *Reader) readString(token lexer.Token) (core.Type, error) {
	str, err := strconv.Unquote(token.Value)
	if err != nil {
		return nil, r.wrapf(err, "cannot read %q as string", token.Value)
	}
	return core.String(str), nil
}

// readQuote quotes the next form.
func (r *Reader) readQuote() (core.Type, error) {
	if r.Done() {
		return nil, r.err("nothing left to quote")
	}
	quoted, err := r.readForm()
	if err != nil {
		return nil, err
	}
	return core.List([]core.Type{core.Symbol("quote"), quoted}), nil
}

func (r *Reader) readDummy(closing lexer.TokenType) (core.Dummy, error) {
	res := core.Dummy{}
	for {
		if r.Done() {
			// Would just say "missing closing paren" but that is not what mal expects.
			return res, r.err("unbalanced expression (dummy)")
		}

		if r.Peek().Type == closing {
			r.position++
			return res, nil
		}

		_, err := r.readForm()
		if err != nil {
			return res, err
		}
	}
}
