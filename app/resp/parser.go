package resp

import (
	"bytes"
	"errors"
	"io"
)

type RESP_TYPE byte

const (
	TYPE_SIMPLE_STRING RESP_TYPE = '+'
	TYPE_BULK_STRING   RESP_TYPE = '$'
	TYPE_ARRAY         RESP_TYPE = '*'
	TYPE_RDB           RESP_TYPE = 'R'
	TYPE_INTEGER       RESP_TYPE = ':'
	TYPE_UNKNOWN       RESP_TYPE = 0
)

func GetType(b byte) (RESP_TYPE, error) {
	switch b {
	case byte(TYPE_SIMPLE_STRING):
		return TYPE_SIMPLE_STRING, nil
	case byte(TYPE_BULK_STRING):
		return TYPE_BULK_STRING, nil
	case byte(TYPE_ARRAY):
		return TYPE_ARRAY, nil
	default:
		return 0, errors.New("Unknown type")
	}
}

// Parse parses the RESP input and returns the value and the rest of the input
func Parse(input []byte) (any, RESP_TYPE, []byte, error) {
	reader := bytes.NewReader(input)
	b, _ := reader.ReadByte()

	switch b {
	case '$':
		// parse bulk string
		size := get_bulk_string_size(reader)

		reader.ReadByte()
		reader.ReadByte()

		value := make([]byte, size)
		for i := 0; i < size; i++ {
			value[i], _ = reader.ReadByte()
		}

		t := TYPE_BULK_STRING // type of the data
		// check if the last two bytes are \r\n
		// if not, move the cursor back
		// INFO: for RDB files, the last two bytes are not \r\n
		if b, _ := reader.ReadByte(); b != '\r' {
			// is RDB file
			reader.Seek(-1, io.SeekCurrent)
			t = TYPE_RDB
		} else {
			// is bulk string
			reader.ReadByte()
		}

		rest := make([]byte, reader.Len())
		reader.Read(rest)

		return string(value), t, rest, nil
	case '*':
		// parse array
		value := make([]any, 0)
		size := get_array_size(reader)

		reader.ReadByte()
		reader.ReadByte()

		for i := 0; i < size; i++ {
			buffered := make([]byte, reader.Len())
			reader.Read(buffered)
			v, _, rest, _ := Parse(buffered)
			value = append(value, v)
			reader.Reset(rest)
		}

		rest := make([]byte, reader.Len())
		reader.Read(rest)

		return value, TYPE_ARRAY, rest, nil
	case '+':
		size := get_simple_string_size(reader)

		reader.Seek(1, io.SeekStart)
		// parse simple string
		value := make([]byte, size)
		reader.Read(value)

		reader.ReadByte()
		reader.ReadByte()

		rest := make([]byte, reader.Len())
		reader.Read(rest)

		return string(value), TYPE_SIMPLE_STRING, rest, nil
	}

	return nil, TYPE_UNKNOWN, make([]byte, 0), errors.New("Unknown type")
}
