package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const CRLF = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	delim := ":"
	idx := bytes.Index(data, []byte(CRLF))
	if idx == -1 {
		return 0, false, nil
	}
	fieldLine := string(data[:idx])
	if len(fieldLine) == 0 {
		return len(CRLF), true, nil
	}
	fieldName, fieldValue, ok := strings.Cut(fieldLine, delim)
	if !ok {
		return 0, false, fmt.Errorf("Malformed field line format")
	}
	if len(strings.TrimSpace(fieldName)) != len(fieldName) {
		return 0, false, fmt.Errorf("Malformed field name format")
	}
	fieldName = strings.ToLower(strings.TrimSpace(fieldName))
	if !validFieldName(fieldName) {
		return 0, false, fmt.Errorf("Field name contains unrecognized characters")
	}
	fieldValue = strings.TrimSpace(fieldValue)
	if len(fieldName) == 0 || len(fieldValue) == 0 {
		return 0, false, fmt.Errorf("Malformed field line format")
	}
	h.Set(fieldName, fieldValue)
	return len(data[:idx]) + len(CRLF), false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	prev, exists := h[key]
	if exists {
		h[key] = prev + ", " + value
	} else {
		h[key] = value
	}
}

func (h Headers) Get(key string) string {
	return h[strings.ToLower(key)]
}

const validCharsMap = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-.^_`|~"

func validFieldName(s string) bool {
	if len(s) < 1 {
		return false
	}
	for _, ch := range s {
		if idx := strings.IndexRune(validCharsMap, ch); idx == -1 {
			return false
		}
	}
	return true
}
