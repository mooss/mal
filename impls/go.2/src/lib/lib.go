// Package lib provides Mal features written in Mal.
package lib

import (
	"fmt"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/typ"
)

var Std = Package{
	Name: "standard",
	Functions: []Function{{
		"load-file",
		`(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\nnil)")))))`,
	}},
}

type Function struct {
	Name string
	Impl string
}

func (fun Function) Load(envr core.Env) error {
	parsed, err := reader.ParseOne(fun.Impl)
	if err != nil {
		return err
	}

	res := typ.Resolve(envr, parsed)
	fatal, isFatal := res.(typ.Fatal)
	if isFatal {
		return fmt.Errorf("function %s: %w", fun.Name, fatal.Unbox())
	}

	return nil
}

type Package struct {
	Name      string
	Functions []Function
}

func (p Package) Load(envr core.Env) error {
	for _, fun := range p.Functions {
		if err := fun.Load(envr); err != nil {
			return fmt.Errorf("package %s: %w", p.Name, err)
		}
	}
	return nil
}
