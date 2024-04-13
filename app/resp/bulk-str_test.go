package resp

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/utl"
)

func TestBulkString(t *testing.T) {
	data := map[string]string{
		"PONG":  "$4\r\nPONG\r\n",
		"HELLO": "$5\r\nHELLO\r\n$5\r\nWORLD\r\n",
		"OK":    "$2\r\nOK\r\n",
	}

	for expected, input := range data {
		b := NewBulkString([]byte(input))
		parsed := b.decode()

		if parsed != expected {
			t.Errorf("Failed to parse bulk string: %v", utl.Unquote(input))
		}
		t.Logf("Test passed successfully: %v -> %v", utl.Unquote(input), parsed)
	}
}

func TestBulkString_UsingHandler(t *testing.T) {
	type testresult = [2]interface{}
	data := map[[2]interface{}][]byte{
		{"PONG", 4}:  []byte("$4\r\nPONG\r\n"),
		{"HELLO", 5}: []byte("$5\r\nHELLO\r\n$5\r\nWORLD\r\n"),
		{"OK", 2}:    []byte("$2\r\nOK\r\n"),
	}

	for expected, input := range data {
		expectedStr := expected[0].(string)
		expectedLen := expected[1].(int)

		b, len := HandleResp(input)
		parsed := b.String()

		if len != expectedLen {
			t.Errorf("Failed to parse bulk string, wrong length: %v, %v != %v", utl.Unquote(string(input)), len, expectedLen)
		}

		if parsed != expectedStr {
			t.Errorf("Failed to parse bulk string: %v != %v", utl.Unquote(string(input)), parsed)
			return
		}
		t.Logf("Test passed successfully: %v -> %v", utl.Unquote(string(input)), parsed)
	}
}
