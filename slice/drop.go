package slice

// Drop remove n amount of item from slice
//   - start: start position
//   - n: number of item deleted
//   - slice: un-deleted slice
//
// start ：
//   - if start < 0,iterator from right to left,i.e -1 represent the last element of slice
func Drop[V any](start, n int, slice []V) []V {
	var s = make([]V, len(slice))
	copy(s, slice)
	if start < 0 {
		start = len(s) + start - n + 1
		if start < 0 {
			start = 0
		}
	}

	end := start + n
	if end > len(s) {
		end = len(s)
	}

	return append(s[:start], s[end:]...)
}
