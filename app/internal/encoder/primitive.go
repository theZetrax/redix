// Description: This file contains the implementation of the primitive encoder.
// Author: Zablon Dawit
// Date: Mar-30-2024
package encoder

import (
	"crypto/sha1"
	"fmt"
)

// ConvertSliceToStringArray converts a slice of unknown types to a slice of strings.
func ConvertSliceToStringArray(slice []interface{}) []string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = fmt.Sprintf("%v", v)
	}
	return strSlice
}

func Sha1Hash(input string) string {
	h := sha1.New()
	h.Write([]byte(input))
	hashed := h.Sum(nil)

	return fmt.Sprintf("%x", hashed)
}
