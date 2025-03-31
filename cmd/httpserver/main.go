package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"fmt"

	"Servus/internal/headers"
	"Servus/internal/request"
	"Servus/internal/response"
	"Servus/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	if w.Response == nil {
		w.Response = &response.Response{}
	}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		w.Response.Code = 400
		w.Response.Message = "Your problem is not my problem\n"
		w.Response.Headers = headers.Headers{}
		w.Response.Headers.AddOverride("Content-Length", fmt.Sprint(len(w.Response.Message)))

	} else if req.RequestLine.RequestTarget == "/myproblem" {
		w.Response.Code = 500
		w.Response.Message = "Woopsie, my bad\n"
		w.Response.Headers = headers.Headers{}
		w.Response.Headers.AddOverride("Content-Length", fmt.Sprint(len(w.Response.Message)))

	} else {
		w.Response.Code = 200
		w.Response.Message = "All good, frfr\n"
		w.Response.Headers = headers.Headers{}
		w.Response.Headers.AddOverride("Content-Length", fmt.Sprint(len(w.Response.Message)))
	}

	w.WriteResponse()
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