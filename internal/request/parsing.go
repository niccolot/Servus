package request

import (
	"fmt"
	"strings"
)

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

//func growBuffer(buffer []byte) []byte {
//	/*
//	* doubles the buffer size
//	*/
//	prevSize := cap(buffer)
//	newBuffer := make([]byte, prevSize, 2*prevSize)
//	copy(newBuffer, buffer)
//
//	return newBuffer
//}
//
//func shrinkBuffer(bufffer []byte, parsedBytes int) []byte {
//	/*
//	* removes the parsed bytes from the buffer
//	*/
//	prevSize := cap(bufffer)
//	newBuffer := make([]byte, prevSize - parsedBytes)
//	copy(newBuffer, bufffer[parsedBytes:])
//
//	return newBuffer
//}

func parseRequestLine(reqLineByteSlice []byte) (*RequestLine, int, error) {
	reqString := string(reqLineByteSlice)
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