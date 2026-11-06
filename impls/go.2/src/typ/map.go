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

// Assoc represents a key-value association.
type Assoc struct {
	Key   core.Value
	Value core.Value
}

// MapForm is an unevaluated sequence of map key/value pairs.
type MapForm []Assoc

// NewMapForm builds a MapForm from a flat sequence of keys and values.
// It returns an error when vals holds an odd number of elements.
func NewMapForm(vals ...core.Value) (MapForm, error) {
	if len(vals)%2 != 0 {
		return nil, errors.New("odd-numbered map literal")
	}
	res := make(MapForm, 0, len(vals)/2)
	for i := 0; i < len(vals); i += 2 {
		res = append(res, Assoc{Key: vals[i], Value: vals[i+1]})
	}
	return res, nil
}

func (MapForm) Type() core.TypeName { return core.MAP_LITERAL }
func (m MapForm) String() string {
	parts := make([]string, 0, len(m)*2)
	for _, assoc := range m {
		parts = append(parts, assoc.Key.String(), assoc.Value.String())
	}
	return "{" + strings.Join(parts, " ") + "}"
}

// Flatten flattens the map form back into a sequence of alternating keys and values.
func (m MapForm) Flatten() []core.Value {
	res := make([]core.Value, 0, len(m)*2)
	for _, assoc := range m {
		res = append(res, assoc.Key, assoc.Value)
	}
	return res
}

// Eval evaluates all keys and values in order, returning a fatal error at the first invalid key
// encountered.
func (m MapForm) Eval(envr core.Env) core.Value {
	res := make(Map, len(m))
	for _, assoc := range m {
		key := assoc.Key.Eval(envr)
		if ut.Is[Fatal](key) {
			return key
		}
		if !validKey(key) {
			return Fatalf("invalid key type in map literal: %s", key.Type())
		}

		value := assoc.Value.Eval(envr)
		if ut.Is[Fatal](value) {
			return value
		}

		res[key] = value
	}
	return res
}
