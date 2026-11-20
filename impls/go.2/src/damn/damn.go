package damn

import (
	"errors"
	"fmt"

	"github.com/mooss/mal/src/ut"
)

// type Position struct {
// 	File      string
// 	Line, Col int
// }
// type CallSite struct {
// 	Func string
// 	Pos  Position
// }

type Error struct {
	cause  error
	prefix string    // Optional prefix qualifying the error.
	pcs    []uintptr // Go trace.
	// mals  []CallSite // Lisp trace, appended per Lisp call.
}

//////////////////////////////
// Standard error interface //

func (e *Error) Error() string {
	res := "<nil>!!!"
	if e.cause != nil {
		res = e.cause.Error()
	}
	if e.prefix != "" {
		res = e.prefix + ": " + res
	}

	return res
}

func (e *Error) Unwrap() error { return e.cause }

//////////////////////////
// Casts and predicates //

// Cast attempts a direct cast of err to E, also trying with the cause if err is an *Error.
func Cast[E error](err error) (E, bool) {
	//nolint:errorlint
	if typed, isE := err.(E); isE {
		return typed, true
	}

	//nolint:errorlint
	damn, isDamn := err.(*Error)
	if !isDamn {
		var zero E
		return zero, false
	}

	//nolint:errorlint
	res, isE := damn.cause.(E)
	return res, isE
}

// Is returns true if the exact type of err is E or if it is an *Error whose direct cause is of type
// E.
func Is[E error](err error) bool {
	if ut.Is[E](err) {
		return true
	}

	//nolint:errorlint
	damn, isDamn := err.(*Error)
	if !isDamn {
		return false
	}

	return ut.Is[E](damn.cause)
}

////////////////////////
// Error construction //

func From(err error) error {
	if err == nil {
		return nil
	}
	return &Error{cause: err, pcs: capture(0)}
}

func New(msg string) error { return &Error{cause: errors.New(msg), pcs: capture(0)} }

func Newf(format string, args ...any) error {
	return &Error{cause: fmt.Errorf(format, args...), pcs: capture(0)}
}

// Qualify returns a new error with an additional prefix added to the error message.
// Preserves nil errors.
func (e *Error) Qualify(prefix string) *Error {
	if e == nil {
		return nil
	}
	if e.prefix != "" {
		prefix = prefix + ": " + e.prefix
	}
	return &Error{cause: e.cause, prefix: prefix, pcs: e.pcs}
}

// Qualify returns a new error with an additional prefix added to the error message.
// Preserves nil errors.
func Qualify(err error, prefix string) error {
	if err == nil {
		return nil
	}

	//nolint:errorlint
	if damn, isDamn := err.(*Error); isDamn {
		return damn.Qualify(prefix)
	}

	return &Error{cause: err, prefix: prefix, pcs: capture(0)}
}

/////////////////
// Local utils //

// capture captures the caller stack.
// Disabled for now until figuring out how to use it and when not to use it (in particular in
// relation to Exceptions, which should not capture).
func capture(int) []uintptr {
	return nil
	// pcs := make([]uintptr, 32)
	// n := runtime.Callers(3+skip, pcs)
	// return pcs[:n]
}
