package util

func Map[T any, S any](items []T, transform func(T) S) []S {
	output := make([]S, len(items))
	for i, item := range items {
		output[i] = transform(item)
	}
	return output
}

func MapIdx[T any, S any](items []T, transform func(T, int) S) []S {
	output := make([]S, len(items))
	for i, item := range items {
		output[i] = transform(item, i)
	}
	return output
}

func Range(from, to, step int) []int {
	n := (to - from) / step
	output := make([]int, n)
	v := from
	for i := 0; v < to; i++ {
		output[i] = v
		v += step
	}
	return output
}

func Chunks[T any](slice []T, size int) [][]T {
	count := len(slice) / size
	chunks := make([][]T, 0, count)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func Reduce[T any, S any](slice []T, initial S, reducer func(S, T) S) S {
	accumulator := initial
	for _, item := range slice {
		accumulator = reducer(accumulator, item)
	}
	return accumulator
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
	output := make([]T, 0, len(slice))
	for _, item := range slice {
		if predicate(item) {
			output = append(output, item)
		}
	}
	return output
}

func Contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
