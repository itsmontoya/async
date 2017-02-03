package aio

import (
	"os"

	"github.com/missionMeteora/toolkit/errors"
)

// New returns a new AIO
func New() *AIO {
	a := AIO{
		rq: make(chan interface{}, 32),
	}

	t := newThread(a.rq)
	go t.listen()
	return &a
}

// AIO does stuff
type AIO struct {
	rq chan interface{}
}

// Open will open a new file for reading
func (a *AIO) Open(key string) (fc <-chan *File, ec <-chan error) {
	return a.OpenFile(key, os.O_RDONLY, 0)
}

// OpenFile will open a new file with flag and perm
func (a *AIO) OpenFile(key string, flag int, perm os.FileMode) (fc <-chan *File, ec <-chan error) {
	var or openRequest
	or.key = key
	or.flag = flag
	or.perm = perm

	or.fileCh = make(chan *File, 1)
	or.errCh = make(chan error, 1)

	fc = or.fileCh
	ec = or.errCh
	a.rq <- &or
	return
}

// Delete will delete a file
func (a *AIO) Delete(key string) (dc <-chan struct{}, ec <-chan error) {
	var dr deleteRequest
	dr.key = key
	dr.doneCh = make(chan struct{}, 1)
	dr.errCh = make(chan error, 1)

	dc = dr.doneCh
	ec = dr.errCh

	a.rq <- &dr
	return
}

func newThread(rq chan interface{}) *thread {
	return &thread{rq}
}

type thread struct {
	rq chan interface{}
}

func (t *thread) open(r *openRequest) {
	if f, err := newFile(r, t.rq); err == nil {
		r.fileCh <- f
	} else {
		r.errCh <- err
	}
}

func (t *thread) read(r *readRequest) {
	var buf [32]byte
	if n, err := r.f.Read(buf[:]); err == nil {
		r.readCh <- buf[:n]
	} else {
		r.errCh <- err
	}
}

func (t *thread) write(r *writeRequest) {
	if n, err := r.f.Write(r.b); err == nil {
		r.doneCh <- n
	} else {
		r.errCh <- err
	}
}

func (t *thread) close(r *closeRequest) {
	if err := r.f.Close(); err == nil {
		r.doneCh <- struct{}{}
	} else {
		r.errCh <- err
	}
}

func (t *thread) delete(r *deleteRequest) {
	if err := os.Remove(r.key); err == nil {
		r.doneCh <- struct{}{}
	} else {
		r.errCh <- err
	}
}

func (t *thread) listen() {
	for req := range t.rq {
		switch r := req.(type) {
		case *openRequest:
			t.open(r)
		case *readRequest:
			t.read(r)
		case *writeRequest:
			t.write(r)
		case *deleteRequest:
			t.delete(r)
		case *closeRequest:
			t.close(r)

		default:
			panic("invalid type")
		}
	}
}

type openRequest struct {
	key  string
	flag int
	perm os.FileMode

	fileCh chan *File
	errCh  chan error
}

type readRequest struct {
	f *os.File

	readCh chan []byte
	errCh  chan error
}

type writeRequest struct {
	f *os.File
	b []byte

	doneCh chan int
	errCh  chan error
}

type closeRequest struct {
	f *os.File

	doneCh chan struct{}
	errCh  chan error
}

type deleteRequest struct {
	key string

	doneCh chan struct{}
	errCh  chan error
}

type sema chan struct{}

func newFile(r *openRequest, rq chan<- interface{}) (fp *File, err error) {
	var f File
	if f.f, err = os.OpenFile(r.key, r.flag, r.perm); err != nil {
		return
	}

	f.rq = rq
	fp = &f
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
func (f *File) Read() (b []byte, err error) {
	bc, ec := f.ReadAsync()

	select {
	case b = <-bc:
	case err = <-ec:
	}

	return
}

// ReadAsync will read a file asynchronously
func (f *File) ReadAsync() (bc <-chan []byte, ec <-chan error) {
	var r readRequest
	r.readCh = make(chan []byte, 1)
	r.errCh = make(chan error, 1)
	r.f = f.f

	bc = r.readCh
	ec = r.errCh

	f.rq <- &r
	return
}

// Write will write to a file
func (f *File) Write(b []byte) (dc <-chan int, ec <-chan error) {
	var r writeRequest
	r.b = b
	r.doneCh = make(chan int, 1)
	r.errCh = make(chan error, 1)
	r.f = f.f

	dc = r.doneCh
	ec = r.errCh

	f.rq <- &r
	return
}

// Delete will delete a file
func (f *File) Delete(key string) (dc <-chan struct{}, ec <-chan error) {
	var r deleteRequest
	r.key = key
	r.doneCh = make(chan struct{}, 1)
	r.errCh = make(chan error, 1)

	dc = r.doneCh
	ec = r.errCh

	f.rq <- &r
	return
}

// Close will close a file
func (f *File) Close() (dc <-chan struct{}, ec <-chan error) {
	var r closeRequest
	r.f = f.f
	r.doneCh = make(chan struct{}, 1)
	r.errCh = make(chan error, 1)

	dc = r.doneCh
	ec = r.errCh

	if f.closed {
		r.errCh <- errors.ErrIsClosed
		return
	}

	f.rq <- &r
	f.closed = true
	f.rq = nil
	f.f = nil
	return
}
