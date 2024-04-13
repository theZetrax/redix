package resp

import (
	"strings"
)

type SimpleString struct {
	resp    Resp
	Parsed  string // Data parsed to simple string
	Decoded []byte // Decoded data
}

func EncodeSimpleString(str string) []byte {
	return []byte("+" + str + CRLF)
}

func NewSimpleString(data []byte) *SimpleString {
	s := &SimpleString{
		resp: Resp{
			Data: data,
		},
	}
	s.decode()

	return s
}

func (s *SimpleString) decode() string {
	data := s.resp.Data
	data = data[1:]
	s.Parsed = strings.TrimRight(string(data), CRLF)
	s.Decoded = EncodeSimpleString(s.Parsed)

	return s.Parsed
}
func (s *SimpleString) Process() []byte {
	switch str := s.String(); str {
	case "PING":
		return EncodeSimpleString("PONG")
	default:
		return EncodeSimpleString(s.Parsed)
	}
}

func (s *SimpleString) String() string {
	s.Parsed = strings.ToUpper(s.decode())
	return s.Parsed
}
