package request

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"Servus/internal/headers"
)

type parserStateType int

const (
	stateInitialized parserStateType = iota
	stateParsingHeaders
	stateParsingBody
	stateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	parserState parserStateType
	Body []byte
	contentLength int
}

func (r *Request) parse(data []byte) (int, error) {
	/*
	* returns the numbers of bytes parsed
	*/
	totalBytesParsed := 0
	for r.parserState != stateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.parserState {
	case stateInitialized:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		} else if n == 0 {
			// need more data
			return 0, nil
		} else {
			r.RequestLine = *reqLine
			r.parserState = stateParsingHeaders
			return n, nil
		}
	
	case stateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			contentLength, ok := r.Headers.Get("Content-Length")
			if ok {
				cLength, err := strconv.Atoi(contentLength)
				if err != nil {
					return 0, fmt.Errorf("failed to parse 'Content-length' header value: %v", err)
				}
				r.Body = make([]byte, 0, cLength)
				r.contentLength = cLength
			}

			r.parserState = stateParsingBody
		}

		return n, nil

	case stateParsingBody:
		_, ok := r.Headers.Get("Content-Length")
		if !ok {
			// in this implementation it is assumed that if 
			// a request has a body it must also contain a 
			// 'Content-Length' header
			r.parserState = stateDone
			return 0, nil
		}

		r.Body = append(r.Body, data...)		
		if len(r.Body) > r.contentLength {
			return 0, fmt.Errorf("actual body size greater than 'Content-Length' header value")
		}

		if len(r.Body) == r.contentLength {
			r.parserState = stateDone
		}


		return len(data), nil
		
	case stateDone:
		return 0, fmt.Errorf("already done parsing")
	
	default:
		return 0, fmt.Errorf("unknown parser state")
	}
}

func (r *Request) PrintRequest() {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", r.RequestLine.Method)
	fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	fmt.Println("Headers:")
	for header, value := range r.Headers {
		fmt.Printf("- %s: %s\n", header, value)
	}
	fmt.Println("Body:")
	fmt.Println(string(r.Body))
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	const buffSize = 8
	buffer := make([]byte, buffSize)
	readToIndex := 0
	reqStruct := &Request{
		parserState: stateInitialized,
		Headers: headers.Headers{},
		Body: make([]byte, 0),
	}
	
	for reqStruct.parserState != stateDone {
		if readToIndex >= len(buffer) {
			buffer = growBuffer(buffer)
		}

		n, err:= reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if reqStruct.parserState != stateDone {
					return nil, fmt.Errorf("incomplete request")
				}

				break
			}

			return nil, err
		}

		readToIndex += n
		parsedBytes, err := reqStruct.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		// overwrite the already parsed data with the data
		// to be processed to avoid growing the buffer too much
		copy(buffer, buffer[parsedBytes:])
		readToIndex -= parsedBytes
	}

	return reqStruct, nil
}

func isMethodValid(method string) bool {
	httpMethods := map[string]bool {
		"GET": true,
		"HEAD": true,
		"POST": true,
		"PUT": true,
		"DELETE": true,
		"TRACE": true,
	}

	_, ok := httpMethods[method]

	return ok
}

func growBuffer(buffer []byte) []byte {
	/*
	* doubles the buffer size
	*/
	newBuffer := make([]byte, 2 * len(buffer))
	copy(newBuffer, buffer)

	return newBuffer
}

func parseRequestLine(reqLineByteSlice []byte) (*RequestLine, int, error) {
	reqString := string(reqLineByteSlice)
	
	// end of request line not yet reached
	if !strings.Contains(reqString, "\r\n") {
		return nil, 0, nil
	}

	reqParts := strings.Split(reqString, "\r\n")
	reqLine := reqParts[0]
	reqLineParts := strings.Split(reqLine, " ")
	if len(reqLineParts) != 3 {
		return nil, 0, fmt.Errorf("invalid request line")
	}

	reqLineStruct := RequestLine{}
	method := reqLineParts[0]
	if !isMethodValid(method) {
		return nil, 0, fmt.Errorf("invalid http method")
	}

	reqLineStruct.Method = method

	target := reqLineParts[1]
	if strings.Contains(target, " ") {
		return nil, 0, fmt.Errorf("invalid target")
	}

	reqLineStruct.RequestTarget = target

	httpVersion := reqLineParts[2]
	if httpVersion != "HTTP/1.1" {
		return nil, 0, fmt.Errorf("http version must be HTTP/1.1")
	}

	reqLineStruct.HttpVersion = strings.ReplaceAll(httpVersion, "HTTP/", "")

	// +2 for the \r\n chars
	return &reqLineStruct, len(reqLine) + 2, nil
}