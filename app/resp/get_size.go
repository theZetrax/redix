// Description: Functions to get the size of the data
package resp

import (
	"bytes"
	"io"
	"strconv"
)

// Reads the size encoded in the data from the reader
// moves the cursor until it finds the end of the size
// returns the size as an int64. The cursor is moved back
// to the last byte of the size
func get_array_size(reader *bytes.Reader) int {
	fb, _ := reader.ReadByte() // first byte
	size_bytes := append(make([]byte, 0), fb)
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		if _, err := strconv.ParseInt(string(b), 10, 0); err != nil {
			break
		} else {
			size_bytes = append(size_bytes, b)
		}
	}
	reader.Seek(int64(-1), io.SeekCurrent)
	size, _ := strconv.Atoi(string(size_bytes))

	return size
}

// Reads the size encoded in the data from the reader
// moves the cursor until it finds the end of the size
// returns the size as an int64. The cursor is moved back
func get_bulk_string_size(reader *bytes.Reader) int {
	fb, _ := reader.ReadByte() // first byte
	size_bytes := append(make([]byte, 0), fb)
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		if _, err := strconv.ParseInt(string(b), 10, 0); err != nil {
			break
		} else {
			size_bytes = append(size_bytes, b)
		}
	}
	reader.Seek(int64(-1), io.SeekCurrent)
	size, _ := strconv.Atoi(string(size_bytes))

	return size
}

func get_simple_string_size(reader *bytes.Reader) int {
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

	return size
}
