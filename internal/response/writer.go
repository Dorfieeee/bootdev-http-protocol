package response

import (
	"fmt"
	"io"

	"google.com/Dorfieeee/bootdev-http-protocol/internal/headers"
)

type writerState int

const (
	writerStateInitialised writerState = iota
	writerStateHeaders
	writerStateBody
	writerStateDone
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateInitialised {
		return fmt.Errorf("Wrong order for writing. Current write state: %d", w.state)
	}
	defer func() { w.state = writerStateHeaders }()
	return WriteStatusLine(w.writer, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateHeaders {
		return fmt.Errorf("Wrong order for writing. Current write state: %d", w.state)
	}
	defer func() { w.state = writerStateBody }()
	return WriteHeaders(w.writer, headers)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("Wrong order for writing. Current write state: %d", w.state)
	}
	defer func() { w.state = writerStateDone }()
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.state != writerStateBody {
		return 0, fmt.Errorf("Wrong order for writing. Current write state: %d", w.state)
	}

}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	defer func() { w.state = writerStateDone }()
	return w.writer.Write([]byte("0\r\n\r\n"))
}
