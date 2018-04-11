package search

import "testing"

func TestGetLineNum(t *testing.T) {
	var num int
	var data []byte
	data = []byte("first line\r\nsecond line\nthird line")
	num = getLineNum(data)
	if num != 3 {
		t.Error("Expected 3, got ", num)
	}
}
