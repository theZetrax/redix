package resp

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type BulkString struct {
	resp    Resp
	Parsed  string // Data parsed to bulk string
	Decoded []byte // Decoded data
}

func EncodeBulkString(str string) []byte {
	size := fmt.Sprintf("%d", len(str))
	return []byte("$" + size + CRLF + str + CRLF)
}

func NewBulkString(data []byte) *BulkString {
	bs := &BulkString{
		resp: Resp{
			Data: data,
		},
	}
	bs.decode()

	return bs
}

func GetBulkStringSize(data []byte) int {
	size, err := strconv.Atoi(string(data[1:2]))
	if err != nil {
		log.Panicln("Failed to parse bulk string, invalid size: ", err)
	}
	return size
}

func (b *BulkString) decode() string {
	data := b.resp.Data
	_, str, found := strings.Cut(string(data), CRLF)
	if !found {
		log.Panicln("Failed to parse bulk string")
	}

	size := GetBulkStringSize([]byte(data))

	b.Parsed = str[:size]
	b.Decoded = EncodeBulkString(b.Parsed)
	return string(b.Decoded)
}

func (b *BulkString) Process() []byte {
	return EncodeBulkString(b.Parsed)
}

func (b *BulkString) String() string {
	return b.Parsed
}
