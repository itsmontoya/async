package aio

import (
	"sync"
)

var readP = sync.Pool{
	New: func() interface{} {
		return &readRequest{
			resp: make(chan *RWResp),
		}
	},
}

var openP = sync.Pool{
	New: func() interface{} {
		return &openRequest{
			resp: make(chan *OpenResp),
		}
	},
}

func acquireReadRequest() (rr *readRequest) {
	var ok bool
	if rr, ok = readP.Get().(*readRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func releaseReadRequest(rr *readRequest) {
	rr.f = nil
	rr.b = nil
	readP.Put(rr)
}

func acquireOpenRequest() (or *openRequest) {
	var ok bool
	if or, ok = openP.Get().(*openRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func releaseOpenRequest(or *openRequest) {
	openP.Put(or)
}
