package utl

import "strings"

func Unquote(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\\r\\n")
}
