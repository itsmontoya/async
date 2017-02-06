package aio

import "os"

func newOpenReq() *openRequest {
	return &openRequest{
		resp: make(chan *OpenResp),
	}
}

type openRequest struct {
	key  string
	flag int
	perm os.FileMode

	resp chan *OpenResp
}

func newReadReq() *readRequest {
	return &readRequest{
		resp: make(chan *RWResp),
	}
}

type readRequest struct {
	f *os.File
	b []byte

	resp chan *RWResp
}

func newWriteReq() *writeRequest {
	return &writeRequest{
		resp: make(chan *RWResp),
	}
}

type writeRequest struct {
	f *os.File
	b []byte

	resp chan *RWResp
}

func newCloseReq() *closeRequest {
	return &closeRequest{
		resp: make(chan error),
	}
}

type closeRequest struct {
	f *File

	resp chan error
}

func newDelReq() *deleteRequest {
	return &deleteRequest{
		resp: make(chan error),
	}
}

type deleteRequest struct {
	key string

	resp chan error
}

func newOpenResp() *OpenResp {
	return &OpenResp{}
}

// OpenResp is a response for open requests
type OpenResp struct {
	F   *File
	Err error
}

func newRWResp() *RWResp {
	return &RWResp{}
}

// RWResp is a response for read/write requests
type RWResp struct {
	N   int
	Err error
}
