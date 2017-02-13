package aio

import (
	"log"
	"runtime"
)

const (
	// WarningInvalidNumThreads is logged when the number of threads are less than one
	WarningInvalidNumThreads = "WARNING: the number of I/O threads is less than 1, setting to 1"
	// WarningTooManyNumThreads is logged when the number of threads specified in New are too much for the current system
	WarningTooManyNumThreads = "WARNING: the number of I/O threads matches or exceeds the number of CPUs"
)

var aio = New(1)

// New returns a new asynchronous I/O manager
func New(numThreads int) *AIO {
	a := AIO{
		// Create request queue
		rq: make(chan Actioner, 1024*32),
	}

	if numThreads < 1 {
		// Log invalid numThreads warning
		log.Println(WarningInvalidNumThreads)
		// numThreads is an invalid value, set to 1
		numThreads = 1
	}

	if numThreads >= runtime.NumCPU() {
		// Log too many numThreads warning
		log.Println(WarningTooManyNumThreads)
	}

	for i := 0; i < numThreads; i++ {
		// Create new thread
		t := newThread(a.rq)
		// Call thread.listen within a new goroutine
		go t.listen()
	}

	return &a
}

// AIO does stuff
type AIO struct {
	rq chan Actioner
	ts []*thread
}

// Set will set the selected instance of AIO's threads to the numThreads value
// Note: -1 will set the value to the current number of CPUs
func (a *AIO) Set(numThreads int) {
	if numThreads == -1 {
		numThreads = runtime.NumCPU()
	} else if numThreads < 0 {
		numThreads = 1
	}

	if delta := numThreads - len(a.ts); delta == 0 {
		return
	} else if delta < 0 {
		a.closeThreads(-delta)
	} else {
		a.openThreads(delta)
	}
}

func (a *AIO) openThreads(n int) {
	if n < 1 {
		// Log invalid numThreads warning
		log.Println(WarningInvalidNumThreads)
		return
	}

	for i := 0; i < n; i++ {
		// Create new thread
		th := newThread(a.rq)
		// Call thread.listen within a new goroutine
		go th.listen()
	}

}

func (a *AIO) closeThreads(n int) {
	if n > len(a.ts) {
		n = len(a.ts)
	}

	for ; n > 0; n-- {
		th := a.ts[0]
		th.Close()
		popThread(a.ts, 0)
	}
}

// Queue will add an item to the request queue
func (a *AIO) Queue(req Actioner) {
	a.rq <- req
}

// Actioner fulfills actions
type Actioner interface {
	Action()
}

// Set is the exported Set func for the global aio
// Note: -1 will set the value to the current number of CPUs
func Set(numThreads int) {
	aio.Set(numThreads)
}

// Queue is the exported Queue func for the global AIO
func Queue(req Actioner) {
	aio.Queue(req)
}

func popThread(ts []*thread, n int) []*thread {
	return append(ts[:n], ts[n+1:]...)
}

// QueueFn is a queue function for sending requests
type QueueFn func(req Actioner)
