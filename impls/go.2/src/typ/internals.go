package typ

func _map[From, To any](fun func(From) To, source []From) []To {
	res := make([]To, len(source))
	for i, el := range source {
		res[i] = fun(el)
	}
	return res
}

func _mapflat[From comparable, To any](fun func(From) To, source map[From]From) []To {
	res := make([]To, len(source)*2)
	i := 0
	for k, v := range source {
		res[i] = fun(k)
		i++
		res[i] = fun(v)
		i++
	}
	return res
}
