package aio

import (
	"os"

	"github.com/missionMeteora/toolkit/errors"
)

func newFile(r *openRequest, rq chan<- interface{}) (f *File, err error) {
	// Acquire file struct from pool
	f = p.acquireFile()
	// Open underlying os.File
	if f.f, err = os.OpenFile(r.key, r.flag, r.perm); err != nil {
		f = nil
		return
	}
	// Set file's request queue (send-only)
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
	// Read and wait for response
	rr := <-f.ReadAsync(b)

	n = rr.N
	err = rr.Err

	// Release response back to the pool
	p.releaseRWResp(rr)
	return
}

// ReadAsync will read a file asynchronously
func (f *File) ReadAsync(b []byte) <-chan *RWResp {
	// Acquire read request from pool
	r := p.acquireReadReq()

	r.b = b
	r.f = f.f

	// Send request to request queue
	f.rq <- r
	return r.resp
}

// Write will write to a file
func (f *File) Write(b []byte) (n int, err error) {
	// Write and wait for response
	rr := <-f.WriteAsync(b)

	n = rr.N
	err = rr.Err

	// Release response back to the pool
	p.releaseRWResp(rr)
	return
}

// WriteAsync will write to a file asynchronously
func (f *File) WriteAsync(b []byte) <-chan *RWResp {
	// Acquire write request from pool
	r := p.acquireWriteReq()

	r.b = b
	r.f = f.f

	// Send request to request queue
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
		// File is already closed, send error to response
		r.resp <- errors.ErrIsClosed
	} else {
		f.closed = true
		r.f = f
		f.rq <- r
	}

	return r.resp
}
