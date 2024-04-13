package resp

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Array struct {
	resp    Resp
	size    int
	Decoded []byte
	Parsed  []any
}

func GetArraySize(data []byte) int {
	data_str, _, _ := strings.Cut(string(data), CRLF)
	size, err := strconv.Atoi(data_str[1:])
	if err != nil {
		log.Panicln("Failed to parse array, invalid size: ", err)
	}
	return size
}

func EncodeArray(data ...[]byte) []byte {
	size := fmt.Sprintf("%d", len(data))
	encoded_array := []byte("*" + size + CRLF)
	for i, d := range data {
		encoded_array = append(encoded_array, d...)

		if i != len(data)-1 {
			// if encoded_array is not empty && doesn't already include CRLF
			if tail := len(encoded_array) - len(CRLF); len(encoded_array) > len(CRLF) {
				if string(encoded_array[tail:]) != CRLF {
					encoded_array = append(encoded_array, []byte(CRLF)...)
				}
			} else {
				encoded_array = append(encoded_array, []byte(CRLF)...)
			}
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

func (a *Array) decode() (data []byte, size int) {
	_, d, found := strings.Cut(string(a.resp.Data), CRLF)
	if !found {
		log.Panicln("Failed to parse array, invalid meta")
	}

	size = GetArraySize(a.resp.Data)
	data = []byte(d)

	a.Parsed = make([]any, size)
	a.Decoded = make([]byte, 0)

	for i := 0; i < size; i++ {
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

	return []byte(d), size
}

func (a *Array) Process() []byte {
	panic("should not call this method")
}

func (a *Array) String() string {
	return string(a.resp.Data)
}
