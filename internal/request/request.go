package request

import (
	"errors"
	"fmt"
	"io"
)

type parserStateType int

const (
	stateInitialized parserStateType = iota
	stateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	ParserState parserStateType
}

func (r *Request) parse(data []byte) (int, error) {
	/*
	* returns the numbers of bytes parsed
	*/
	switch r.ParserState {
	case stateInitialized:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		} else if n == 0 {
			return n, nil
		} else {
			r.RequestLine = *reqLine
			r.ParserState = stateDone
			return n, nil
		}
	
	case stateDone:
		return 0, fmt.Errorf("already done parsing")
	
	default:
		return 0, fmt.Errorf("unknown parser state")
	}
}

func (r *Request) PrintRequestLine() {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", r.RequestLine.Method)
	fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 8)
	readToIndex := 0
	reqStruct := Request{
		ParserState: stateInitialized,
	}
	
	for reqStruct.ParserState != stateDone {
		if len(buffer) == cap(buffer) {
			buffer = growBuffer(buffer)
		}

		n, err:= reader.Read(buffer[readToIndex:])
		if errors.Is(err, io.EOF) {
			reqStruct.ParserState = stateDone
			break
		}

		readToIndex += n
		parsedBytes, err := reqStruct.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		buffer = shrinkBuffer(buffer, parsedBytes)
		readToIndex -= parsedBytes
	}

	return &reqStruct, nil
}