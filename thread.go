package aio

import (
	"runtime"
	"sync/atomic"

	"github.com/missionMeteora/toolkit/errors"
	"sync"
)

func newThread(rq chan Actioner) *thread {
	return &thread{rq: rq}
}

type thread struct {
	rq chan Actioner
	wg sync.WaitGroup

	closed int32
}

func (t *thread) listen() {
	t.wg.Add(1)
	runtime.LockOSThread()

	for req := range t.rq {
		req.Action()

		if t.isClosed() {
			break
		}
	}

	runtime.UnlockOSThread()
	t.wg.Done()
}

func (t *thread) isClosed() bool {
	return atomic.LoadInt32(&t.closed) == 1
}

func (t *thread) Close() (err error) {
	if !atomic.CompareAndSwapInt32(&t.closed, 0, 1) {
		return errors.ErrIsClosed
	}

	t.wg.Wait()
	return
}
