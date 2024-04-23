package resp

import "fmt"

func EncodeFileContent(data []byte) []byte {
	size := len(data)
	return append([]byte("$"+fmt.Sprintf("%d", size)+CRLF), data...)
}
