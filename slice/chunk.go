package slice

import "errors"

// Chunk split origin slice into even slice
func Chunk[V any](slice []V, size int) ([][]V, error) {
	if len(slice) == 0 {
		return nil, errors.New("empty collection")
	}

	if size < 1 {
		return nil, errors.New("size must be greater than 0")
	}

	result := make([][]V, 0, (len(slice)+size-1)/size)
	for size < len(slice) {
		slice, result = slice[size:], append(result, slice[0:size])
	}
	return append(result, slice), nil
}

// ChunkV2 faster version split origin slice into even slice
func ChunkV2[V any](slice []V, size int) ([][]V, error) {
	if len(slice) == 0 {
		return nil, errors.New("empty slice")
	}
	divided := make([][]V, (len(slice)+size-1)/size)
	prev := 0
	i := 0
	till := len(slice) - size
	for prev < till {
		next := prev + size
		divided[i] = slice[prev:next]
		prev = next
		i++
	}
	divided[i] = slice[prev:]
	return divided, nil
}
