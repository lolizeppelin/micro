package utils

func Supported(supports []string, target string) bool {
	for _, value := range supports {
		if value == "*" || value == target {
			return true
		}
	}
	return false
}

func SliceReverse[T any](s []T) []T {
	newS := make([]T, len(s))
	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		newS[i], newS[j] = s[j], s[i]
	}
	return newS
}
