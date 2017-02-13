package io

import (
	"io"

	"github.com/itsmontoya/aio"
)

func read(r io.Reader, b []byte) (n int, err error) {
	resp := <-readAsync(r, b)
	n = resp.N
	err = resp.Err
	p.releaseRWResp(resp)
	return
}

func readAsync(r io.Reader, b []byte) <-chan *RWResp {
	// Acquire request
	req := p.acquireReadReq()

	// Set values
	req.r = r
	req.b = b

	// Send request to queue
	aio.Queue(req)
	return req.resp
}

func write(w io.Writer, b []byte) (n int, err error) {
	resp := <-writeAsync(w, b)
	n = resp.N
	err = resp.Err
	p.releaseRWResp(resp)
	return
}

func writeAsync(w io.Writer, b []byte) <-chan *RWResp {
	// Acquire request
	req := p.acquireWriteReq()

	// Set values
	req.w = w
	req.b = b

	// Send request to queue
	aio.Queue(req)
	return req.resp
}

func close(c io.Closer) (err error) {
	return <-closeAsync(c)
}

func closeAsync(c io.Closer) <-chan error {
	// Acquire request
	req := p.acquireCloseReq()

	// Set values
	req.c = c

	// Send request to queue
	aio.Queue(req)
	return req.resp
}
