package server

import (
	"bytes"
	"fmt"
	"log"
	"net"

	"sync/atomic"

	"Servus/internal/request"
	"Servus/internal/response"
)

type Server struct {
	Port int
	closed atomic.Bool
	listener net.Listener
	handlerFunc Handler
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return 
			}

			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	err := s.listener.Close()
	s.closed.Store(true)

	return err
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		handlerErr := HandlerError{
			Code: response.CodeBadRequest,
			Message: err.Error(),
		}

		handlerErr.WriteHandlerError(conn)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	handlerErr := s.handlerFunc(buffer, req)
	if handlerErr != nil {
		handlerErr.WriteHandlerError(conn)
		return
	}

	data := buffer.Bytes()
	response.WriteStatusLine(conn, response.CodeOK)
	respHeaders := response.GetDefaultHeaders(len(data))
	response.WriteHeaders(conn, respHeaders)
	conn.Write(data)
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := Server{
		Port: port,
		listener: l,
		handlerFunc: handler,
	}

	server.closed.Store(false)
	go server.listen()

	return &server, err
}

