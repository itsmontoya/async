package file

import (
	"os"

	"github.com/itsmontoya/aio"
	"github.com/missionMeteora/toolkit/errors"
)

func newFile(r *openRequest, a *aio.AIO) (f *File, err error) {
	// Acquire file struct from pool
	f = p.acquireFile()
	// Open underlying os.File
	if f.f, err = os.OpenFile(r.key, r.flag, r.perm); err != nil {
		f = nil
		return
	}

	// Set file's internal aio
	f.a = a
	return
}

// File is a file
type File struct {
	f *os.File
	// Reference AIO instance
	a *aio.AIO
	// Closed state
	closed bool
}

// Read will read a file
func (f *File) Read(b []byte) (n int, err error) {
	// Read and wait for response
	resp := <-f.ReadAsync(b)

	n = resp.N
	err = resp.Err

	// Release response back to the pool
	p.releaseRWResp(resp)
	return
}

// ReadAsync will read a file asynchronously
func (f *File) ReadAsync(b []byte) <-chan *RWResp {
	// Acquire read request from pool
	req := p.acquireReadReq()

	req.b = b
	req.f = f.f

	// Send request to request queue
	f.a.Queue(req)
	return req.resp
}

// Write will write to a file
func (f *File) Write(b []byte) (n int, err error) {
	// Write and wait for response
	resp := <-f.WriteAsync(b)

	n = resp.N
	err = resp.Err

	// Release response back to the pool
	p.releaseRWResp(resp)
	return
}

// WriteAsync will write to a file asynchronously
func (f *File) WriteAsync(b []byte) <-chan *RWResp {
	// Acquire write request from pool
	req := p.acquireWriteReq()

	req.b = b
	req.f = f.f

	// Send request to request queue
	f.a.Queue(req)
	return req.resp
}

// Seek will seek within a file
func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	// Seek and wait for response
	resp := <-f.SeekAsync(offset, whence)

	ret = resp.Ret
	err = resp.Err

	// Release response back to the pool
	p.releaseSeekResp(resp)
	return
}

// SeekAsync will seek within a file asynchronously
func (f *File) SeekAsync(offset int64, whence int) <-chan *SeekResp {
	// Acquire seek request from pool
	req := p.acquireSeekReq()

	req.f = f.f
	req.offset = offset
	req.whence = whence

	// Send request to request queue
	f.a.Queue(req)
	return req.resp
}

// Sync will sync a file
func (f *File) Sync() (err error) {
	// Sync and wait for response
	return <-f.SyncAsync()
}

// SyncAsync will sync a file asynchronously
func (f *File) SyncAsync() <-chan error {
	// Acquire seek request from pool
	req := p.acquireSyncReq()

	req.f = f.f

	// Send request to request queue
	f.a.Queue(req)
	return req.resp
}

// Close will close a file
func (f *File) Close() error {
	return <-f.CloseAsync()
}

// CloseAsync will close a file asynchronously
func (f *File) CloseAsync() <-chan error {
	req := p.acquireCloseReq()
	if f.closed {
		// File is already closed, send error to response
		req.resp <- errors.ErrIsClosed
	} else {
		f.closed = true
		req.f = f
		f.a.Queue(req)
	}

	return req.resp
}
