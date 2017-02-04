package main

import (
	"github.com/itsmontoya/aio"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
)

func main() {
	s := &srv{
		aio: aio.New(),
	}

	r := httprouter.New()
	r.GET("/a", s.handleA)
	r.GET("/b", s.handleB)
	http.ListenAndServe(":1337", r)
}

type srv struct {
	aio *aio.AIO
}

func (s *srv) handleA(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	or := <-s.aio.Open("../testing/declarationOfIndependence.txt")
	if or.Err != nil {
		panic(or.Err)
	}

	if _, err := io.Copy(w, or.F); err != nil {
		panic(err)
	}

	if err := or.F.Close(); err != nil {
		panic(err)
	}
}

func (s *srv) handleB(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	f, err := os.Open("../testing/declarationOfIndependence.txt")
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(w, f); err != nil {
		panic(err)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}
}
