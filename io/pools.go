package io

import "sync"

// Global pool for requests and responses
var p = newPools()

func newPools() *pools {
	var p pools

	// Requests
	p.readReqs.New = func() interface{} {
		return newReadReq()
	}

	p.writeReqs.New = func() interface{} {
		return newWriteReq()
	}

	p.closeReqs.New = func() interface{} {
		return newCloseReq()
	}

	p.copyReqs.New = func() interface{} {
		return newCopyReq()
	}

	// Responses
	p.rwResps.New = func() interface{} {
		return newRWResp()
	}

	p.copyResps.New = func() interface{} {
		return newCopyResp()
	}

	return &p
}

type pools struct {
	readReqs  sync.Pool
	writeReqs sync.Pool
	closeReqs sync.Pool
	copyReqs  sync.Pool

	rwResps   sync.Pool
	copyResps sync.Pool
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

func (p *pools) acquireCloseReq() (req *closeRequest) {
	var ok bool
	if req, ok = p.closeReqs.Get().(*closeRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireCopyReq() (req *copyRequest) {
	var ok bool
	if req, ok = p.copyReqs.Get().(*copyRequest); !ok {
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

func (p *pools) acquireCopyResp() (resp *CopyResp) {
	var ok bool
	if resp, ok = p.copyResps.Get().(*CopyResp); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) releaseReadReq(req *readRequest) {
	req.r = nil
	req.b = nil
	p.readReqs.Put(req)
}

func (p *pools) releaseWriteReq(req *writeRequest) {
	req.w = nil
	req.b = nil
	p.writeReqs.Put(req)
}

func (p *pools) releaseCloseReq(req *closeRequest) {
	req.c = nil
	req.b = nil
	p.closeReqs.Put(req)
}

func (p *pools) releaseCopyReq(req *copyRequest) {
	req.w = nil
	req.r = nil
	p.copyReqs.Put(req)
}

func (p *pools) releaseRWResp(resp *RWResp) {
	p.rwResps.Put(resp)
}

func (p *pools) releaseCopyResp(resp *CopyResp) {
	p.copyResps.Put(resp)
}
