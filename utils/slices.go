package utils

func Supported(supports []string, target string) bool {
	for _, value := range supports {
		if value == "*" || value == target {
			return true
		}
	}
	return false
}
