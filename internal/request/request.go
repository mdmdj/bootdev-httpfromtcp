package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

const bufferSize = 8

type RequestParseState int

const (
	Initialized RequestParseState = iota
	Done
)

func (s RequestParseState) String() string {
	switch s {
	case Initialized:
		return "Initialized"
	case Done:
		return "Done"
	default:
		return "Unknown"
	}
}

type Request struct {
	RequestLine RequestLine
	State       RequestParseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) parse(data []byte) (bytesConsumed int, err error) {
	switch r.State {
	case Initialized:
		var requestLine *RequestLine
		bytesConsumed, requestLine, err = parseRequestLine(string(data))
		if err != nil {
			return
		}

		if bytesConsumed == 0 {
			return
		}

		r.RequestLine = *requestLine
		r.State = Done
		return
	case Done:
		err = errors.New("ERROR in Request.parse: trying to read data in a done state")
		return
	}

	err = errors.New("ERROR in Request.parse: unknown state " + r.State.String())
	return
}

func RequestFromReader(reader io.Reader) (request *Request, err error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	//parsedToIndex := 0

	request = &Request{}
	request.State = Initialized

	for request.State != Done {
		if readToIndex >= cap(buf) {
			newbuf := make([]byte, cap(buf)*2)
			copy(newbuf, buf)
			buf = newbuf
		}

		//Read from the io.Reader into the buffer starting at readToIndex
		var bytesRead int
		bytesRead, err = reader.Read(buf[readToIndex:])
		if err != nil && !errors.Is(err, io.EOF) {
			return
		}

		if bytesRead == 0 {
			continue
		}

		readToIndex += bytesRead

		var bytesParsed int
		bytesParsed, err = request.parse(buf[:readToIndex])
		if err != nil {
			return
		}

		// Remove the data that was parsed successfully from the buffer
		buf = buf[bytesParsed:]

		readToIndex -= bytesParsed

		if errors.Is(err, io.EOF) {
			request.State = Done
			break
		}

		//parsedToIndex = readToIndex

	}

	return
}

func parseRequestLine(content string) (bytesConsumed int, result *RequestLine, err error) {
	//If there is no \r\n yet, there is not enough data to have a request line
	firstRN := strings.Index(content, "\r\n")

	//There is no \r\n yet, so we need to read more data
	if firstRN == -1 {
		bytesConsumed = 0
		return
	}

	//The content starts with \r\n, so it's invalid
	if firstRN == 0 {
		err = errors.New("ERROR in parseRequestLine: invalid request, does not start with a method token")
		return
	}

	requestLine := content[:firstRN]
	requestLineSplit := strings.Split(requestLine, " ")
	if len(requestLineSplit) != 3 {
		err = errors.New("ERROR in parseRequestLine: invalid request line, not 3 parts")
		return
	}

	requestMethod := requestLineSplit[0]
	requestTarget := requestLineSplit[1]
	httpVersionString := requestLineSplit[2]

	// check if anything is missing
	if requestMethod == "" || requestTarget == "" || httpVersionString == "" {
		err = errors.New("ERROR in parseRequestLine: invalid request line, empty parts")
		return
	}

	// check if request method is uppercase
	if requestMethod != strings.ToUpper(requestMethod) {
		err = errors.New("ERROR in parseRequestLine: invalid request method, not uppercase")
		return
	}

	// check if request method contains only alphabetic characters (limit to lower ASCII)
	for _, r := range requestMethod {
		if !unicode.IsLetter(r) || !unicode.Is(unicode.Latin, r) {
			err = errors.New("ERROR in parseRequestLine: invalid request method, contains non-alphabetic characters")
			return
		}
	}

	if httpVersionString != "HTTP/1.1" {
		err = errors.New("ERROR in parseRequestLine: only HTTP/1.1 is supported")
		return
	}

	httpVersionSplit := strings.Split(httpVersionString, "/")

	httpVersion := httpVersionSplit[1]

	result = &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        requestMethod,
	}

	bytesConsumed = len(requestLine)

	return

}
