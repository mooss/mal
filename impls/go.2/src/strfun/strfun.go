// Package strfun implements string manipulation functions.
package strfun

import (
	"errors"
	"fmt"
)

// MalUnquote takes a double quoted string, including the quotes, and returns its unquoted contents
// with escape sequences resolved according to Mal rules.
//
//revive:disable:cognitive-complexity
func MalUnquote(quoted string) (string, error) {
	if len(quoted) < 2 || quoted[0] != '"' || quoted[len(quoted)-1] != '"' {
		return "", errors.New("mal-unquote: missing surrounding quotes")
	}
	inner := quoted[1 : len(quoted)-1]
	out := make([]byte, 0, len(inner))
	for i := 0; i < len(inner); i++ {
		if inner[i] == '\\' {
			if i+1 >= len(inner) {
				return "", errors.New("mal-unquote: invalid escape \\")
			}
			next := inner[i+1]
			switch next {
			case 'n':
				out = append(out, '\n')
				i++
				continue
			case '"', '\\':
				out = append(out, next)
				i++
				continue
			default:
				return "", fmt.Errorf("mal-unquote: invalid escape \\%c", next)
			}
		}
		out = append(out, inner[i])
	}
	return string(out), nil
}

// MalQuote returns s as a double quoted string, including the quotes, escaping backslashes, double
// quotes, and newlines according to Mal rules.
func MalQuote(s string) string {
	out := make([]byte, 0, len(s)+2)
	out = append(out, '"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			out = append(out, '\\', '\\')
		case '"':
			out = append(out, '\\', '"')
		case '\n':
			out = append(out, '\\', 'n')
		default:
			out = append(out, s[i])
		}
	}
	out = append(out, '"')
	return string(out)
}
