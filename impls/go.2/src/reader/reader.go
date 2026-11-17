package reader

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lexer"
	"github.com/mooss/mal/src/strfun"
	"github.com/mooss/mal/src/typ"
)

// Reader constructs Mal types from a string.
type Reader struct {
	tokens   []lexer.Token
	position int
	tokerr   func() error
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
// Returns an error wrapping io.EOF when there are no more tokens left to read.
func (r *Reader) ReadForm() (core.Value, error) {
	if err := r.tokerr(); err != nil {
		return nil, err
	}
	if r.Done() {
		return nil, r.wrap(io.EOF)
	}
	return r.readForm()
}

// readForm constructs the next type from the input tokens.
func (r *Reader) readForm() (core.Value, error) {
	token := r.Next()
	switch token.Type {
	// Balanced expressions.
	case lexer.TOKEN_DELIMITER:
		return r.dispatchDelim(token.Value)
	// Quotes, atoms and errors.
	case lexer.TOKEN_ERROR:
		return nil, r.err("unexpected error token")
	case lexer.TOKEN_SIGIL:
		return r.readSigil(token)
	case lexer.TOKEN_DQS:
		return r.readString(token)
	case lexer.TOKEN_INT:
		return r.readInt(token)
	case lexer.TOKEN_SYMBOL:
		return r.readSymbol(token)
	default:
		return nil, r.err("unexpected token type")
	}
}

// dispatchDelim dispatches a read delimiter to the appropriate read* method.
func (r *Reader) dispatchDelim(value string) (core.Value, error) {
	switch value {
	case "(":
		return r.readSequence[typ.List](")")
	case ")":
		return nil, r.err("unexpected closing paren")
	case "[":
		return r.readSequence[typ.VectorForm]("]")
	case "]":
		return nil, r.err("unexpected closing bracket")
	case "{":
		return r.readMapLit()
	case "}":
		return nil, r.err("unexpected closing brace")
	default:
		return nil, r.err("unexpected delimiter")
	}
}

// seekDelim seeks a balanced expression delimiter.
// Returns:
//   - An error when the reader is done.
//   - nil, nil when the delimiter is found.
//   - The next form otherwise.
func (r *Reader) seekDelim(delim string) (core.Value, error) {
	if r.Done() {
		return nil, r.errf("unbalanced expression (missing %q)", delim)
	}

	if r.Peek().Type == lexer.TOKEN_DELIMITER && r.Peek().Value == delim {
		r.position++
		return nil, nil
	}

	return r.readForm()
}

// readSequence constructs a sequence type from the remaining input tokens (assumes that the opening
// delimiter has already been consumed).
func (r *Reader) readSequence[T ~[]core.Value](closing string) (T, error) {
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
		key, err := r.seekDelim("}")
		if err != nil {
			return nil, err
		}
		if key == nil {
			return res, nil
		}

		value, err := r.seekDelim("}")
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
	if token.Value[0] == ':' {
		return typ.Keyword(token.Value[1:]), nil
	}
	return typ.Symbol(token.Value), nil
}

// readString constructs a string from the given input token.
func (r *Reader) readString(token lexer.Token) (core.Value, error) {
	res, err := strfun.MalUnquote(token.Value)
	if err != nil {
		return nil, r.wrap(err)
	}

	return typ.String(res), nil
}

// readSigil dispatches on the sigil value, allowing longer prefixes (e.g. ->) later.
func (r *Reader) readSigil(token lexer.Token) (core.Value, error) {
	switch token.Value {
	case "'":
		return r.expandSigil("quote")
	case "`":
		return r.expandSigil("quasiquote")
	case "~":
		return r.expandSigil("unquote")
	case "@":
		return r.expandSigil("deref")
	case "~@":
		return r.expandSigil("splice-unquote")
	default:
		return nil, r.errf("unexpected sigil %q", token.Value)
	}
}

// expandSigil wraps the next form within a list application starting with the given symbol.
func (r *Reader) expandSigil(symbol string) (core.Value, error) {
	if r.Done() {
		return nil, r.errf("nothing left to %s", symbol)
	}
	quoted, err := r.readForm()
	if err != nil {
		return nil, err
	}

	return typ.NewList(typ.Symbol(symbol), quoted), nil
}
