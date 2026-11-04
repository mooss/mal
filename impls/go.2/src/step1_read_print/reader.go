package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

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

//////////////////
// Read methods //

// ReadForm constructs the next type from the input tokens, after ensuring that there was no
// tokenization error and that there are tokens left to read.
func (r *Reader) ReadForm() (MalType, error) {
	if err := r.tokerr(); err != nil {
		return nil, r.wrap(err)
	}
	if r.Done() {
		return nil, r.err("no tokens left to read")
	}
	return r.readForm()
}

// readForm constructs the next type from the input tokens.
func (r *Reader) readForm() (MalType, error) {
	token := r.Next()
	switch token.Type {
	case lexer.TOKEN_LPAREN:
		return r.readList()
	case lexer.TOKEN_RPAREN:
		return nil, r.err("unexpected closing paren")
	case lexer.TOKEN_ERROR:
		return nil, r.err("unexpected error token")
	case lexer.TOKEN_DQS, lexer.TOKEN_SYMBOL, lexer.TOKEN_INT:
		return r.readAtom(token)
	default:
		return nil, r.err("unexpected token type")
	}
}

// readList constructs a list from the remaining input tokens (assumes that the opening paren has
// already been consumed).
func (r *Reader) readList() (List, error) {
	res := List{}
	for {
		if r.Done() {
			return nil, r.err("missing closing paren")
		}

		if r.Peek().Type == lexer.TOKEN_RPAREN {
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

// readAtom constructs an atom from the next input token.
func (r *Reader) readAtom(token lexer.Token) (MalType, error) {
	if token.Type == lexer.TOKEN_INT {
		integer, err := strconv.ParseInt(token.Value, 10, 64)
		if err != nil {
			return nil, r.wrapf(err, "cannot read %q as int", token.Value)
		}
		return Integer(integer), nil
	}

	return Symbol(token.Value), nil
}
