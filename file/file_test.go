package file

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

var testBuf = bytes.NewBuffer(nil)

const (
	testFilePath = "./testing/declarationOfIndependence.txt"
)

func TestBasic(t *testing.T) {
	var (
		f   *File
		wf  *File
		buf [32]byte
		n   int
		err error
	)

	if f, err = Open("./testing/helloWorld.txt"); err != nil {
		t.Fatal(err)
	}

	if wf, err = OpenFile("./testing/testWrite.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		t.Fatal(err)
	}

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

	if err = <-Delete("./testing/testWrite.txt"); err != nil {
		t.Fatal(err)
	}

	if err = f.Close(); err != nil {
		t.Fatal(err)
	}

	if f, err = Open(testFilePath); err != nil {
		t.Fatal(err)
	}

	if _, err = io.Copy(testBuf, f); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkAIO(b *testing.B) {
	var (
		f   *File
		err error
	)

	for i := 0; i < b.N; i++ {
		testBuf.Reset()
		if f, err = Open(testFilePath); err != nil {
			b.Fatal(err)
		}

		if _, err = io.Copy(testBuf, f); err != nil {
			b.Fatal(err)
		}

		if err = f.Close(); err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkStdlib(b *testing.B) {
	var (
		f   *os.File
		err error
	)

	for i := 0; i < b.N; i++ {
		testBuf.Reset()
		if f, err = os.Open(testFilePath); err != nil {
			b.Fatal(err)
		}

		if _, err = io.Copy(testBuf, f); err != nil {
			b.Fatal(err)
		}

		if err = f.Close(); err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkAIOPara(b *testing.B) {
	b.SetParallelism(64)
	b.RunParallel(func(pb *testing.PB) {
		buf := bytes.NewBuffer(nil)
		for pb.Next() {
			f, err := Open(testFilePath)
			if err != nil {
				b.Fatal(err)
			}

			if _, err = io.Copy(buf, f); err != nil {
				b.Fatal(err)
			}

			if err = f.Close(); err != nil {
				b.Fatal(err)
			}

			buf.Reset()
		}
	})

	b.ReportAllocs()
}

func BenchmarkStdlibPara(b *testing.B) {
	b.SetParallelism(64)
	b.RunParallel(func(pb *testing.PB) {
		buf := bytes.NewBuffer(nil)
		for pb.Next() {
			f, err := os.Open(testFilePath)
			if err != nil {
				b.Fatal(err)
			}

			if _, err = io.Copy(buf, f); err != nil {
				b.Fatal(err)
			}

			if err = f.Close(); err != nil {
				b.Fatal(err)
			}

			buf.Reset()
		}
	})

	b.ReportAllocs()
}
