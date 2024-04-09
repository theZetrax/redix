package decoder

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/r3labs/diff/v3"
)

func TestParseSimpleString(t *testing.T) {
	result, err := parseSimpleString("+OK\r\n")
	if err != nil {
		t.Error(err)
	}

	if result != "OK" {
		t.Error("Expected OK, got ", result)
	}

	t.Log("TestParseSimpleString case 1 passed")

	result, err = parseSimpleString("OK\r\n")

	if err == nil {
		t.Error("Failed to detect missing prefix")
	}

	t.Log("TestParseSimpleString case 2 passed")
}

func TestParseBulkString(t *testing.T) {
	input_expected_map := map[string]string{
		"$5\r\nhello\r\n":            "hello",
		"$0\r\n\r\n":                 "",
		"$3\r\nfoo\r\nsomethingelse": "foo",
	}
	input_msg_map := map[string]string{
		"$5hello": "Failed to detect wrong format",
	}

	count := 1
	for input, expected := range input_expected_map {
		result, err := parseBulkString(input)
		if err != nil {
			t.Error(err)
		}

		if result != expected {
			t.Errorf("[ERROR] Expected: `%s`, Input: %s, Result: `%s`", expected, strings.ReplaceAll(input, CRLF, "[CRLF]"), result)
		}

		t.Logf("TestParseBulkString case %d passed", count)
		count += 1
	}

	for input, msg := range input_msg_map {
		_, err := parseBulkString(input)
		if err == nil {
			t.Error(msg)
		}

		t.Logf("TestParseBulkString case %d passed", count)
		count += 1
	}
}

func TestParseInt(t *testing.T) {
	input_expected_map := map[string]int{
		"-1\r\n": -1,
		"+1\r\n": 1,
		"+0\r\n": 0,
		"+2\r\n": 2,
		"-2\r\n": -2,
	}

	count := 1
	for input, expected := range input_expected_map {
		result, err := parseInt(input)
		if err != nil {
			t.Error(err)
		}

		if result != expected {
			t.Logf("Expected %d for input %s, got %d", expected, input, result)
			t.Error(diff.Diff(input, expected))
		}

		t.Log(fmt.Sprintf("TestParseInt case %d passed", count))
		count += 1
	}
}

func TestParseArray(t *testing.T) {
	input_expected_map := map[string][]any{
		"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n":              {"foo", "bar"},
		"*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n": {"foo", "bar", "baz"},
		"*0\r\n":                                    {},
		"*1\r\n:1\r\n":                              {1},
		"*1\r\n+SIMPLE\r\n":                         {"SIMPLE"},
		"*1\r\n$12\r\nhello world!\r\n":             {"hello world!"},
		"*1\r\n$0\r\n\r\n":                          {""},
		"*1\r\n*2\r\n:1\r\n:2\r\n\r\n":              {[]any{1, 2}},
		"*1\r\n*2\r\n+NO\r\n+YES\r\n\r\n":           {[]any{"NO", "YES"}},
		"*1\r\n*2\r\n$2\r\nNO\r\n$3\r\nYES\r\n\r\n": {[]any{"NO", "YES"}},
		"*1\r\n$4\r\nPING\r\n":                      {"PING"},
	}

	for input, expected := range input_expected_map {
		result, err := ParseArray(input)
		if err != nil {
			t.Error(err)
		}

		if len(result) != len(expected) {
			t.Log("Result: ", result)
			t.Errorf("Expected %d elements, got %d", len(expected), len(result))
		}

		for i, v := range result {
			t_expected := reflect.TypeOf(expected[i])

			switch t_expected.Kind() {
			case reflect.Slice:
				if reflect.TypeOf(v).Kind() != reflect.Slice {
					t.Error("Expected slice, got", reflect.TypeOf(v))
				}
				for j := range expected[i].([]any) {
					val_entry := reflect.ValueOf(v).Index(j)
					expected_entry := reflect.ValueOf(expected[i]).Index(j)

					if reflect.DeepEqual(val_entry, expected_entry) {
						t.Error("Expected", expected_entry, ", got", val_entry)
						t.Error("Expected", reflect.TypeOf(expected_entry).Kind(), ", got", reflect.TypeOf(val_entry).Kind())
					}
				}
				break
			default:

				if v != expected[i] {
					t.Error("Expected", expected[i], fmt.Sprintf("(%s)", reflect.TypeOf(expected[i])), ", got", v, fmt.Sprintf("(%s)", reflect.TypeOf(v)))
				}
			}
		}

		t.Log("TestParseArray passed")
	}

}
