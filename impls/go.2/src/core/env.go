package core

import "maps"

// Bindings is a primitive environment associating symbols to their values.
type Bindings map[string]Value

func (b Bindings) AsEnv() Env { return Env{inner: b} }

// Env is an environment associating symbols to their values.
type Env struct {
	inner Bindings
	outer *Env
}

// NewEnv creates a new environment, optionally pre-populated with the given bindings.
func NewEnv(bindings ...Bindings) Env {
	store := make(map[string]Value)
	for _, b := range bindings {
		maps.Copy(store, b)
	}
	return Env{inner: store}
}

// Get returns the value bound to key, nil if the key is unbound.
func (e Env) Get(key string) Value {
	value, present := e.inner[key]
	if present {
		return value
	}

	for outer := e.outer; outer != nil; outer = outer.outer {
		value, present := outer.inner[key]
		if present {
			return value
		}
	}
	return nil
}

// Set binds key to val in the environment.
func (e Env) Set(key string, val Value) {
	e.inner[key] = val
}
