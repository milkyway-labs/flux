package utils

// Returns the elements of `slice` that satisfy `predicate`.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Converts the element of the provided `slice` using the provided `mapFunc`.
func Map[T any, R any](slice []T, mapFunc func(T) R) []R {
	result := make([]R, len(slice))
	for i, value := range slice {
		result[i] = mapFunc(value)
	}
	return result
}
