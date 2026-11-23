package util

func TruncateString(str string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}
	runes := []rune(str)
	if len(runes) <= maxLength {
		return str
	}

	return string(runes[:maxLength])
}