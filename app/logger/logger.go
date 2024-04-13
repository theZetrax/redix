package logger

import (
	"fmt"
	"strings"
)

func LogResp(message string, data []byte) {
	fmt.Println(message, strings.ReplaceAll(string(data), "\r\n", "\\r\\n"))
}
