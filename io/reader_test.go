package io

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	var (
		f   *os.File
		err error
	)

	if f, err = os.Open("../file/testing/helloWorld.txt"); err != nil {
		t.Fatal(err)
	}

	rdr := NewReader(f)
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, rdr)

	if str := buf.String(); str != "Hello world!" {
		t.Fatal(fmt.Errorf("invalid string\nExpected: %s\nReceived: %s\n", "Hello world!", str))
	}
}
