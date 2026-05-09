package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"google.com/Dorfieeee/bootdev-http-protocol/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	parseState  int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const CRLF = "\r\n"

const (
	requestStateParsingInitialised int = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateParsingDone
)

func (r *Request) parse(data []byte) (int, error) {
	switch r.parseState {
	case requestStateParsingInitialised:
		n, rl, err := parseRequestLine(data)
		if err == nil && n > 0 {
			r.RequestLine = *rl
			r.parseState = requestStateParsingHeaders
		}
		return n, err
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if done {
			r.parseState = requestStateParsingBody
			return n, nil
		}
		return n, err
	case requestStateParsingBody:
		contentLength, err := strconv.Atoi(r.Headers.Get("content-length"))
		if err != nil || contentLength == 0 {
			r.parseState = requestStateParsingDone
			return 0, nil
		}
		r.Body = append(r.Body, data...)
		parsedBytes := len(data)
		bodyLength := len(r.Body)
		if bodyLength > contentLength {
			return 0, fmt.Errorf("Body longer than reported content length")
		}
		if bodyLength == contentLength {
			r.parseState = requestStateParsingDone
			return parsedBytes, nil
		}
		return parsedBytes, nil
	case requestStateParsingDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, nil
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := Request{
		RequestLine: RequestLine{},
		Headers:     headers.NewHeaders(),
		parseState:  requestStateParsingInitialised,
	}
	bufferSize := 8
	buffer := make([]byte, bufferSize)
	cursorIdx := 0
	for r.parseState != requestStateParsingDone {
		if cursorIdx+bufferSize >= len(buffer) {
			resizedBuffer := make([]byte, 2*len(buffer))
			copy(resizedBuffer, buffer)
			buffer = resizedBuffer
		}
		nRead, err := reader.Read(buffer[cursorIdx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if r.parseState != requestStateParsingDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", r.parseState, nRead)
				}
				break
			}
			return nil, err
		}
		cursorIdx += nRead
		bufToParse := buffer[:cursorIdx]
		totalBytesParsed := 0
		for r.parseState != requestStateParsingDone {
			bytesParsed, err := r.parse(bufToParse[totalBytesParsed:])
			if err != nil {
				return nil, err
			}
			totalBytesParsed += bytesParsed
			if bytesParsed == 0 {
				break
			}
		}
		cursorIdx -= totalBytesParsed
		copy(buffer, buffer[totalBytesParsed:])
	}
	return &r, nil
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(CRLF))
	if idx == -1 {
		return 0, nil, nil
	}
	rl := strings.Fields(string(data[:idx]))
	if len(rl) != 3 {
		return 0, nil, fmt.Errorf("Invalid request line")
	}
	method := rl[0]
	if match, _ := regexp.MatchString("^[A-Z]+$", method); !match {
		return 0, nil, fmt.Errorf("Invalid HTTP method")
	}
	requestTarget := rl[1]
	httpPart := rl[2]
	versionPart := strings.Split(httpPart, "/")
	if len(versionPart) != 2 {
		return 0, nil, fmt.Errorf("Invalid HTTP version")
	}
	if versionPart[0] != "HTTP" {
		return 0, nil, fmt.Errorf("Invalid HTTP version")
	}
	if versionPart[1] != "1.1" {
		return 0, nil, fmt.Errorf("Invalid HTTP version")
	}
	return len(data[:idx]) + len(CRLF), &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionPart[1],
	}, nil
}
