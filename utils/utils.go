package utils

func FixNilArrayToEmpty[T any](a []T) []T {
	if a == nil {
		return []T{}
	}

	return a
}
