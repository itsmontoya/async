package aio

import (
	"os"
	"runtime"
)

func newThread(rq chan interface{}) *thread {
	return &thread{rq}
}

type thread struct {
	rq chan interface{}
}

func (t *thread) listen() {
	runtime.LockOSThread()

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

	runtime.UnlockOSThread()
}

func (t *thread) open(r *openRequest) {
	resp := p.acquireOpenResp()
	resp.F, resp.Err = newFile(r, t.rq)
	r.resp <- resp
	p.releaseOpenReq(r)
}

func (t *thread) read(r *readRequest) {
	resp := p.acquireRWResp()
	resp.N, resp.Err = r.f.Read(r.b)
	r.resp <- resp
	p.releaseReadReq(r)
}

func (t *thread) write(r *writeRequest) {
	resp := p.acquireRWResp()
	resp.N, resp.Err = r.f.Write(r.b)
	r.resp <- resp
	p.releaseWriteReq(r)
}

func (t *thread) seek(r *seekRequest) {
	resp := p.acquireSeekResp()
	resp.Ret, resp.Err = r.f.Seek(r.offset, r.whence)
	r.resp <- resp
	p.releaseSeekReq(r)
}

func (t *thread) sync(r *syncRequest) {
	r.resp <- r.f.Sync()
	p.releaseSyncReq(r)
}

func (t *thread) close(r *closeRequest) {
	r.resp <- r.f.f.Close()
	p.releaseFile(r.f)
	p.releaseCloseReq(r)
}

func (t *thread) delete(r *deleteRequest) {
	r.resp <- os.Remove(r.key)
	p.releaseDelReq(r)
}
