// Package envs provides pre-defined environments.
package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

type _int = typ.Integer

var (
	MalDefault  = core.CatEnv(MalSpecial, Math, MalFun)
	All         = core.CatEnv(MalDefault, ExtraSpecial)
	arityGt1Err = typ.Fatalf("expected at least 1 argument")
	arityGt2Err = typ.Fatalf("expected at least 2 argument")
)
