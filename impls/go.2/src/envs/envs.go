// Package envs provides pre-defined environments.
package envs

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

type _int = typ.Integer

var (
	MalValues = core.Bindings{
		"*ARGV*":          typ.List{},
		"*host-language*": typ.NewString("go.2"),
		"false":           typ.False,
		"true":            typ.True,
	}.AsEnv()
	MalStandard = core.CatEnv(MalValues, MalSpecial, Math, MalFun)
	All         = core.CatEnv(MalStandard, ExtraSpecial)
	arityGt1Err = typ.Fatalf("expected at least 1 argument")
	arityGt2Err = typ.Fatalf("expected at least 2 argument")
)
