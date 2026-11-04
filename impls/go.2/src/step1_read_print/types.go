package main

type TypeName string

const (
	INTEGER TypeName = "integer"
	SYMBOL  TypeName = "symbol"
	LIST    TypeName = "list"
)

type MalType interface {
	Type() TypeName
}

type Integer int64
type Symbol string
type List []MalType

func (Integer) Type() TypeName { return INTEGER }
func (Symbol) Type() TypeName  { return SYMBOL }
func (List) Type() TypeName    { return LIST }

func (lst List) Map[To any](fun func(MalType) To) []To {
	res := make([]To, len(lst))
	for i, el := range lst {
		res[i] = fun(el)
	}
	return res
}
