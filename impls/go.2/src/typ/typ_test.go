package typ_test

// Internal utilities for the typ_test package.

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

type evalCounter struct {
	count *int
}

func (evalCounter) Type() core.TypeName { return core.INTEGER }
func (evalCounter) String() string      { return "counter" }
func (v evalCounter) Eval(core.Env) core.Value {
	(*v.count)++
	return typ.Integer(1)
}

type evalOrder struct {
	name   string
	result core.Value
	order  *[]string
}

func (v evalOrder) Type() core.TypeName { return v.result.Type() }
func (v evalOrder) String() string      { return v.name }
func (v evalOrder) Eval(core.Env) core.Value {
	*v.order = append(*v.order, v.name)
	return v.result
}

func requireError(t *testing.T, value core.Value) {
	t.Helper()
	if !ut.Is[typ.Fatal](value) {
		t.Fatalf("Eval() = %#v, want %T", value, typ.Fatal{})
	}
}
