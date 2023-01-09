package buffer

import (
	"bytes"
	"testing"
)

func TestBufWriteAtStartingPosition(t *testing.T) {
	testdata := "This is a single line string"
	b := NewBuffer().(*buf)
	b.WriteAt(0, []byte(testdata))

	r := string(b.data[:bytes.IndexByte(b.data[:], 0)])
	if r != testdata {
		t.Errorf("(%s) is different from expected: %s", r, testdata)
	}
}

func TestBufWriteAtMiddle(t *testing.T) {
	testPos, testRunes := 22, []byte{'s', 's'}
	expected := "This is a single line ssstring"
	b := NewBuffer().(*buf)
	b.WriteAt(0, []byte("This is a single line string"))

	b.WriteAt(testPos, testRunes)

	r := string(b.data[:bytes.IndexByte(b.data[:], 0)])
	if r != expected {
		t.Errorf("(%s) is different from expected: %s", r, expected)
	}
}
