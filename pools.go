package aio

import (
	"sync"
)

func newPools() *pools {
	var p pools
	p.openReqs.New = func() interface{} {
		return newOpenReq()
	}

	p.readReqs.New = func() interface{} {
		return newReadReq()
	}

	p.writeReqs.New = func() interface{} {
		return newWriteReq()
	}

	p.seekReqs.New = func() interface{} {
		return newSeekReq()
	}

	p.syncReqs.New = func() interface{} {
		return newSyncReq()
	}

	p.closeReqs.New = func() interface{} {
		return newCloseReq()
	}

	p.delReqs.New = func() interface{} {
		return newDelReq()
	}

	p.openResps.New = func() interface{} {
		return newOpenResp()
	}

	p.rwResps.New = func() interface{} {
		return newRWResp()
	}

	p.seekResps.New = func() interface{} {
		return newSeekResp()
	}

	p.files.New = func() interface{} {
		return &File{}
	}

	return &p
}

type pools struct {
	openReqs  sync.Pool
	readReqs  sync.Pool
	writeReqs sync.Pool
	seekReqs  sync.Pool
	syncReqs  sync.Pool
	closeReqs sync.Pool
	delReqs   sync.Pool

	openResps sync.Pool
	rwResps   sync.Pool
	seekResps sync.Pool

	files sync.Pool
}

func (p *pools) acquireOpenReq() (req *openRequest) {
	var ok bool
	if req, ok = p.openReqs.Get().(*openRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireReadReq() (req *readRequest) {
	var ok bool
	if req, ok = p.readReqs.Get().(*readRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireWriteReq() (req *writeRequest) {
	var ok bool
	if req, ok = p.writeReqs.Get().(*writeRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireSeekReq() (req *seekRequest) {
	var ok bool
	if req, ok = p.seekReqs.Get().(*seekRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireSyncReq() (req *syncRequest) {
	var ok bool
	if req, ok = p.syncReqs.Get().(*syncRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireCloseReq() (req *closeRequest) {
	var ok bool
	if req, ok = p.closeReqs.Get().(*closeRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireDelReq() (req *deleteRequest) {
	var ok bool
	if req, ok = p.delReqs.Get().(*deleteRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireOpenResp() (resp *OpenResp) {
	var ok bool
	if resp, ok = p.openResps.Get().(*OpenResp); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireRWResp() (resp *RWResp) {
	var ok bool
	if resp, ok = p.rwResps.Get().(*RWResp); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireSeekResp() (resp *SeekResp) {
	var ok bool
	if resp, ok = p.seekResps.Get().(*SeekResp); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireFile() (f *File) {
	var ok bool
	if f, ok = p.files.Get().(*File); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) releaseOpenReq(req *openRequest) {
	p.openReqs.Put(req)
}

func (p *pools) releaseReadReq(req *readRequest) {
	req.f = nil
	req.b = nil
	p.readReqs.Put(req)
}

func (p *pools) releaseWriteReq(req *writeRequest) {
	req.f = nil
	req.b = nil
	p.writeReqs.Put(req)
}

func (p *pools) releaseSeekReq(req *seekRequest) {
	req.f = nil
	req.offset = 0
	req.whence = 0
	p.seekReqs.Put(req)
}

func (p *pools) releaseSyncReq(req *syncRequest) {
	req.f = nil
	p.syncReqs.Put(req)
}

func (p *pools) releaseCloseReq(req *closeRequest) {
	p.closeReqs.Put(req)
}

func (p *pools) releaseDelReq(req *deleteRequest) {
	p.delReqs.Put(req)
}

func (p *pools) releaseOpenResp(resp *OpenResp) {
	resp.F = nil
	resp.Err = nil
	p.openResps.Put(resp)
}

func (p *pools) releaseRWResp(resp *RWResp) {
	resp.N = 0
	resp.Err = nil
	p.rwResps.Put(resp)
}

func (p *pools) releaseSeekResp(resp *SeekResp) {
	resp.Ret = 0
	resp.Err = nil
	p.seekResps.Put(resp)
}

func (p *pools) releaseFile(f *File) {
	f.f = nil
	f.closed = false
	f.rq = nil
	p.files.Put(f)
}
