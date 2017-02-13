package reader

import (
	"io"

	"github.com/itsmontoya/aio"
)

// New will return a new asynchronous reader
func New(r io.Reader) *Reader {
	return &Reader{r}
}

// Reader is an asynchronous wrapper for an io.Reader
type Reader struct {
	// Internal reader
	r io.Reader
}

// Read will read
func (r *Reader) Read(b []byte) (n int, err error) {
	resp := <-r.ReadAsync(b)
	n = resp.N
	err = resp.Err
	p.releaseResponse(resp)
	return
}

// ReadAsync will read asynchronously
func (r *Reader) ReadAsync(b []byte) <-chan *RWResp {
	// Acquire request
	req := p.acquireRequest()

	// Set values
	req.r = r.r
	req.b = b

	// Send request to queue
	aio.Queue(req)
	return req.resp
}
