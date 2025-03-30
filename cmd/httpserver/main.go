package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Servus/internal/request"
	"Servus/internal/server"
)

const port = 42069

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	handlerErr := server.HandlerError{}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handlerErr.Code = 400
		handlerErr.Message = "Your problem is not my problem\n"
		handlerErr.WriteHandlerError(w)
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		handlerErr.Code = 500
		handlerErr.Message = "Woopsie, my bad\n"
		handlerErr.WriteHandlerError(w)
	} else {
		handlerErr.Code = 200
		handlerErr.Message = "All good, frfr\n"
		handlerErr.WriteHandlerError(w)
	}

	return &handlerErr
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}