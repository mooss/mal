package typ

import (
	"errors"
	"strings"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/ut"
)

// Map is a mapping from values to values.
type Map map[core.Value]core.Value

func (Map) Type() core.TypeName { return core.MAP }
func (m Map) String() string {
	return "{" + strings.Join(_mapflat(core.Value.String, m), " ") + "}"
}
func (m Map) Eval(core.Env) core.Value { return m }

// MapLit is an unevaluated sequence of map keys and values.
type MapLit []core.Value

func (MapLit) Type() core.TypeName { return core.MAP_LITERAL }
func (m MapLit) String() string {
	return "{" + strings.Join(_map(core.Value.String, m), " ") + "}"
}
func (m MapLit) ValidateEven() error {
	if len(m)%2 != 0 {
		return errors.New("odd-numbered map literal")
	}
	return nil
}

// Eval evaluates all keys and values in order, returning a fatal error when its length is not even
// or at the first invalid key encountered.
func (m MapLit) Eval(envr core.Env) core.Value {
	if err := m.ValidateEven(); err != nil {
		return Fatal{err}
	}

	res := make(Map, len(m)/2)
	for i := 0; i < len(m); i += 2 {
		key := m[i].Eval(envr)
		if ut.Is[Fatal](key) {
			return key
		}
		if !validKey(key) {
			return Fatalf("invalid key type in map literal: %s", key.Type())
		}

		value := m[i+1].Eval(envr)
		if ut.Is[Fatal](value) {
			return value
		}

		res[key] = value
	}
	return res
}
