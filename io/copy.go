package io

import (
	"io"

	"github.com/itsmontoya/async"
)

// Copy will copy
func Copy(dst io.Writer, src io.Reader) (n int64, err error) {
	resp := <-CopyAsync(dst, src)
	n = resp.N
	err = resp.Err
	p.releaseCopyResp(resp)
	return
}

// CopyAsync will copy asynchronously
func CopyAsync(dst io.Writer, src io.Reader) <-chan *CopyResp {
	// Acquire request
	req := p.acquireCopyReq()

	// Set values
	req.w = dst
	req.r = src

	// Send request to queue
	async.Queue(req)
	return req.resp
}
