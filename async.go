package async

import (
	"log"
	"runtime"
	"sync"
)

const (
	// WarningInvalidNumThreads is logged when the number of threads are less than one
	WarningInvalidNumThreads = "WARNING: the number of I/O threads is less than 1, setting to 1"
	// WarningTooManyNumThreads is logged when the number of threads specified in New are too much for the current system
	WarningTooManyNumThreads = "WARNING: the number of I/O threads matches or exceeds the number of CPUs"
)

var async = New(1, 1024*32)

// New returns a new asynchronous I/O manager
func New(numThreads, queueLen int) *Async {
	a := Async{
		// Create request queue
		rq: make(chan Actioner, queueLen),
	}

	a.Set(numThreads)
	return &a
}

// Async does stuff
type Async struct {
	mux sync.Mutex

	rq chan Actioner
	ts []*thread
}

// Set will set the selected instance of Async's threads to the numThreads value
// Note: -1 will set the value to the current number of CPUs
func (a *Async) Set(numThreads int) {
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

func (a *Async) openThreads(n int) {
	if n < 1 {
		// Log invalid numThreads warning
		log.Println(WarningInvalidNumThreads)
		return
	}

	a.mux.Lock()
	for i := 0; i < n; i++ {
		// Create new thread
		th := newThread(a.rq)
		a.ts = append(a.ts, th)
		// Call thread.listen within a new goroutine
		go th.listen()
	}
	a.mux.Unlock()
}

func (a *Async) closeThreads(n int) {
	if n < 1 {
		// Log invalid numThreads warning
		log.Println(WarningInvalidNumThreads)
		return
	}

	var i int
	if i = len(a.ts) - 1; n > i {
		n = i
	}

	a.mux.Lock()
	for {
		th := a.ts[i]
		th.Close()
		a.ts = a.ts[:i]

		if n == 0 {
			break
		}

		n--
		i--
	}
	a.mux.Unlock()
}

// Queue will add an item to the request queue
func (a *Async) Queue(req Actioner) {
	a.rq <- req
}

// Actioner fulfills actions
type Actioner interface {
	Action()
}

// Set is the exported Set func for the global async
// Note: -1 will set the value to the current number of CPUs
func Set(numThreads int) {
	async.Set(numThreads)
}

// Queue is the exported Queue func for the global Async
func Queue(req Actioner) {
	async.Queue(req)
}

func popThread(ts []*thread, n int) []*thread {
	return append(ts[:n], ts[n+1:]...)
}

// QueueFn is a queue function for sending requests
type QueueFn func(req Actioner)
