package reader

import (
	"fmt"
	"slices"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lexer"
)

// ReadStr returns a Reader parsing the given string.
func ReadStr(input string) Reader {
	tokenIt, errFn := lexer.TokenIterator(input)
	return Reader{
		tokens:   slices.Collect(tokenIt),
		position: 0,
		tokerr:   errFn,
	}
}

// ReadProgram reads the input string as one expression.
func ReadProgram(input string) (core.Value, error) {
	rdr := ReadStr(input)
	parsed, err := rdr.ReadForm()
	if err != nil {
		return nil, err
	}

	if rest := rdr.Rest(); rest != nil {
		return nil, fmt.Errorf("leftover tokens remain: %v", rest)
	}

	return parsed, nil
}
