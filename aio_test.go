package aio

import (
	"fmt"
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	var (
		f   *File
		wf  *File
		bs  []byte
		err error
	)

	aio := New()
	fc, ec := aio.Open("./test.txt")

	select {
	case f = <-fc:
	case err = <-ec:
		t.Fatal(err)
		return
	}

	bc, ec := f.Read()

	select {
	case bs = <-bc:
	case err = <-ec:
		t.Fatal(err)
		return
	}

	fmt.Println(string(bs))

	fc, ec = aio.OpenFile("./testWrite.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	select {
	case wf = <-fc:
	case err = <-ec:
		t.Fatal(err)
		return
	}

	dc, ec := wf.Write([]byte("hai hai hai!"))
	select {
	case <-dc:
	case err = <-ec:
		t.Fatal(err)
		return
	}

	sc, ec := wf.Close()
	select {
	case <-sc:
	case err = <-ec:
		t.Fatal(err)
		return
	}

	sc, ec = aio.Delete("./testWrite.txt")
	select {
	case <-sc:
	case err = <-ec:
		t.Fatal(err)
		return
	}
}
