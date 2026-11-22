package main

import (
	"os"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/lib"
	"github.com/mooss/mal/src/repl"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

func main() {
	environment := core.CatEnv(envs.All)
	environment.Set("*step*", typ.String("A"))
	ut.Must(lib.StdMal.Load(environment))
	ut.Must(repl.OrExecute(environment, os.Args[1:]...))
}
