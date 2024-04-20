package resp

import (
	"bytes"
	"errors"
	"io"
)

func parse(input []byte) (any, []byte, error) {
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

		// check if the last two bytes are \r\n
		// if not, move the cursor back
		if b, _ := reader.ReadByte(); b != '\r' {
			reader.Seek(-1, io.SeekCurrent)
		} else {
			reader.ReadByte()
		}

		rest := make([]byte, reader.Len())
		reader.Read(rest)

		return string(value), rest, nil
	case '*':
		// parse array
		value := make([]any, 0)
		size := get_array_size(reader)

		reader.ReadByte()
		reader.ReadByte()

		for i := 0; i < size; i++ {
			buffered := make([]byte, reader.Len())
			reader.Read(buffered)
			v, rest, _ := parse(buffered)
			value = append(value, v)
			reader.Reset(rest)
		}

		rest := make([]byte, reader.Len())
		reader.Read(rest)

		return value, rest, nil
	case '+':
		size := 0

		for {
			b, _ := reader.ReadByte()
			size += 1

			if b == '\r' {
				b, _ = reader.ReadByte()
				size += 1

				if b == '\n' {
					break
				}
			}

			if reader.Len() == 0 {
				break
			}
		}
		size -= 2 // for the last \r\n

		reader.Seek(1, io.SeekStart)
		// parse simple string
		value := make([]byte, size)
		reader.Read(value)

		reader.ReadByte()
		reader.ReadByte()

		rest := make([]byte, reader.Len())
		reader.Read(rest)

		return string(value), rest, nil
	}

	return nil, make([]byte, 0), errors.New("Unknown type")
}
