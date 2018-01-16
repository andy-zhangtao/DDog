package io_ut

import (
	"bytes"
	"io"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	//we don't actually have to do anything here, since the buffer is
	//and the error is initialized to no-error
	return
}

func GetReadCloser(str string) io.ReadCloser {
	cb := &ClosingBuffer{bytes.NewBufferString(str)}
	var rc io.ReadCloser
	rc = cb

	return rc
}
