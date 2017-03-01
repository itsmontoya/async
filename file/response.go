package file

import "os"

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

func newSeekResp() *SeekResp {
	return &SeekResp{}
}

// SeekResp is the response for seek requests
type SeekResp struct {
	Ret int64
	Err error
}

func newStatResp() *StatResp {
	return &StatResp{}
}

// StatResp is the response for stat requests
type StatResp struct {
	Fi  os.FileInfo
	Err error
}
