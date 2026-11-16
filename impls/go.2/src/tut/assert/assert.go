package assert

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/tut"
)

// T holds a *testing.T and provides methods wrapping assert functions.
// Every function in this package whose first argument is a *testing.T is registered as a T
// method, forwarding T and its arguments to the function.
//
// For instance, the function:
//
//	func Must[T any](t *testing.T, v T, err error) T
//
// Is registered as:
//
//	func (w T) Must[T any](v T, err error) T { return Must(w.T, v, err) }
type T struct {
	T *testing.T
}

// New creates a new wrapper holding t.
func New(t *testing.T) T { return T{t} }

func (T) MustEval(env core.Env, program string) core.Value {
	return tut.MustEval(env, program)
}

func (tr T) MustEvalAs[T core.Value](env core.Env, program string) T {
	tr.T.Helper()
	return tut.MustEvalAs[T](tr.T, env, program)
}

func (tr T) RunCases[Case any](
	groupname string, fun func(T, Case), cases ...Case,
) {
	if len(cases) == 0 {
		tr.T.Fatalf("%s was given no test cases", groupname)
	}

	tr.T.Run(groupname, func(t *testing.T) {
		t.Helper()
		for num, c := range cases {
			name := fmt.Sprintf("n°%d", num)
			if pname := getStringField("Name", c); pname != nil && *pname != "" {
				name = *pname
			}

			t.Run(name, func(t *testing.T) {
				t.Helper()
				fun(New(t), c)
			})
		}
	})
}

//////////////////////////
// testing.T forwarding //

func (tr T) Fatalf(format string, args ...any) {
	tr.T.Fatalf(format, args...)
}

func (tr T) Run(name string, fun func(T)) {
	tr.T.Run(name, func(t *testing.T) {
		t.Helper()
		fun(New(t))
	})
}

func (tr T) TempDir() string {
	return tr.T.TempDir()
}

/////////////////
// Local utils //

func getStringField(field string, data any) *string {
	v := reflect.Indirect(reflect.ValueOf(data))
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return nil
	}

	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.String {
		return nil
	}

	res := f.String()
	return &res
}
