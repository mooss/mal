package main

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/repl"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func main() {
	environment := core.CatEnv(envs.All)
	rep := func(input string) (string, error) {
		parsed, err := reader.ParseOne(input)
		if err != nil {
			return "", err
		}

		return typ.Resolve(environment, parsed).String(), nil
	}
	ut.Must(repl.Run(rep))
}
