package typ

import "github.com/mooss/mal/src/core"

// Function is a binary operation on integers.
type Function func(Integer, Integer) Integer

func (Function) Type() core.TypeName          { return core.FUNCTION }
func (Function) String() string               { return "<function>" }
func (fun Function) Eval(core.Env) core.Value { return fun }
