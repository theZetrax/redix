// @author Zablon Dawit <zablon@qebero.dev>
// @date 2021/08/01
package internal

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

const (
	CRLF = "\r\n"
)

// Types
const (
	T_SIMPLE_STRING = '+'
	T_BULK_STRING   = '$'
	T_ARRAY         = '*'
	T_INTEGER       = ':'
	T_INVALID       = '0'
)

func parseRaw(raw string) (any, error) {
	data_type := raw[0]
	switch data_type {
	case T_SIMPLE_STRING:
		return parseSimpleString(raw)
	case T_BULK_STRING:
		return parseBulkString(raw)
	case T_ARRAY:
		return parseArray(raw)
	case T_INTEGER:
		return parseInt(raw)
	default:
		return nil, errors.New("Invalid raw string")
	}
}

func parseArray(raw string) ([]any, error) {
	size, err := strconv.Atoi(raw[1:2])
	// split meta & data
	_, data, found := strings.Cut(raw, CRLF)
	if !found || err != nil {
		return nil, errors.New("Invalid Array")
	}

	entries := make([]any, 0)

	// split the array into individual elements
	for i := 0; i < size; i++ {
		if dataType, hasMeta, err := checkDataType(data); err != nil {
			return nil, err // exit `parseArray` on error
		} else {
			if hasMeta {
				entry := []any{}
				entry_meta, entry_data_with_rest, found := strings.Cut(data, CRLF)
				if !found {
					return nil, errors.New("Invalid Array")
				}

				entry_data, rest, found_data := strings.Cut(entry_data_with_rest, CRLF)
				if !found_data {
					entry_data = entry_data_with_rest
				} else {
					entry_data = data
				}

				switch dataType {
				case T_BULK_STRING:
					parsed_entry, err := parseBulkString(entry_data)
					if err != nil {
						return nil, err
					}
					entries = append(entries, parsed_entry)
					break
				case T_ARRAY:
					parsed_entry, err := parseArray(entry_data)
					if err != nil {
						return nil, err
					}
					entries = append(entries, parsed_entry)
					break
				default:
					entry = append(entry, entry_meta, entry_data)
					log.Println("Entry: ", entry)
					entries = append(entries, entry)
					break
				}

				if found {
					data = rest
				}
			} else {
				entry, rest, found := strings.Cut(data, CRLF)
				if !found {
					return nil, errors.New("Invalid Entry")
				}

				// Parse the entry here
				switch dataType {
				case T_INTEGER:
					// parse int
					parsed_entry, err := parseInt(entry)
					if err != nil {
						return nil, err
					}

					entries = append(entries, parsed_entry)
					break
				case T_SIMPLE_STRING:
					// parse simple string
					parsed_entry, err := parseSimpleString(entry)
					if err != nil {
						return nil, err
					}
					entries = append(entries, parsed_entry)
					break
				default:
					entries = append(entries, data)
				}

				data = rest // set the data variable to the remaining unparsed data
			}
		}
	}

	return entries, nil
}

func checkDataType(raw string) (byte, bool, error) {
	data_type := raw[0]

	switch data_type {
	case T_ARRAY:
		return T_ARRAY, true, nil
	case T_BULK_STRING:
		return T_BULK_STRING, true, nil
	case T_SIMPLE_STRING:
		return T_SIMPLE_STRING, false, nil
	case T_INTEGER:
		return T_INTEGER, false, nil
	default:
		return T_INVALID, false, errors.New("Invalid raw string")
	}
}

func parseInt(raw string) (int, error) {
	prefix := raw[0]
	var data string
	if prefix == '-' {
		data = strings.TrimSuffix(raw, CRLF)
	} else {
		data = strings.TrimSuffix(raw[1:], CRLF)
	}
	return strconv.Atoi(data)
}

func parseBulkString(raw string) (string, error) {
	meta, data, found := strings.Cut(raw, CRLF)
	if !found {
		return "", errors.New("Invalid Bulk String")
	}

	size, err := strconv.Atoi(meta[1:])
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(data[:size]), nil
}

func parseSimpleString(raw string) (string, error) {
	data_type := raw[0]
	if data_type != T_SIMPLE_STRING {
		return "", errors.New("Invalid Simple String")
	}

	return strings.TrimSpace(raw)[1:], nil
}
