package aio

import (
	"os"

	"github.com/missionMeteora/toolkit/errors"
)

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
	return rr.n, rr.err
}

// ReadAsync will read a file asynchronously
func (f *File) ReadAsync(b []byte) <-chan *RWResp {
	var r readRequest
	r.b = b
	r.resp = make(chan *RWResp, 1)
	r.f = f.f
	f.rq <- &r
	return r.resp
}

// Write will write to a file
func (f *File) Write(b []byte) (n int, err error) {
	rr := <-f.WriteAsync(b)
	return rr.n, rr.err
}

// WriteAsync will write to a file asynchronously
func (f *File) WriteAsync(b []byte) <-chan *RWResp {
	var r writeRequest
	r.b = b
	r.resp = make(chan *RWResp, 1)
	r.f = f.f
	f.rq <- &r
	return r.resp
}

// Delete will delete a file
func (f *File) Delete(key string) error {
	return <-f.DeleteAsync(key)
}

// DeleteAsync will delete a file asynchronously
func (f *File) DeleteAsync(key string) <-chan error {
	var r deleteRequest
	r.key = key
	r.errCh = make(chan error, 1)
	f.rq <- &r
	return r.errCh
}

// Close will close a file
func (f *File) Close() error {
	return <-f.CloseAsync()
}

// CloseAsync will close a file asynchronously
func (f *File) CloseAsync() <-chan error {
	var r closeRequest
	r.errCh = make(chan error, 1)

	if f.closed {
		r.errCh <- errors.ErrIsClosed
	} else {
		r.f = f.f
		f.rq <- &r
		f.closed = true
		f.rq = nil
		f.f = nil
	}

	return r.errCh
}
