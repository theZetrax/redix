// Description: This file contains the implementation of the primitive encoder.
// Author: Zablon Dawit
// Date: Mar-30-2024
package encoder

import (
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
