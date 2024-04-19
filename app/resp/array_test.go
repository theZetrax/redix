package resp

import (
	"reflect"
	"strings"
	"testing"
)

func TestArray(t *testing.T) {
	data := map[string][]string{
		"*1\r\n$2\r\nOK\r\n":                              {"OK"},
		"*1\r\n$4\r\nPONG\r\n":                            {"PONG"},
		"*2\r\n$3\r\nGET\r\n$3\r\nKEY\r\n":                {"GET", "KEY"},
		"*2\r\n$5\r\nHELLO\r\n$5\r\nWORLD\r\n":            {"HELLO", "WORLD"},
		"*3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n": {"SET", "KEY", "VALUE"},
	}

	for input, expected := range data {
		a := NewArray([]byte(input))
		parsed := a.Parsed

		same := true
		for i, b := range parsed {
			if b != expected[i] {
				same = false
			}
		}

		if !same {
			t.Errorf("Failed to parse array: %v", input)
			return
		}
		t.Logf("Test passed successfully: %v", strings.ReplaceAll(input, "\r\n", "\\r\\n"))
	}
}

func TestArray_Multiple(t *testing.T) {
	data := map[string][][]string{
		"*1\r\n$2\r\nOK\r\n*1\r\n$4\r\nPONG\r\n":                                                     {{"OK"}, {"PONG"}},
		"*2\r\n$3\r\nGET\r\n$3\r\nKEY\r\n*2\r\n$5\r\nHELLO\r\n$5\r\nWORLD\r\n":                       {{"GET", "KEY"}, {"HELLO", "WORLD"}},
		"*3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n*1\r\n$2\r\nOK\r\n":                          {{"SET", "KEY", "VALUE"}, {"OK"}},
		"*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\n123\r\n*3\r\n$3\r\nSET\r\n$3\r\nbar\r\n$3\r\n456\r\n": {{"SET", "foo", "123"}, {"SET", "bar", "456"}},
	}

	for input_str, expected := range data {
		input := []byte(input_str)

		arr := NewArray(input)
		parsed := arr.Parsed

		for i, b := range parsed {
			if reflect.TypeOf(b).Kind() != reflect.Slice {
				t.Errorf("Failed to parse array: %v", strings.ReplaceAll(string(input), "\r\n", "\\r\\n"))
			}

			same := true
			for j, c := range b.([]any) {
				if c != expected[i][j] {
					same = false
				}
			}

			if !same {
				t.Errorf("Failed to parse array: %v", input)
				return
			}
		}

		t.Logf("Test passed successfully: %v", strings.ReplaceAll(input_str, "\r\n", "\\r\\n"))
	}
}

func Test_IsValid(t *testing.T) {
	// *3 $3 SET $3 bar $3 456 \n
	input := "*3\r\n$3\r\nSET\r\n$3\r\nbar\r\n$3\r\n456\r\n\n"
	expected := true

	if IsValidArray([]byte(input)) != expected {
		t.Errorf("Failed to validate array: %v", input)
		return
	}
	t.Logf("Test passed successfully: %v", strings.ReplaceAll(input, "\r\n", "\\r\\n"))
}
