package io

import (
	"io"

	"github.com/itsmontoya/async"
)

// Copy will copy
// Note: This should be used only with disk to disk copying (or at least syscall reader to syscall writer)
// If only one of the args is syscall-blocking, use the Reader or Writer wrapper for that particular item
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
