package request

import (
	"bytes"
	"fmt"
	"io"
)

type ParserState int

const (
	StateInit ParserState = 0
	StateDone ParserState = 1
)

type Request struct {
	RequestLine RequestLine
	Status      ParserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func newRequest() *Request {
	return &Request{
		Status: 0,
	}
}

var bufferSize = 1024
var sep = "\r\n"
var (
	ErrInvalidRequest     = fmt.Errorf("Invalid request")
	ErrInvalidRequestLine = fmt.Errorf("Invalid request line format")
	ErrInvalidStatus      = fmt.Errorf("Request have invalid status (0 - init, 1 - done)")
)

func (r *Request) parse(data []byte) (int, error) {
	switch r.Status {
	case StateInit:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, fmt.Errorf("Parsing request line error: %w", err)
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *rl
		r.Status = StateDone
		return n, nil
	case StateDone:
		return 0, nil
	default:
		return 0, ErrInvalidStatus
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, bufferSize, bufferSize)
	bufLen := 0

	for request.Status != StateDone {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}

func parseRequestLine(request []byte) (*RequestLine, int, error) {
	idx := bytes.Index(request, []byte(sep))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineString := request[:idx]

	requestLineParts := bytes.Split(requestLineString, []byte(" "))
	if len(requestLineParts) != 3 {
		return nil, 0, ErrInvalidRequestLine
	}

	versionParts := bytes.Split(requestLineParts[2], []byte("/"))
	if len(versionParts) != 2 || string(versionParts[0]) != "HTTP" || string(versionParts[1]) != "1.1" {
		return nil, 0, ErrInvalidRequestLine
	}

	rl := &RequestLine{
		Method:        string(requestLineParts[0]),
		RequestTarget: string(requestLineParts[1]),
		HttpVersion:   string(versionParts[1]),
	}

	if !rl.isValidMethod() {
		return nil, 0, ErrInvalidRequestLine
	}

	return rl, len(requestLineString) + len(sep), nil
}

func (rl *RequestLine) isValidMethod() bool {
	for _, ch := range rl.Method {
		if ch < 'A' || ch > 'Z' {
			return false
		}
	}

	return true
}
