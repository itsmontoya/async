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
}

// Queue will add an item to the request queue
func (a *AIO) Queue(req Actioner) {
	a.rq <- req
}

// Actioner fulfills actions
type Actioner interface {
	Action()
}
