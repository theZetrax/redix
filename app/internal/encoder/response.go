// Description: Response encoder, encodes the response to the client
// Author: Zablon Dawit
// Date: Mar-30-2024
package encoder

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/parser"
)

func NewBulkString(value string) string {
	str := strings.Join([]string{
		string(parser.T_BULK_STRING) + fmt.Sprint(len(value)),
		value,
	}, parser.CRLF) + parser.CRLF
	return str
}

func NewError(value string) string {
	return string(parser.T_ERROR) + value + parser.CRLF
}
