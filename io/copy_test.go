package io

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/itsmontoya/async"
)

const testDOILoc = "../testing/declarationOfIndependence.txt"

func TestCopy(t *testing.T) {
	async.Set(2)

	buf := bytes.NewBuffer(nil)
	f, err := os.Open(testDOILoc)
	if err != nil {
		t.Fatal(err)
	}

	n, err := Copy(buf, f)
	if err != nil {
		t.Fatal(err)
	}

	if n != 9284 {
		t.Fatalf("unexpected amount copied\nExpected: %v\nReturned: %v\n", 9284, n)
	}
}

func BenchmarkCopy(b *testing.B) {
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		w, err := ioutil.TempFile("", "async_copy_test_")
		if err != nil {
			b.Fatal(err)
		}
		fl := w.Name()

		r, err := os.Open(testDOILoc)
		if err != nil {
			b.Fatal(err)
		}

		for pb.Next() {
			if _, err = Copy(w, r); err != nil {
				b.Fatal(err)
			}
			if _, err = r.Seek(0, 0); err != nil {
				b.Fatal(err)
			}
		}

		r.Close()
		w.Close()
		os.Remove(fl)
	})

	b.ReportAllocs()
}

func BenchmarkStdlibCopy(b *testing.B) {
	b.SetParallelism(8)
	b.RunParallel(func(pb *testing.PB) {
		w, err := ioutil.TempFile("", "async_copy_test_")
		if err != nil {
			b.Fatal(err)
		}
		fl := w.Name()

		r, err := os.Open(testDOILoc)
		if err != nil {
			b.Fatal(err)
		}

		for pb.Next() {
			if _, err = io.Copy(w, r); err != nil {
				b.Fatal(err)
			}
			if _, err = r.Seek(0, 0); err != nil {
				b.Fatal(err)
			}
		}

		r.Close()
		w.Close()
		os.Remove(fl)
	})

	b.ReportAllocs()
}
