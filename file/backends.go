package file

import "os"

// GetOS is an OpenFile func for the os backend
// Matches the OpenFunc type
func GetOS(name string, flag int, perm os.FileMode) (fi Interface, err error) {
	return os.OpenFile(name, flag, perm)
}

// GetAsync is an OpenFile func for the async backend
// Matches the OpenFunc type
func GetAsync(name string, flag int, perm os.FileMode) (fi Interface, err error) {
	return OpenFile(name, flag, perm)
}

// OpenFunc is the func which produces Interface
type OpenFunc func(name string, flag int, perm os.FileMode) (Interface, error)

// Interface is a file interface
type Interface interface {
	Seek(offset int64, whence int) (ret int64, err error)
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Sync() error
	Stat() (os.FileInfo, error)
	Close() error
}
