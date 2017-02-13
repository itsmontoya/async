package file

import (
	"os"

	"github.com/itsmontoya/aio"
)

// New will return a new manager
func New(a *aio.AIO) *Manager {
	return &Manager{a}
}

// Manager will manage files
type Manager struct {
	a *aio.AIO
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
	req.a = m.a

	// Send request to queue
	m.a.Queue(req)
	return req.resp
}

// Delete will delete a file
func (m *Manager) Delete(key string) <-chan error {
	// Acquire delete request from pool
	req := p.acquireDelReq()

	// Set delete request key
	req.key = key

	// Send request to queue
	m.a.Queue(req)
	return req.resp
}
