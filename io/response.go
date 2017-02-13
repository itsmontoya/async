package io

func newRWResp() *RWResp {
	return &RWResp{}
}

// RWResp is a response for a read or write request
type RWResp struct {
	N   int
	Err error
}

func newCopyResp() *CopyResp {
	return &CopyResp{}
}

// CopyResp is a response for a copy request
type CopyResp struct {
	N   int64
	Err error
}
