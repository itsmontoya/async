package aio

import (
	"runtime"
)

func newThread(rq chan Actioner) *thread {
	return &thread{rq}
}

type thread struct {
	rq chan Actioner
}

func (t *thread) listen() {
	runtime.LockOSThread()

	for req := range t.rq {
		req.Action()
	}

	runtime.UnlockOSThread()
}
