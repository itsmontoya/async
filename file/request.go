package file

import (
	"os"

	"github.com/itsmontoya/aio"
)

func newOpenReq() *openRequest {
	return &openRequest{
		resp: make(chan *OpenResp),
	}
}

type openRequest struct {
	key  string
	flag int
	perm os.FileMode

	a    *aio.AIO
	resp chan *OpenResp
}

func (req *openRequest) Action() {
	resp := p.acquireOpenResp()
	resp.F, resp.Err = newFile(req, req.a)
	req.resp <- resp
	p.releaseOpenReq(req)
}

func newReadReq() *readRequest {
	return &readRequest{
		resp: make(chan *RWResp),
	}
}

type readRequest struct {
	f *os.File
	b []byte

	resp chan *RWResp
}

func (req *readRequest) Action() {
	resp := p.acquireRWResp()
	resp.N, resp.Err = req.f.Read(req.b)
	req.resp <- resp
	p.releaseReadReq(req)
}

func newWriteReq() *writeRequest {
	return &writeRequest{
		resp: make(chan *RWResp),
	}
}

type writeRequest struct {
	f *os.File
	b []byte

	resp chan *RWResp
}

func (req *writeRequest) Action() {
	resp := p.acquireRWResp()
	resp.N, resp.Err = req.f.Write(req.b)
	req.resp <- resp
	p.releaseWriteReq(req)
}

func newSeekReq() *seekRequest {
	return &seekRequest{
		resp: make(chan *SeekResp),
	}
}

type seekRequest struct {
	f *os.File

	offset int64
	whence int

	resp chan *SeekResp
}

func (req *seekRequest) Action() {
	resp := p.acquireSeekResp()
	resp.Ret, resp.Err = req.f.Seek(req.offset, req.whence)
	req.resp <- resp
	p.releaseSeekReq(req)
}

func newSyncReq() *syncRequest {
	return &syncRequest{
		resp: make(chan error),
	}
}

type syncRequest struct {
	f *os.File

	resp chan error
}

func (req *syncRequest) Action() {
	req.resp <- req.f.Sync()
	p.releaseSyncReq(req)
}

func newCloseReq() *closeRequest {
	return &closeRequest{
		resp: make(chan error),
	}
}

type closeRequest struct {
	f *File

	resp chan error
}

func (req *closeRequest) Action() {
	req.resp <- req.f.f.Close()
	p.releaseFile(req.f)
	p.releaseCloseReq(req)
}

func newDelReq() *deleteRequest {
	return &deleteRequest{
		resp: make(chan error),
	}
}

type deleteRequest struct {
	key string

	resp chan error
}

func (req *deleteRequest) Action() {
	req.resp <- os.Remove(req.key)
	p.releaseDelReq(req)
}
