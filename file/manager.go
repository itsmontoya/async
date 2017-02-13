package file

import (
	"os"

	"github.com/itsmontoya/aio"
)

var mngr = New(nil)

// New will return a new manager
func New(a *aio.AIO) *Manager {
	var m Manager
	if a == nil {
		m.qfn = aio.Queue
	} else {
		m.qfn = a.Queue
	}

	return &m
}

// Manager will manage files
type Manager struct {
	qfn func(aio.Actioner)
}

// Open will open a new file for reading
func (m *Manager) Open(key string) (f *File, err error) {
	return m.OpenFile(key, os.O_RDONLY, 0)
}

// OpenFile will open a new file with flag and perm
func (m *Manager) OpenFile(key string, flag int, perm os.FileMode) (f *File, err error) {
	// Call OpenFileAsync and wait for the channel to return
	resp := <-m.OpenFileAsync(key, flag, perm)

	// Set f and err from response
	f = resp.F
	err = resp.Err

	// Return response to the pool
	p.releaseOpenResp(resp)
	return
}

// OpenFileAsync will open a new file with flag and perm asynchronously
func (m *Manager) OpenFileAsync(key string, flag int, perm os.FileMode) <-chan *OpenResp {
	// Acquire open request from pool
	req := p.acquireOpenReq()

	// Set open request values
	req.key = key
	req.flag = flag
	req.perm = perm
	req.qfn = m.qfn

	// Send request to queue
	m.qfn(req)
	return req.resp
}

// Delete will delete a file
func (m *Manager) Delete(key string) <-chan error {
	// Acquire delete request from pool
	req := p.acquireDelReq()

	// Set delete request key
	req.key = key

	// Send request to queue
	m.qfn(req)
	return req.resp
}

// Open is the exported Open func for the global Manager
func Open(key string) (f *File, err error) {
	return mngr.Open(key)
}

// OpenFile is the exported OpenFile func for the global Manager
func OpenFile(key string, flag int, perm os.FileMode) (f *File, err error) {
	return mngr.OpenFile(key, flag, perm)
}

// OpenFileAsync is the exported OpenFileAsync func for the global Manager
func OpenFileAsync(key string, flag int, perm os.FileMode) <-chan *OpenResp {
	return mngr.OpenFileAsync(key, flag, perm)
}

// Delete is the exported Delete func for the global Manager
func Delete(key string) <-chan error {
	return mngr.Delete(key)
}
