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

var closeP = sync.Pool{
	New: func() interface{} {
		return &closeRequest{
			resp: make(chan error),
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

func acquireCloseRequest() (cr *closeRequest) {
	var ok bool
	if cr, ok = closeP.Get().(*closeRequest); !ok {
		panic("invalid pool type")
	}

	return
}

func releaseCloseRequest(cr *closeRequest) {
	cr.f = nil
	closeP.Put(cr)
}
