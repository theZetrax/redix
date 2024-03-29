package internal

import (
	"log"
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
	CMD string
}

// Parses the incoming buffer and returns a Request object
func ParseRequest(buffer []byte) Request {
	raw := string(buffer)
	// remove the prefix and suffix from the raw request
	// the prefix and suffix are defined by the protocol
	// and are not part of the actual command
	log.Println("Raw request: ", raw)
	if parsed, err := parseRaw(raw); err != nil {
		log.Println("Error parsing request: ", err)
	} else {
		log.Println("Parsed request: ", parsed)
	}

	return Request{
		CMD: raw,
	}
}
