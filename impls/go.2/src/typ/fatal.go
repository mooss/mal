package typ

import (
	"fmt"

	"github.com/mooss/mal/src/core"
)

///////////
// Fatal //

// Fatal is an error that interrupts the flow of execution.
// It represent two distincts sets of error that both interrupt the execution:
//  1. An irrecoverable internal error (when error is set).
//  2. An exception that can be recovered from with try*/catch* (when error is not set).
type Fatal struct {
	error
	Exception core.Value
}

// Fatalf creates a Fatal error with a formatted message.
func Fatalf(msg string, args ...any) Fatal { return Fatal{error: fmt.Errorf(msg, args...)} }

// NewFatal creates a fatal error wrapping the given source error.
func NewFatal(source error) Fatal { return Fatal{error: source} }

func Exception(val core.Value) Fatal { return Fatal{Exception: val} }

func Exceptionf(format string, args ...any) Fatal {
	return Fatal{Exception: String(fmt.Sprintf(format, args...))}
}

// Equal returns true if both are nil errors and both their exception are equal.
func (f Fatal) Equal(other core.Value) bool {
	o, ok := other.(Fatal)
	if !ok || f.error != nil || o.error != nil {
		return false
	}

	if f.Exception == nil || o.Exception == nil {
		return f.Exception == o.Exception
	}

	return f.Exception.Equal(o.Exception)
}

func (f Fatal) Error() string {
	if err := f.Unbox(); err != nil {
		return err.Error()
	}
	return ""
}

func (f Fatal) Eval(core.Env) core.Value { return f }

// NotNil returns true when the embedded error is not nil or the embedded exception is not nil.
func (f Fatal) NotNil() bool { return f.error != nil || f.Exception != nil }

// Qualify returns a new fatal error with "{prefix}: " added in front of the error.
// When the error was a recoverable exception, it gets upgraded to a irrecoverable error.
func (f Fatal) Qualify(prefix string) Fatal {
	return Fatal{error: fmt.Errorf("%s: %w", prefix, f.Unbox()), Exception: f.Exception}
}

func (f Fatal) String() string {
	if f.error != nil {
		return fmt.Sprintf("(fatal %q)", f.Error())
	}
	return fmt.Sprintf("Exception %s", f.Exception)
}

// Recoverable returns true, when this error can be recovered from via the try*/catch* mechanism.
func (f Fatal) Recoverable() bool { return f.error == nil }

func (Fatal) Type() core.TypeName { return core.FATAL }

func (f Fatal) Unbox() error {
	if f.error == nil && f.Exception != nil {
		return fmt.Errorf("(throw %s)", f.Exception)
	}
	return f.error
}
