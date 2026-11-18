package reader

import (
	"errors"
	"fmt"
	"io"
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

// ParseOne parses the input string as one expression, failing when there are leftover tokens.
// Will return a error wrapping io.EOF when input contains no expression.
func ParseOne(input string) (core.Value, error) {
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

// ParseAll parses the input string as a sequence of expressions.
// Empty input yields an empty slice and no error.
func ParseAll(input string) ([]core.Value, error) {
	rdr := ReadStr(input)
	var forms []core.Value
	for {
		form, err := rdr.ReadForm()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		forms = append(forms, form)
	}

	if rest := rdr.Rest(); rest != nil {
		return nil, fmt.Errorf("leftover tokens remain: %v", rest)
	}

	return forms, nil
}
