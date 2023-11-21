package slice

// FillBy fill value into slice
func FillBy[V any](slice []V, fn func(index int, val V) V) []V {
	for i, v := range slice {
		slice[i] = fn(i, v)
	}
	return slice
}

// FillByCopy same as FillBy but will not change the origin slice
func FillByCopy[V any](slice []V, fn func(index int, val V) V) []V {
	var s = make([]V, len(slice))
	copy(s, slice)
	return FillBy(s, fn)
}
