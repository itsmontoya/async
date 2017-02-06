package aio

import (
	"os"

	"github.com/missionMeteora/toolkit/errors"
)

func newFile(r *openRequest, rq chan<- interface{}) (f *File, err error) {
	f = p.acquireFile()
	if f.f, err = os.OpenFile(r.key, r.flag, r.perm); err != nil {
		f = nil
		return
	}

	f.rq = rq
	return
}

// File is a file
type File struct {
	f *os.File
	// Request queue
	rq chan<- interface{}
	// Closed state
	closed bool
}

// Read will read a file
func (f *File) Read(b []byte) (n int, err error) {
	rr := <-f.ReadAsync(b)
	n = rr.N
	err = rr.Err
	p.releaseRWResp(rr)
	return
}

// ReadAsync will read a file asynchronously
func (f *File) ReadAsync(b []byte) <-chan *RWResp {
	r := p.acquireReadReq()
	r.b = b
	r.f = f.f
	f.rq <- r
	return r.resp
}

// Write will write to a file
func (f *File) Write(b []byte) (n int, err error) {
	rr := <-f.WriteAsync(b)
	return rr.N, rr.Err
}

// WriteAsync will write to a file asynchronously
func (f *File) WriteAsync(b []byte) <-chan *RWResp {
	var r writeRequest
	r.b = b
	r.resp = make(chan *RWResp)
	r.f = f.f
	f.rq <- &r
	return r.resp
}

// Close will close a file
func (f *File) Close() error {
	return <-f.CloseAsync()
}

// CloseAsync will close a file asynchronously
func (f *File) CloseAsync() <-chan error {
	r := p.acquireCloseReq()
	if f.closed {
		r.resp <- errors.ErrIsClosed
	} else {
		r.f = f
		f.rq <- r
	}

	return r.resp
}
