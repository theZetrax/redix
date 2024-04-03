// Description: Complex structure encoder, encodes the complex structures spicific to RESP
// to the client.
// Author: Zablon Dawit
// Date: Mar-30-2024
package encoder

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
)

func NewBinaryString(value []byte) string {
	str := strings.Join([]string{
		string(decoder.T_BULK_STRING) + fmt.Sprint(len(value)),
		string(value),
	}, decoder.CRLF)
	return str
}

func NewBulkString(value string) string {
	str := strings.Join([]string{
		string(decoder.T_BULK_STRING) + fmt.Sprint(len(value)),
		value,
	}, decoder.CRLF) + decoder.CRLF
	return str
}

func NewError(err error) string {
	return string(decoder.T_ERROR) + err.Error() + decoder.CRLF
}

func NewSimpleString(value string) string {
	return string(decoder.T_SIMPLE_STRING) + value + decoder.CRLF
}

func NewNil() string {
	return string(decoder.T_BULK_STRING) + fmt.Sprint(-1) + decoder.CRLF
}

func NewArray(entries ...string) (raw string) {
	size := fmt.Sprint(len(entries))
	raw = string(decoder.T_ARRAY) + size + decoder.CRLF
	for _, entry := range entries {
		hasSuffix := strings.HasSuffix(entry, decoder.CRLF)
		raw += entry
		if !hasSuffix {
			raw += decoder.CRLF
		}
	}

	return raw // two CRLF to end the array
}
