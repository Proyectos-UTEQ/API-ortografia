package utils

func ContainsString(arr []string, str string) bool {
	for _, elemento := range arr {
		if elemento == str {
			return true
		}
	}
	return false
}
