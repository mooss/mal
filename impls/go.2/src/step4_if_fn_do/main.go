package main

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/repl"
	"github.com/mooss/mal/src/ut"
)

func main() {
	environment := core.CatEnv(envs.All)
	ut.Must(repl.RunEnv(environment))
}
