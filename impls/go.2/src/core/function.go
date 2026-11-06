package core

type Function func(Integer, Integer) Integer

func (Function) Type() TypeName { return FUNCTION }
func (Function) String() string { return "<function>" }
