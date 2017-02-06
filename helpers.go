package aio

func newOpenReq() *openRequest {
	return &openRequest{
		resp: make(chan *OpenResp),
	}
}

func newReadReq() *readRequest {
	return &readRequest{
		resp: make(chan *RWResp),
	}
}

func newWriteReq() *writeRequest {
	return &writeRequest{
		resp: make(chan *RWResp),
	}
}

func newDelReq() *deleteRequest {
	return &deleteRequest{
		resp: make(chan error),
	}
}

func newCloseReq() *closeRequest {
	return &closeRequest{
		resp: make(chan error),
	}
}

func newOpenResp() *OpenResp {
	return &OpenResp{}
}

func newRWResp() *RWResp {
	return &RWResp{}
}
