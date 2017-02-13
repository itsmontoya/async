package io

import "io"

// NewReader will return a new asynchronous reader
func NewReader(r io.Reader) *Reader {
	return &Reader{r}
}

// Reader is an asynchronous wrapper for an io.Reader
type Reader struct {
	// Internal reader
	r io.Reader
}

// Read will read
func (r *Reader) Read(b []byte) (n int, err error) {
	return read(r.r, b)
}

// ReadAsync will read asynchronously
func (r *Reader) ReadAsync(b []byte) <-chan *RWResp {
	return readAsync(r.r, b)
}

// NewReadCloser will return a new asynchronous read closer
func NewReadCloser(rc io.ReadCloser) *ReadCloser {
	return &ReadCloser{rc}
}

// ReadCloser is an asynchronous wrapper for an io.ReadCloser
type ReadCloser struct {
	// Internal reader
	rc io.ReadCloser
}

// Read will read
func (rc *ReadCloser) Read(b []byte) (n int, err error) {
	return read(rc.rc, b)
}

// ReadAsync will read asynchronously
func (rc *ReadCloser) ReadAsync(b []byte) <-chan *RWResp {
	return readAsync(rc.rc, b)
}

// Close will close
func (rc *ReadCloser) Close() (err error) {
	return close(rc.rc)
}

// CloseAsync will close asynchronously
func (rc *ReadCloser) CloseAsync() <-chan error {
	return closeAsync(rc.rc)
}
