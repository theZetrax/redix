package resp

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Array struct {
	resp    Resp
	size    int
	Decoded []byte // decoded data, raw
	Parsed  []any  // parsed data, into SimpleString, BulkString, or Array
}

func GetArraySize(data []byte) int {
	data_str, _, _ := strings.Cut(string(data), CRLF)
	size, err := strconv.Atoi(data_str[1:])
	if err != nil {
		log.Panicln("Failed to parse array, invalid size: ", err)
	}
	return size
}

func IsValidArray(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	str := string(data)
	meta, parts, found := strings.Cut(str, CRLF)
	if !found || meta[0] != '*' {
		return false
	}

	size, err := strconv.Atoi(meta[1:])
	if err != nil {
		log.Println("Failed to parse array, invalid size: ", err)
		return false
	}
	actual_parts := strings.Split(parts, CRLF)

	// if size is less than the actual parts
	return size < len(actual_parts)
}

// Check if the array data is a nested array
func IsNestedArray(data []any) bool {
	for _, d := range data {
		if reflect.TypeOf(d).Kind() == reflect.Slice {
			return true
		}
	}
	return false
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

// Would parse the array data into `[]any` type data.
// if there are multiple array data encoded in the bytes
// the result would be an `[][]any` type data
func NewArray(data []byte) *Array {
	arr := &Array{
		resp: Resp{
			Data: data,
		},
	}
	arr.decode()

	return arr
}

func (a *Array) decode() ([]byte, int) {
	data := a.resp.Data

	size := 0

	a.Parsed = make([]any, 0)
	a.Decoded = make([]byte, 0)
	a.size = size

	for {
		_, d, found := strings.Cut(string(data), CRLF)
		if !found {
			log.Panicln("Failed to parse array, invalid meta")
		}
		current_segment_size := GetArraySize(data)
		size += current_segment_size
		data = []byte(d) // take the data after the metadata
		parsed_collection := make([]any, 0)
		for i := 0; i < current_segment_size; i++ {
			switch handler, _ := HandleResp(data); handler.(type) {
			case *BulkString:
				bs := handler.(*BulkString) // bulk-string
				data = data[len(bs.Decoded):]
				parsed_collection = append(parsed_collection, bs.Parsed)
				a.Decoded = append(a.Decoded, bs.Decoded...)
			// case *Array:
			// 	arr := handler.(*Array)
			// 	data = data[len(arr.Decoded):]
			// 	parsed_collection = append(parsed_collection, arr.Parsed)
			// 	a.Decoded = append(a.Decoded, arr.Decoded...)
			case *SimpleString:
				str := handler.(*SimpleString)
				data = data[len(str.Decoded):]
				parsed_collection = append(parsed_collection, str.Parsed)
				a.Decoded = append(a.Decoded, str.Decoded...)
			}
			// a.Parsed = append(a.Parsed, parsed_collection)
			// check if there is more array left in the data
			//	if there is more array left, then check if it's valid
			//	if there is more array left, the result would be an `[][]any` type data
			// if len(data) == 0 || !IsValidArray(data) {
			// 	break
			// }
		}

		a.Parsed = append(a.Parsed, parsed_collection)

		if IsValidArray(data) {
			continue
		}
		break
	}

	if len(a.Parsed) == 1 {
		a.Parsed = a.Parsed[0].([]any)
	}

	return a.Decoded, size
}

func (a *Array) Process() []byte {
	panic("should not call this method, should be processed by cmd.NewCMD")
}

func (a *Array) String() string {
	return string(a.resp.Data)
}
