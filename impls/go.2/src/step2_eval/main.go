package main

import (
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/repl"
	"github.com/mooss/mal/src/ut"
)

//revive:disable:cognitive-complexity
func main() {
	rep := func(input string) (string, error) {
		parsed, err := reader.ParseOne(input)
		if err != nil {
			return "", err
		}

		return parsed.Eval(envs.Math).String(), nil
	}
	ut.Must(repl.Run(rep))
}
