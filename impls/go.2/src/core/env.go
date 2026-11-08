package core

import (
	"maps"
	"slices"
)

// Bindings is a primitive environment associating symbols to their values.
type Bindings map[string]Value

func (b Bindings) AsEnv() Env {
	if b == nil {
		b = map[string]Value{}
	}
	return Env{inner: b}
}

// Env is an environment associating symbols to their values.
// The zero value is invalid and will panic when setting.
type Env struct {
	inner Bindings
	outer *Env
}

// NewEnv creates a new environment, optionally pre-populated with the given bindings.
func NewEnv(bindings ...Bindings) Env {
	store := map[string]Value{}
	for _, b := range bindings {
		maps.Copy(store, b)
	}
	return Env{inner: store}
}

// Cat concatenates environments, creating a new flat environment without any outer environment.
// The right-most environments take precedence, as well as the inner-most bindings.
func CatEnv(envs ...Env) Env {
	res := NewEnv()
	for _, env := range slices.Backward(envs) {
		for penv := &env; penv != nil; penv = penv.outer {
			for key, value := range penv.inner {
				if _, exists := res.inner[key]; !exists {
					res.inner[key] = value
				}
			}
		}
	}
	return res
}

// Get returns the value bound to key, nil if the key is unbound.
func (e Env) Get(key string) Value {
	for penv := &e; penv != nil; penv = penv.outer {
		value, present := penv.inner[key]
		if present {
			return value
		}
	}
	return nil
}

// Nest returns a new environment using inner as inner bindings and the current environment as the
// outer environment.
func (e Env) Nest(inner Bindings) Env {
	if inner == nil {
		inner = Bindings{}
	}
	return Env{inner: inner, outer: &e}
}

// Set binds key to val in the environment.
func (e Env) Set(key string, val Value) {
	e.inner[key] = val
}
