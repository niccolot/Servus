package request

import (
	"errors"
	"fmt"
	"io"

	"Servus/internal/headers"
)

type parserStateType int

const (
	stateInitialized parserStateType = iota
	stateParsingHeaders
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
	ParserState parserStateType
}

func (r *Request) parse(data []byte) (int, error) {
	/*
	* returns the numbers of bytes parsed
	*/
	totalBytesParsed := 0
	for r.ParserState != stateDone {
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
	switch r.ParserState {
	case stateInitialized:
		reqLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		} else if n == 0 {
			// need more data
			return 0, nil
		} else {
			r.RequestLine = *reqLine
			r.ParserState = stateParsingHeaders
			return n, nil
		}
	
	case stateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if done {
			r.ParserState = stateDone
		}

		return n, nil
		
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
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 8)
	readToIndex := 0
	reqStruct := &Request{
		ParserState: stateInitialized,
		Headers: headers.Headers{},
	}
	
	for reqStruct.ParserState != stateDone {
		if readToIndex >= len(buffer) {
			buffer = growBuffer(buffer)
		}

		n, err:= reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if reqStruct.ParserState != stateDone {
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