package reader

import (
	"io"
	"sync"
)

// Global pool for requests and responses
// TODO: Decide if we want to bring the pools to the AIO-level, and give AIO's the ability to utilize their own pools
var p = newPools()

func newPools() *pools {
	var p pools

	p.requests.New = func() interface{} {
		return newReadReq()
	}

	p.responses.New = func() interface{} {
		return newResp()
	}

	return &p
}

type pools struct {
	requests  sync.Pool
	responses sync.Pool
}

func (p *pools) acquireRequest() (req *readRequest) {
	var ok bool
	if req, ok = p.requests.Get().(*readRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) acquireResponse() (resp *RWResp) {
	var ok bool
	if resp, ok = p.responses.Get().(*RWResp); !ok {
		panic("invalid pool type")
	}

	return
}

func (p *pools) releaseRequest(req *readRequest) {
	p.requests.Put(req)
}

func (p *pools) releaseResponse(resp *RWResp) {
	p.responses.Put(resp)
}

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
	resp := newResp()
	resp.N, resp.Err = req.r.Read(req.b)
	req.resp <- resp
	p.releaseRequest(req)
}

func newResp() *RWResp {
	return &RWResp{}
}

// RWResp is a response for a read request
type RWResp struct {
	N   int
	Err error
}
