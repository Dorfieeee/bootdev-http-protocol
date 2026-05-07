package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	ParseState  int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const CRLF = "\r\n"
const INITIALIZED = 0
const DONE = 1

func (r *Request) parse(data []byte) (int, error) {
	if r.ParseState == DONE {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}

	if r.ParseState == INITIALIZED {
		n, rl, err := parseRequestLine(data)
		if err == nil && n > 0 {
			r.RequestLine = *rl
			r.ParseState = DONE
		}
		return n, err
	}

	return 0, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := Request{}
	var buffer bytes.Buffer
	for r.ParseState != DONE {
		chunk := make([]byte, 8)
		nRead, err := reader.Read(chunk)
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.ParseState = DONE
				break
			}
			fmt.Printf("Error: %v\n", err)
			break
		}

		buffer.Write(chunk[:nRead])
		nParsed, err := r.parse(buffer.Bytes())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}
		if nParsed > 0 {
			tmp := buffer.Bytes()
			buffer.Reset()
			buffer.Write(tmp[nParsed:])
		}
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
