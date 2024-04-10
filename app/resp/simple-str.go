package resp

import (
	"fmt"
	"strings"
)

type SimpleString struct {
	resp       Resp
	ParsedData string
}

func EncodeSimpleString(str string) []byte {
	return []byte(fmt.Sprintf("+%s"+CRLF, str))
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
	return strings.TrimRight(string(data), CRLF)
}
func (s *SimpleString) Process() []byte {
	fmt.Println("Data: ", s.decode())
	switch s.ParsedData = strings.ToUpper(s.decode()); s.ParsedData {
	case "PING":
		return EncodeSimpleString("PONG")
	default:
		return EncodeSimpleString(s.ParsedData)
	}
}
