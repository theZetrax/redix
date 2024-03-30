package internal

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal/parser"
)

// Request types
const (
	BULK_STRING = "$"
	ARRAY       = "*"
)

// Constants
const (
	PREFIX = "*1\r\n$4\r\n"
)

type Request struct {
	CMD parser.CMD
}

// Parses the incoming buffer and returns a Request object
func ParseRequest(buffer []byte) Request {
	raw := string(buffer)
	// remove the prefix and suffix from the raw request
	// the prefix and suffix are defined by the protocol
	// and are not part of the actual command
	var cmd parser.CMD
	if parsed, err := parser.ParseArray(raw); err != nil {
		log.Panicln("Error parsing request: ", err)
	} else {
		cmd, err = parser.NewCMD(parsed)
		if err != nil {
			log.Panicln("Error parsing request: ", err)
		}
	}

	return Request{
		CMD: cmd,
	}
}
