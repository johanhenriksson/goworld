package util

func MapValues[K comparable, V any, T any](m map[K]V, transform func(V) T) []T {
	output := make([]T, len(m))
	i := 0
	for _, value := range m {
		output[i] = transform(value)
		i++
	}
	return output
}

func MapKeys[K comparable, V any, T any](m map[K]V, transform func(K) T) []T {
	output := make([]T, len(m))
	i := 0
	for key := range m {
		output[i] = transform(key)
		i++
	}
	return output
}
