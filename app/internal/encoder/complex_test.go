package encoder

import (
	"testing"
)

func TestNewArray(t *testing.T) {
	expected := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	result := NewArray(
		NewBulkString("foo"),
		NewBulkString("bar"),
	)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
		return
	}

	t.Log("TestNewArray passed")
}
