package tut

import (
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/typ"
)

// TestVal is a minimal distinct Value for env tests.
type TestVal string

func (TestVal) Type() core.TypeName        { return core.TypeName("test") }
func (v TestVal) String() string           { return string(v) }
func (v TestVal) Eval(core.Env) core.Value { return v }
func (v TestVal) Equal(o core.Value) bool {
	other, ok := o.(TestVal)
	return ok && v == other
}

// Probe is a Value that records Eval order and returns Ret.
type Probe struct {
	ID  string
	Log *[]string
	Ret core.Value
}

func (p Probe) Type() core.TypeName {
	if p.Ret != nil {
		return p.Ret.Type()
	}
	return core.TypeName("probe")
}
func (p Probe) String() string      { return p.ID }
func (Probe) Equal(core.Value) bool { return false }
func (p Probe) Eval(core.Env) core.Value {
	if p.Log != nil {
		*p.Log = append(*p.Log, p.ID)
	}
	if p.Ret != nil {
		return p.Ret
	}
	return typ.Integer(1)
}
