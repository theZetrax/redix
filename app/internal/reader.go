package internal

import (
	"log"

	"github.com/codecrafters-io/redis-starter-go/app/internal/decoder"
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
	CMD decoder.CMD
}

// Parses the incoming buffer and returns a Request object
func ParseRequest(buffer []byte) (Request, error) {
	raw := string(buffer)

	// remove the prefix and suffix from the raw request
	// the prefix and suffix are defined by the protocol
	// and are not part of the actual command
	var cmd decoder.CMD
	if parsed, err := decoder.ParseArray(raw); err != nil {
		log.Println("Error parsing request: ", err)
		return Request{}, err
	} else {
		cmd, err = decoder.NewCMD(parsed)
		if err != nil {
			log.Println("Error parsing request: ", err)
			return Request{}, err
		}
	}

	return Request{
		CMD: cmd,
	}, nil
}
