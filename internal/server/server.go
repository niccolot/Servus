package server

import (
	"fmt"
	"log"
	"net"

	"sync/atomic"

	"Servus/internal/headers"
	"Servus/internal/request"
	"Servus/internal/response"
)

type Server struct {
	Port int
	closed atomic.Bool
	listener net.Listener
	handlerFunc response.Handler
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
		headers := headers.GetDefaultHeaders(len(err.Error()))
		resp := response.Response{
			Code: response.CodeBadRequest,
			Message: []byte(err.Error()),
			Headers: headers,
		}

		respWriter := response.NewResponseWriter(conn)
		respWriter.Response = &resp
		respWriter.WriteResponse()
		return
	}

	respWriter := response.NewResponseWriter(conn)
	s.handlerFunc(&respWriter, req)
}

func Serve(port int, handler response.Handler) (*Server, error) {
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

