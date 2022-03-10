package util

import "strings"

func CStrings(strings []string) []string {
	return Map(strings, func(str string) string {
		return CString(str)
	})
}

func CString(str string) string {
	if strings.HasSuffix(str, "\x00") {
		return str
	}
	return str + "\x00"
}
