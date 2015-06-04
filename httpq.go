package httpq

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Httpq struct {
	q     Queue
	debug bool
}

func NewHttpq(q Queue, debug bool) *Httpq {
	return &Httpq{q, debug}
}

func (h *Httpq) PushRequest(r *http.Request) error {
	requestBytes, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}

	if h.debug {
		fmt.Printf("pushed request:\n%s", string(requestBytes))
	}

	if err := h.q.Push(requestBytes); err != nil {
		return err
	}

	return nil
}

func (h *Httpq) PopRequestBytes() ([]byte, error) {
	requestBytes, err := h.q.Pop()
	if err != nil {
		return nil, err
	}

	if h.debug {
		fmt.Printf("popped request:\n%s", string(requestBytes))
	}

	return requestBytes, nil
}

func (h *Httpq) PopRequest() (*http.Request, error) {
	requestBytes, err := h.PopRequestBytes()
	if err != nil {
		return nil, err
	}

	if requestBytes == nil {
		return nil, nil
	}

	buf := bytes.NewBuffer(requestBytes)
	requestBytesReader := bufio.NewReader(buf)
	request, err := http.ReadRequest(requestBytesReader)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (h *Httpq) Size() (uint64, error) {
	return h.q.Size()
}
