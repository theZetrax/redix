package resp

import (
	"reflect"
	"strings"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/utl"
)

func escape_all(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\\r\\n")
}

func Test_Parser(t *testing.T) {
	data := map[string]string{
		"$5\r\nhello\r\n$4\r\nping\r\n": "hello",
		"+simple\r\n":                   "simple",
	}

	for k, v := range data {
		result, _, _ := parse([]byte(k))
		if result.(string) != v {
			t.Errorf("Expected %v, got %v", v, result)
		}
	}
}

func Test_Large_String(t *testing.T) {
	data := map[string]string{
		"$100\r\n" + strings.Repeat("a", 100) + "\r\n":   strings.Repeat("a", 100),
		"$1000\r\n" + strings.Repeat("a", 1000) + "\r\n": strings.Repeat("a", 1000),
	}

	for input, expected := range data {
		result, _, _ := parse([]byte(input))

		if result.(string) != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}
}

func Test_Large_RDB(t *testing.T) {
	type test_data struct {
		expected string
		rest     string
	}
	data := map[string]test_data{
		"$100\r\n" + strings.Repeat("a", 100) + "*1\r\n$3\r\nfoo\r\n": {
			expected: strings.Repeat("a", 100),
			rest:     "*1\r\n$3\r\nfoo\r\n",
		},
	}

	for input, expected := range data {
		result, rest, _ := parse([]byte(input))

		if result.(string) != expected.expected {
			t.Errorf("Expected %v, got %v", expected.expected, result)
		}
		if string(rest) != expected.rest {
			t.Errorf("Expected Rest %q, got %q", string(expected.rest), string(rest))
		}
	}
}

func Test_Array_Single(t *testing.T) {
	type test_data struct {
		Expected []string
		Left     string
	}
	data := map[string]test_data{
		"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*1\r\n$3\r\nbaz\r\n": {[]string{"foo", "bar"}, "*1\r\n$3\r\nbaz\r\n"},
		"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n":                    {[]string{"foo", "bar"}, ""},
		"*1\r\n$3\r\nfoo\r\n":                                 {[]string{"foo"}, ""},
		"*0\r\n":                                              {[]string{}, ""},
		"*-1\r\n":                                             {[]string{}, ""},
		"*11\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n": {[]string{"foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo", "foo"}, ""},
	}

	for input, exp_results := range data {
		result, left, _ := parse([]byte(input))
		result_str := utl.ToStringSlice(result.([]interface{}))

		for i, v := range result_str {
			if v != exp_results.Expected[i] {
				t.Errorf("Expected %v, got %v", exp_results.Expected[i], v)
			}
		}

		if string(left) != exp_results.Left {
			t.Errorf("Expected Left %v, got %v", escape_all(exp_results.Left), escape_all(string(left)))
		}
	}
}

func Test_Array_Nested(t *testing.T) {
	type test_data struct {
		Expected []string
		Left     string
	}
	data := map[string]test_data{
		"*1\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*1\r\nbar\r\n": {
			Expected: []string{"foo", "bar"},
			Left:     "*1\r\nbar\r\n",
		},
	}

	for input, exp_results := range data {
		result, left, _ := parse([]byte(input))
		result_str := utl.ToStringSlice(result.([]interface{})[0].([]interface{}))

		for i, v := range result_str {
			if v != exp_results.Expected[i] {
				t.Errorf("Expected %v, got %v", exp_results.Expected[i], v)
			}
		}

		if string(left) != exp_results.Left {
			t.Errorf("Expected Left %v, got %v", escape_all(exp_results.Left), escape_all(string(left)))
		}
	}
}

func Test_Array_RDB(t *testing.T) {
	rdb_hex := []byte("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
	rdb_enc, _ := utl.DecodeHexToBinary(rdb_hex)

	input := "+FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0\r\n" + string(EncodeFileContent(rdb_enc)) + "*3\r\n$8\r\nREPLCONF\r\n$6\r\nGETACK\r\n$1\r\n*\r\n"
	expected := []string{
		"FULLRESYNC 75cd7bc10c49047e0d163660f3b90625b1af31dc 0",
		string(rdb_enc),
	}
	expected_arr := []string{"REPLCONF", "GETACK", "*"}

	value1, rest1, _ := parse([]byte(input))
	if value1.(string) != expected[0] {
		t.Errorf("Expected %q, got %q", expected[0], value1)
	}

	value2, rest2, _ := parse(rest1)
	if value2.(string) != expected[1] {
		t.Errorf("Expected %v, got %v", expected[1], value2.(string))
	}

	value3, rest3, _ := parse(rest2)
	if reflect.DeepEqual(value3, expected_arr) {
		t.Errorf("Expected %v, got %v", expected_arr, value3)
	}

	if string(rest3) != "" {
		t.Errorf("Expected empty rest, got %q", rest3)
	}
}
