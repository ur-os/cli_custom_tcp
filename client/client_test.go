package main

import (
	mock_server "github"
	"testing"
)

func TestUnit(t *testing.T) {
	buffer := []byte("qwerty")
	left := 0
	bytesNumb := int32(len(buffer))
	expected1 := []byte("qwerty")
	expected2 := len(buffer)

	result1, result2 := readSliceBytePacket(buffer, left, bytesNumb)

	if result1[3] != expected1[3] {
		t.Errorf("Expect %s, got %s", expected1, result1)
	}

	if result2 != expected2 {
		t.Errorf("Expect %d, got %d", expected2, result2)
	}

}
