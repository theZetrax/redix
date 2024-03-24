package internal

import "strings"

const (
	CLRF   = "\r\n"
	PREFIX = "*1\r\n$4\r\n"
)

type Request struct {
	CMD string
}

// Parses the incoming buffer and returns a Request object
func ParseRequest(buffer []byte) Request {
	raw := string(buffer)
	// remove the prefix and suffix from the raw request
	// the prefix and suffix are defined by the protocol
	// and are not part of the actual command
	raw = strings.TrimPrefix(raw, PREFIX)
	raw = strings.TrimSuffix(raw, CLRF+CLRF)
	raw = strings.ToUpper(strings.TrimSpace(raw))

	return Request{
		CMD: raw,
	}

}
