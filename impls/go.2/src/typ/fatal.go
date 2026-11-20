package typ

import (
	"fmt"

	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/damn"
)

///////////
// Fatal //

// Fatal is an error that interrupts the flow of execution.
type Fatal struct {
	Wrapped error
}

// Fatalf creates a Fatal error with a formatted message.
func Fatalf(msg string, args ...any) Fatal { return Fatal{damn.Newf(msg, args...)} }

// NewFatal creates a fatal error wrapping the given source error.
func NewFatal(source error) Fatal { return Fatal{damn.From(source)} }

func Exception(val core.Value) Fatal { return Fatal{damn.Throw(val)} }

func Exceptionf(format string, args ...any) Fatal {
	return Fatal{damn.Throw(String(fmt.Sprintf(format, args...)))}
}

// Fatal errors are never equal.
func (Fatal) Equal(core.Value) bool { return false }

func (f Fatal) Error() string {
	if f.Wrapped == nil {
		return "<nil>"
	}
	if exc, isExc := damn.Cast[*damn.Exception](f.Wrapped); isExc {
		return "Exception: " + exc.Payload.String()
	}
	return "Fatal error: " + f.Wrapped.Error()
}

func (f Fatal) Eval(core.Env) core.Value { return f }

// NotNil returns true when the embedded error is not nil.
func (f Fatal) NotNil() bool { return f.Wrapped != nil }

// Qualify returns a new fatal error with "{prefix}: " added in front of the error.
func (f Fatal) Qualify(prefix string) Fatal {
	return Fatal{damn.Qualify(f.Wrapped, prefix)}
}

func (f Fatal) String() string {
	if f.Wrapped != nil {
		return fmt.Sprintf("(fatal %q)", f.Wrapped.Error())
	}
	return "(fatal)"
}

func (Fatal) Type() core.TypeName { return core.FATAL }

func (f Fatal) Unbox() error { return f.Wrapped }
