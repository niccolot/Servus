package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func isUpper(s string) bool {
	for _, c := range s {
		if !unicode.IsUpper(c) && unicode.IsLetter(c) {
			return false
		}
	}

	return true
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqByteSlice, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	reqString := string(reqByteSlice)
	reqParts := strings.Split(reqString, "\r\n")
	reqLine := reqParts[0]
	reqLineParts := strings.Split(reqLine, " ")

	reqStruct := Request{}
	method := reqLineParts[0]
	if !isUpper(method) {
		return nil, fmt.Errorf("method must be capitalized in request line")
	}

	reqStruct.RequestLine.Method = method

	// to do: validate target
	reqStruct.RequestLine.RequestTarget = reqLineParts[1]

	httpVersion := reqLineParts[2]
	if httpVersion != "HTTP/1.1" {
		return nil, fmt.Errorf("http version must be HTTP/1.1")
	}

	reqStruct.RequestLine.HttpVersion = strings.Replace(httpVersion, "HTTP/", "", -1)
	fmt.Println(reqStruct.RequestLine.HttpVersion)

	return &reqStruct, nil
}