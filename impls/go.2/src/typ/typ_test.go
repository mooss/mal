package typ_test

// Internal utilities for the typ_test package.

import (
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

type evalOrder struct {
	name   string
	result core.Value
	order  *[]string
}

func (v evalOrder) Type() core.TypeName { return v.result.Type() }
func (v evalOrder) String() string      { return v.name }
func (evalOrder) Equal(core.Value) bool { return false }
func (v evalOrder) Eval(core.Env) core.Value {
	*v.order = append(*v.order, v.name)
	return v.result
}

// logger records evaluation order and returns val.
type logger struct {
	id  string
	log *[]string
	val core.Value
}

func (logger) Type() core.TypeName   { return core.LIST }
func (p logger) String() string      { return p.id }
func (logger) Equal(core.Value) bool { return false }
func (p logger) Eval(core.Env) core.Value {
	*p.log = append(*p.log, p.id)
	return p.val
}

func requireError(t *testing.T, value core.Value) {
	t.Helper()
	if !ut.Is[typ.Fatal](value) {
		t.Fatalf("Eval() = %#v, want %T", value, typ.Fatal{})
	}
}
