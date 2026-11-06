package reader

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lexer"
	"github.com/mooss/mal/src/typ"
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
func (r *Reader) ReadForm() (core.Value, error) {
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
func (r *Reader) readForm() (core.Value, error) {
	token := r.Next()
	switch token.Type {
	// Balanced expressions.
	case lexer.TOKEN_LPAREN:
		return r.readSequence[typ.List](lexer.TOKEN_RPAREN)
	case lexer.TOKEN_RPAREN:
		return nil, r.err("unexpected closing paren")
	case lexer.TOKEN_LBRACKET:
		return r.readSequence[typ.Vector](lexer.TOKEN_RBRACKET)
	case lexer.TOKEN_RBRACKET:
		return nil, r.err("unexpected closing bracket")
	case lexer.TOKEN_LBRACE:
		return r.readMapLit()
	case lexer.TOKEN_RBRACE:
		return nil, r.err("unexpected closing brace")

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

// seekDelim seeks a balanced expression delimiter.
// Returns:
//   - An error when the reader is done.
//   - nil, nil when the delimiter is found.
//   - The next form otherwise.
func (r *Reader) seekDelim(delim lexer.TokenType) (core.Value, error) {
	if r.Done() {
		return nil, r.errf("unbalanced expression (missing %s)", delim.What())
	}

	if r.Peek().Type == delim {
		r.position++
		return nil, nil
	}

	return r.readForm()
}

// readSequence constructs a sequence type from the remaining input tokens (assumes that the opening
// delimiter has already been consumed).
func (r *Reader) readSequence[T ~[]core.Value](closing lexer.TokenType) (T, error) {
	var res T
	for {
		next, err := r.seekDelim(closing)
		if err != nil {
			return nil, err
		}

		if next == nil {
			return res, nil
		}

		res = append(res, next)
	}
}

// readMapLit constructs a map literal from the remaining input tokens (assumes that the opening
// delimiter has already been consumed).
// Map literals are represented as a slice of key/value pairs and only instantiated as maps during
// evaluation.
func (r *Reader) readMapLit() (typ.MapForm, error) {
	var res typ.MapForm
	for {
		key, err := r.seekDelim(lexer.TOKEN_RBRACE)
		if err != nil {
			return nil, err
		}
		if key == nil {
			return res, nil
		}

		value, err := r.seekDelim(lexer.TOKEN_RBRACE)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, r.err("odd-numbered map literal")
		}

		res = append(res, typ.Assoc{Key: key, Value: value})
	}
}

// readInt constructs an int from the given input token.
func (r *Reader) readInt(token lexer.Token) (core.Value, error) {
	integer, err := strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		return nil, r.wrapf(err, "cannot read %q as int", token.Value)
	}
	return typ.Integer(integer), nil
}

// readSymbol constructs a symbol from the given input token.
func (*Reader) readSymbol(token lexer.Token) (core.Value, error) {
	if token.Value == "nil" {
		return typ.Nil{}, nil
	}
	return typ.Symbol(token.Value), nil
}

// readString constructs a string from the given input token.
func (r *Reader) readString(token lexer.Token) (core.Value, error) {
	str, err := strconv.Unquote(token.Value)
	if err != nil {
		return nil, r.wrapf(err, "cannot read %q as string", token.Value)
	}
	return typ.String(str), nil
}

// readQuote quotes the next form.
func (r *Reader) readQuote() (core.Value, error) {
	if r.Done() {
		return nil, r.err("nothing left to quote")
	}
	quoted, err := r.readForm()
	if err != nil {
		return nil, err
	}
	return typ.List([]core.Value{typ.Symbol("quote"), quoted}), nil
}
