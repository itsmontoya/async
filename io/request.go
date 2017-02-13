package io

import "io"

func newReadReq() *readRequest {
	return &readRequest{
		resp: make(chan *RWResp),
	}
}

type readRequest struct {
	r io.Reader
	b []byte

	resp chan *RWResp
}

func (req *readRequest) Action() {
	resp := newRWResp()
	resp.N, resp.Err = req.r.Read(req.b)
	req.resp <- resp
	p.releaseReadReq(req)
}

func newWriteReq() *writeRequest {
	return &writeRequest{
		resp: make(chan *RWResp),
	}
}

type writeRequest struct {
	w io.Writer
	b []byte

	resp chan *RWResp
}

func (req *writeRequest) Action() {
	resp := newRWResp()
	resp.N, resp.Err = req.w.Write(req.b)
	req.resp <- resp
	p.releaseWriteReq(req)
}

func newCloseReq() *closeRequest {
	return &closeRequest{
		resp: make(chan error),
	}
}

type closeRequest struct {
	c io.Closer
	b []byte

	resp chan error
}

func (req *closeRequest) Action() {
	req.resp <- req.c.Close()
	p.releaseCloseReq(req)
}

func newCopyReq() *copyRequest {
	return &copyRequest{
		resp: make(chan *CopyResp),
	}
}

type copyRequest struct {
	w io.Writer
	r io.Reader

	resp chan *CopyResp
}

func (req *copyRequest) Action() {
	resp := newCopyResp()
	resp.N, resp.Err = io.Copy(req.w, req.r)
	req.resp <- resp
	p.releaseCopyReq(req)
}
