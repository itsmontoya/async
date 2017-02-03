package aio

import (
	"fmt"
	//	"os"
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	var (
		f     *File
		wf    *File
		oresp *OpenResp
	)

	aio := New()

	if oresp = <-aio.Open("./test.txt"); oresp.err != nil {
		t.Fatal(oresp.err)
	}

	f = oresp.f

	if oresp = <-aio.OpenFile("./testWrite.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); oresp.err != nil {
		t.Fatal(oresp.err)
	}

	wf = oresp.f

	var (
		buf [32]byte
		n   int
		err error
	)

	if n, err = f.Read(buf[:]); err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(buf[:n]))

	if _, err = wf.Write([]byte("hai hai hai!")); err != nil {
		t.Fatal(err)
	}

	if err = wf.Close(); err != nil {
		t.Fatal(err)
	}

	if err = <-aio.Delete("./testWrite.txt"); err != nil {
		t.Fatal(err)
	}
}
