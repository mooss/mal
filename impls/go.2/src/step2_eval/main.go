package main

import (
	"github.com/mooss/mal/src/env"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/repl"
	"github.com/mooss/mal/src/ut"
)

//revive:disable:cognitive-complexity
func main() {
	ut.Must(repl.Run(func(input string) (string, error) {
		parsed, err := reader.ReadProgram(input)
		if err != nil {
			return "", err
		}

		return parsed.Eval(env.Initial).String(), nil
	}))
}
