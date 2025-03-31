package response

import (
	"fmt"
	"net"
	"strconv"

	"Servus/internal/headers"
	"Servus/internal/request"
)

type StatusCode int

const (
	CodeOK StatusCode = 200
	CodeBadRequest StatusCode = 400
	CodeInternalServerError StatusCode = 500
)

type Response struct {
	Code StatusCode
	Message []byte
	Headers headers.Headers
}

type Handler func(w *Writer, req *request.Request)

type WriterStatus int

const (
	StatusWriteResponseLine WriterStatus = iota
	StatusWriteHeaders
	StatusWriteBody
	StatusDone
)

type Writer struct {
	Status WriterStatus 
	Response *Response
	Connection net.Conn
}

func NewResponseWriter(conn net.Conn) Writer {
	return Writer{
		Status: StatusWriteResponseLine,
		Connection: conn,
	}
}

func (w *Writer) WriteStatusLine(code StatusCode) error {
	var err error
	if w.Status != StatusWriteResponseLine {
		return fmt.Errorf("invalid response writer status")
	}

	switch w.Response.Code {
	case CodeOK:
		statusLine := "HTTP/1.1 " + strconv.Itoa(int(w.Response.Code)) + " OK\r\n"
		_, err = w.Connection.Write([]byte(statusLine))
	
	case CodeBadRequest:
		statusLine := "HTTP/1.1 " + strconv.Itoa(int(w.Response.Code)) + " Bad Request\r\n"
		_, err = w.Connection.Write([]byte(statusLine))

	case CodeInternalServerError:
		statusLine := "HTTP/1.1 " + strconv.Itoa(int(w.Response.Code)) + " Internal Server Error\r\n"
		_, err = w.Connection.Write([]byte(statusLine))
	}

	w.Status = StatusWriteHeaders

	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.Status != StatusWriteHeaders {
		return fmt.Errorf("invalid response writer status")
	}

	for key, val := range headers {
		headerString := key + ": " + val + "\r\n"
		_, err := w.Connection.Write([]byte(headerString))
		if err != nil {
			return err
		}
	}

	_, err := w.Connection.Write([]byte("\r\n"))

	w.Status = StatusWriteBody

	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.Status != StatusWriteBody {
		return 0, fmt.Errorf("invalid response writer status")
	}

	n, err := w.Connection.Write(p)
	w.Status = StatusDone

	return n, err
}

func (w *Writer) WriteResponse() (int, error) {
	err := w.WriteStatusLine(w.Response.Code)
	if err != nil {
		return 0, err
	}

	err = w.WriteHeaders(w.Response.Headers)
	if err != nil {
		return 0, err
	}

	n, err := w.WriteBody(w.Response.Message) 
	if err != nil {
		return 0, err
	}

	return n, nil
}
/*
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
*/

