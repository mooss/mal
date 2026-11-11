package core

import (
	"iter"
	"maps"
	"slices"

	"github.com/mooss/mal/src/fun"
	"github.com/mooss/mal/src/lazy"
)

// Bindings is a primitive environment associating symbols to their values.
type Bindings map[string]Value

func (b Bindings) AsEnv() Env {
	if b == nil {
		b = map[string]Value{}
	}
	return Env{inner: b}
}

func (b Bindings) Len() int { return len(b) }

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
		for key, value := range env.IterAssocs() {
			if _, exists := res.inner[key]; !exists {
				res.inner[key] = value
			}
		}
	}
	return res
}

// Get returns the value bound to key, nil if the key is unbound.
func (e Env) Get(key string) Value {
	for inner := range e.IterBindings() {
		value, present := inner[key]
		if present {
			return value
		}
	}
	return nil
}

// IterAssocs returns an iterator on all the key-value associations defined in this environment,
// including those that are shadowed.
// The iteration starts with inner bindings.
func (e Env) IterAssocs() iter.Seq2[string, Value] {
	return func(yield func(string, Value) bool) {
		for bindings := range e.IterBindings() {
			for key, value := range bindings {
				if !yield(key, value) {
					return
				}
			}
		}
	}
}

// IterBindings returns an iterator on all the inner bindings defined in this environment, starting
// with the innermost bindings.
func (e Env) IterBindings() iter.Seq[Bindings] {
	return func(yield func(Bindings) bool) {
		for penv := &e; penv != nil; penv = penv.outer {
			if !yield(penv.inner) {
				return
			}
		}
	}
}

// NAssocs returns the total number of associations in the environment, including those that are
// shadowed.
func (e Env) NAssocs() int {
	lengths := lazy.Map(Bindings.Len, e.IterBindings())
	return fun.Reduce(func(l, r int) int { return l + r }, 0, lengths)
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
