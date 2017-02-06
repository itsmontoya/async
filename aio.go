package aio

import (
	"os"
	//	"github.com/missionMeteora/toolkit/errors"
)

var p = newPools()

// New returns a new AIO
func New(numThreads int) *AIO {
	a := AIO{
		rq: make(chan interface{}, 1024*32),
	}

	for i := 0; i < numThreads; i++ {
		t := newThread(a.rq)
		go t.listen()
	}

	return &a
}

// AIO does stuff
type AIO struct {
	rq chan interface{}
}

// Open will open a new file for reading
func (a *AIO) Open(key string) (f *File, err error) {
	return a.OpenFile(key, os.O_RDONLY, 0)
}

// OpenFile will open a new file with flag and perm
func (a *AIO) OpenFile(key string, flag int, perm os.FileMode) (f *File, err error) {
	resp := <-a.OpenFileAsync(key, flag, perm)
	f = resp.F
	err = resp.Err
	p.releaseOpenResp(resp)
	return
}

// OpenFileAsync will open a new file with flag and perm asynchronously
func (a *AIO) OpenFileAsync(key string, flag int, perm os.FileMode) <-chan *OpenResp {
	or := p.acquireOpenReq()
	or.key = key
	or.flag = flag
	or.perm = perm
	a.rq <- or
	return or.resp
}

// Delete will delete a file
func (a *AIO) Delete(key string) <-chan error {
	dr := p.acquireDelReq()
	dr.key = key
	a.rq <- dr
	return dr.resp
}

// OpenResp is a response for open requests
type OpenResp struct {
	F   *File
	Err error
}

// RWResp is a response for read/write requests
type RWResp struct {
	N   int
	Err error
}

type openRequest struct {
	key  string
	flag int
	perm os.FileMode

	resp chan *OpenResp
}

type readRequest struct {
	f *os.File
	b []byte

	resp chan *RWResp
}

type writeRequest struct {
	f *os.File
	b []byte

	resp chan *RWResp
}

type closeRequest struct {
	f *File

	resp chan error
}

type deleteRequest struct {
	key string

	resp chan error
}

type sema chan struct{}

func newFile(r *openRequest, rq chan<- interface{}) (f *File, err error) {
	f = p.acquireFile()
	if f.f, err = os.OpenFile(r.key, r.flag, r.perm); err != nil {
		f = nil
		return
	}

	f.rq = rq
	return
}
