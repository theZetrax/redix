package resp

import "fmt"

func EncodeInteger(val int) []byte {
	return []byte(string(TYPE_INTEGER) + fmt.Sprintf("%d", val) + CRLF)
}
