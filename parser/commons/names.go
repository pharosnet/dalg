package commons

import "strings"

func NormalizeNameAndUpper(name string) string {
	return strings.ToUpper(strings.ReplaceAll(name, "`", ""))
}

func NormalizeName(name string) string {
	return strings.ReplaceAll(name, "`", "")
}

func SplitFullName(fullName string) (name1 string, name2 string) {
	fullName = NormalizeName(fullName)
	if strings.Index(fullName, ".") > 0 {
		names := strings.Split(fullName, ".")
		name1 = strings.TrimSpace(names[0])
		name2 = strings.TrimSpace(names[1])
	} else {
		name2 = fullName
	}
	return
}
