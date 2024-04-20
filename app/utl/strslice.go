package utl

import "fmt"

// ConvertSliceToStringArray converts a slice of unknown types to a slice of strings.
func ToStringSlice(slice []any) []string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprintf("%v", v)
	}
	return strSlice
}
