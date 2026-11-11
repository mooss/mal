package main

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/repl"
	"github.com/mooss/mal/src/ut"
)

//revive:disable:cognitive-complexity
func main() {
	environment := core.CatEnv(envs.MalDefault, envs.ExtraSpecial)
	rep := func(input string) (string, error) {
		parsed, err := reader.ParseOne(input)
		if err != nil {
			return "", err
		}

		return parsed.Eval(environment).String(), nil
	}
	ut.Must(repl.Run(rep))
}
