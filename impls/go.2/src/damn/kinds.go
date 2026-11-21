package damn

import (
	"fmt"

	"github.com/mooss/mal/src/core"
)

type Throwable interface {
	error
	Mal() core.Value
}

///////////////
// Exception //

type Exception struct{ Payload core.Value }

func (e *Exception) Error() string {
	if e.Payload == nil {
		return "<nil>"
	}
	return e.Payload.String()
}

func (e *Exception) Mal() core.Value { return e.Payload }
func Throw(v core.Value) error       { return &Exception{v} }

///////////
// Arity //

type ArityError struct {
	Want string
	Got  int
}

func (e *ArityError) Error() string {
	return fmt.Sprintf("expected %s args, got %d", e.Want, e.Got)
}

//////////
// Type //

type TypeError struct {
	Want string
	Got  core.Value
}

func (e *TypeError) Error() string {
	if e.Got == nil {
		return fmt.Sprintf("expected %s but got nil (nil)", e.Want)
	}
	return fmt.Sprintf("expected %s but got %s (%v)", e.Want, e.Got.Type(), e.Got)
}

func TypeMismatch[T ~string](want T, got core.Value) error {
	return &TypeError{string(want), got}
}

////////////
// Bounds //

type BoundsError struct {
	Index int
	Len   int
}

func (e *BoundsError) Error() string {
	return fmt.Sprintf("index %d out of bounds for slice of size %d", e.Index, e.Len)
}

////////////
// Syntax //

type SyntaxError struct {
	File      string
	Line, Col int
	Msg       string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", e.File, e.Line, e.Col, e.Msg)
}
