// This file implements trivial types, i.e. types whose representation and evaluation amounts to
// basically nothing.

package core

import (
	"fmt"
	"strconv"
)

type Error struct{ error }

func Errorf(msg string, args ...any) Error { return Error{error: fmt.Errorf(msg, args...)} }
func (Error) Type() TypeName               { return ERROR }
func (e Error) String() string             { return fmt.Sprintf("(error %q)", e.Error()) }

type Integer int64

func (Integer) Type() TypeName   { return INTEGER }
func (i Integer) String() string { return strconv.FormatInt(int64(i), 10) }

type Nil struct{}

func (Nil) Type() TypeName { return NIL }
func (Nil) String() string { return "nil" }

type String string

func (String) Type() TypeName   { return STRING }
func (s String) String() string { return strconv.Quote(string(s)) }

type Symbol string

func (Symbol) Type() TypeName   { return SYMBOL }
func (s Symbol) String() string { return string(s) }
