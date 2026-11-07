package main

import (
	"github.com/mooss/mal/src/repl"
	"github.com/mooss/mal/src/ut"
)

func main() {
	ut.Must(repl.Run(func(input string) (string, error) {
		return input, nil
	}))
}
