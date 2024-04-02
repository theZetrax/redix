// Description: Complex structure encoder, encodes the complex structures spicific to RESP
// to the client.
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

func NewError(err error) string {
	return string(parser.T_ERROR) + err.Error() + parser.CRLF
}

func NewSimpleString(value string) string {
	return string(parser.T_SIMPLE_STRING) + value + parser.CRLF
}

func NewNil() string {
	return string(parser.T_BULK_STRING) + fmt.Sprint(-1) + parser.CRLF
}

func NewArray(entries ...string) (raw string) {
	size := fmt.Sprint(len(entries))
	raw = string(parser.T_ARRAY) + size + parser.CRLF
	for _, entry := range entries {
		hasSuffix := strings.HasSuffix(entry, parser.CRLF)
		raw += entry
		if !hasSuffix {
			raw += parser.CRLF
		}
	}

	return raw // two CRLF to end the array
}
