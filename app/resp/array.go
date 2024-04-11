package resp

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/logger"
)

type Array struct {
	resp    Resp
	size    int
	Decoded []byte
	Parsed  []any
}

func EncodeArray(data ...[]byte) []byte {
	size := fmt.Sprintf("%d", len(data))
	encoded_array := []byte("*" + size + CRLF)
	for i, d := range data {
		encoded_array = append(encoded_array, d...)

		if i != len(data)-1 {
			encoded_array = append(encoded_array, []byte(CRLF)...)
		}
	}
	return encoded_array
}

func NewArray(data []byte) *Array {
	arr := &Array{
		resp: Resp{
			Data: data,
		},
	}
	arr.decode()

	return arr
}

func GetArraySize(data []byte) int {
	size, err := strconv.Atoi(string(data[1:2]))
	if err != nil {
		log.Panicln("Failed to parse array, invalid size: ", err)
	}
	return size
}

func (a *Array) decode() (data []byte, size int) {
	_, d, found := strings.Cut(string(a.resp.Data), CRLF)
	if !found {
		log.Panicln("Failed to parse array, invalid meta")
	}

	size = GetArraySize(a.resp.Data)
	data = []byte(d)

	log.Println("Decoded Size: ", size)
	logger.LogResp("Decoded Data: ", []byte(data))

	a.Parsed = make([]any, size)
	a.Decoded = make([]byte, 0)

	for i := 0; i < size; i++ {
		logger.LogResp("Data: ", data)

		switch handler, _ := HandleResp(data); handler.(type) {
		case *BulkString:
			bs := handler.(*BulkString)
			data = data[len(bs.Decoded):]
			a.Parsed[i] = bs.Parsed
			a.Decoded = append(a.Decoded, bs.Decoded...)
		case *Array:
			arr := handler.(*Array)
			data = data[len(arr.Decoded):]
			a.Parsed[i] = arr.Parsed
			a.Decoded = append(a.Decoded, arr.Decoded...)
		case *SimpleString:
			str := handler.(*SimpleString)
			data = data[len(str.Decoded):]
			a.Parsed[i] = str.Parsed
			a.Decoded = append(a.Decoded, str.Decoded...)
		}
	}

	log.Println(a.Parsed)
	logger.LogResp("Decoded Data: ", a.Decoded)

	return []byte(d), size
}

func (a *Array) Process() []byte {
	panic("should not call this method")
}

func (a *Array) String() string {
	return string(a.resp.Data)
}
