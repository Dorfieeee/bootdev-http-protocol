package response

import (
	"fmt"
	"io"

	"google.com/Dorfieeee/bootdev-http-protocol/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

const CRLF = "\r\n"
const SP = " "

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := fmt.Sprintf("%s%s%d%s", "HTTP/1.1", SP, statusCode, SP)
	switch statusCode {
	case StatusOK:
		statusLine += "OK"
	case StatusBadRequest:
		statusLine += "Bad Request"
	case StatusInternalServerError:
		statusLine += "Internal Status Error"
	}
	_, err := fmt.Fprintf(w, "%s%s", statusLine, CRLF)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	// Content-Length (Set to the given size)
	// Connection (Set to close because we're not doing keep-alive's yet)
	// Content-Type (Set to text/plain)
	header := headers.NewHeaders()
	header.Add("Content-Length", fmt.Sprintf("%d", contentLen))
	header.Add("Connection", "close")
	header.Add("Content-Type", "text/plain")
	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k := range headers {
		_, err := fmt.Fprintf(w, "%s:%s%s%s", k, SP, headers.Get(k), CRLF)
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte(CRLF))
	if err != nil {
		return err
	}
	return nil
}
