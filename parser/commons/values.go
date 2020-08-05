package commons

import "strings"

func NormalizeUpperValue(v string) string {
	return strings.TrimSpace(strings.ToUpper(strings.ReplaceAll(v, "'", "")))
}

func NormalizeValue(v string) string {
	return strings.TrimSpace(strings.ReplaceAll(v, "'", ""))
}
