package aio

import (
	"sync"
)

var readP = sync.Pool{
	New: func() interface{} {
		return &readRequest{
			resp: make(chan *RWResp, 1),
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
