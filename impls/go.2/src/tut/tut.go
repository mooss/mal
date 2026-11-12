// Package tut provides test utilities.
package tut

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/typ"
)

func MustEval(envr core.Env, input string) core.Value {
	program, err := reader.ParseOne(input)
	if err != nil {
		panic(err)
	}

	return typ.Resolve(envr, program)
}
