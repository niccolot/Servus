package response

import (
	"fmt"
	"io"
	"strconv"

	"Servus/internal/headers"
)

type StatusCode int

const (
	CodeOK StatusCode = 200
	CodeBadRequest StatusCode = 400
	CodeInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case CodeOK:
		statusLine := "HTTP/1.1 " + strconv.Itoa(int(statusCode)) + " OK\r\n" 
		_, err = w.Write([]byte(statusLine))
	
	case CodeBadRequest:
		statusLine := "HTTP/1.1 " + strconv.Itoa(int(statusCode)) + " Bad Request\r\n"
		_, err = w.Write([]byte(statusLine)) 
	
	case CodeInternalServerError:
		statusLine := "HTTP/1.1 " + strconv.Itoa(int(statusCode)) + " Internal Server Error\r\n"
		_, err = w.Write([]byte(statusLine))
	} 

	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.Headers{}
	headers.Add("Content-Length", fmt.Sprint(contentLen))
	headers.Add("Connection", "close")
	headers.Add("Content-Type", "text/plain")

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		headerString := key + ": " + val + "\r\n"
		_, err := w.Write([]byte(headerString))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))

	return err
}
