package assert

import (
	"fmt"
	"reflect"
	"testing"
)

// Tester holds a *testing.T and provides methods wrapping assert functions.
// Every function in this package whose first argument is a *testing.T is registered as a Tester
// method, forwarding T and its arguments to the function.
//
// For instance, the function:
//
//	func Must[T any](t *testing.T, v T, err error) T
//
// Is registered as:
//
//	func (w Tester) Must[T any](v T, err error) T { return Must(w.T, v, err) }
type Tester struct {
	T *testing.T
}

// New creates a new wrapper holding t.
func New(t *testing.T) Tester { return Tester{t} }

func (tr Tester) RunCases[Case any](
	groupname string, fun func(Tester, Case), cases ...Case,
) {
	tr.T.Helper() // REVIEW: Not sure if this should be here (or the other T.Helper()).
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
				fun(Tester{t}, c)
			})
		}
	})
}

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
