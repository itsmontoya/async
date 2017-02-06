package aio

import (
	"log"
	"os"
	"runtime"
)

// Global pool for requests and responses
// TODO: Decide if we want to bring the pools to the AIO-level, and give AIO's the ability to utilize their own pools
var p = newPools()

const (
	// WarningInvalidNumThreads is logged when the number of threads are negativ
	WarningInvalidNumThreads = "WARNING: the number of I/O threads is less than 1, setting to 1"
	// WarningTooManyNumThreads is logged when the number of threads specified in New are too much for the current system
	WarningTooManyNumThreads = "WARNING: the number of I/O threads matches or exceeds the number of CPUs"
)

// New returns a new asynchronous I/O manager
func New(numThreads int) *AIO {
	a := AIO{
		// Create request queue
		rq: make(chan interface{}, 1024*32),
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
	rq chan interface{}
}

// Open will open a new file for reading
func (a *AIO) Open(key string) (f *File, err error) {
	return a.OpenFile(key, os.O_RDONLY, 0)
}

// OpenFile will open a new file with flag and perm
func (a *AIO) OpenFile(key string, flag int, perm os.FileMode) (f *File, err error) {
	// Call OpenFileAsync and wait for the channel to return
	resp := <-a.OpenFileAsync(key, flag, perm)

	// Set f and err from response
	f = resp.F
	err = resp.Err

	// Return response to the pool
	p.releaseOpenResp(resp)
	return
}

// OpenFileAsync will open a new file with flag and perm asynchronously
func (a *AIO) OpenFileAsync(key string, flag int, perm os.FileMode) <-chan *OpenResp {
	// Acquire open request from pool
	or := p.acquireOpenReq()

	// Set open request values
	or.key = key
	or.flag = flag
	or.perm = perm

	// Send request to queue
	a.rq <- or
	return or.resp
}

// Delete will delete a file
func (a *AIO) Delete(key string) <-chan error {
	// Acquire delete request from pool
	dr := p.acquireDelReq()

	// Set delete request key
	dr.key = key

	// Send request to queue
	a.rq <- dr
	return dr.resp
}
