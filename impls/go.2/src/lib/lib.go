// Package lib provides Mal features written in Mal.
package lib

import (
	"fmt"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/envs"
	"github.com/mooss/mal/src/reader"
	"github.com/mooss/mal/src/typ"
)

// StdMal implements the Mal standard library.
var StdMal = Package{
	Name: "standard",
	GuestFun: []GuestFun{{
		"load-file",
		`(def! load-file (fn* (f)
  (do (eval (read-string (str "(do " (slurp f) "\nnil)")))) ))`,
	}, {
		"cond",
		`(defmacro! cond (fn* (& xs)
  (if (> (count xs) 0)
      (list 'if (first xs)
             (if (> (count xs) 1)
                 (nth xs 1)
                 (throw "odd number of forms to cond"))
             (cons 'cond (rest (rest xs)))))))`,
	}},
	BoundPrimitives: []EnvBound{{
		Name:  "eval",
		Maker: pipe(envs.MkEval, typ.I),
	}},
}

/////////////
// Package //

// Package is a named collection of guest functions and bound primitives that can be loaded into an
// environment.
type Package struct {
	Name            string
	GuestFun        []GuestFun
	BoundPrimitives []EnvBound
}

// Load registers all the package's guest functions and bound primitives in the given environment.
func (p Package) Load(root core.Env) error {
	for _, fun := range p.GuestFun {
		if err := fun.Load(root); err != nil {
			return fmt.Errorf("package %s: %w", p.Name, err)
		}
	}

	for _, envcap := range p.BoundPrimitives {
		root.Set(envcap.Name, envcap.Maker(root))
	}

	return nil
}

/////////////////////////////////
// Functions (guest and bound) //

// GuestFun is a function implemented in the guest language.
type GuestFun struct {
	Name string
	Impl string
}

// Load evaluates the function's implementation and registers it in the given environment.
func (fun GuestFun) Load(root core.Env) error {
	parsed, err := reader.ParseOne(fun.Impl)
	if err != nil {
		return err
	}

	res := typ.Resolve(root, parsed)
	fatal, isFatal := res.(typ.Fatal)
	if isFatal {
		return fmt.Errorf("function %s: %w", fun.Name, fatal.Unbox())
	}

	return nil
}

// EnvBound is a function bound to the root environment.
type EnvBound struct {
	Name  string
	Maker func(core.Env) typ.Function
}

/////////////////////
// Local utilities //

// pipe combines two functions into one, applying inner first and then outer.
func pipe[A, B, C any](inner func(A) B, outer func(B) C) func(A) C {
	return func(input A) C { return outer(inner(input)) }
}
