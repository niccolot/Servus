package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"Servus/internal/request"
	"Servus/internal/response"
	"Servus/internal/server"
	"Servus/internal/html"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {
	if w.Response == nil {
		w.Response = &response.Response{}
	}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		//w.Response.Code = 400
		//w.Response.Message = []byte("Your problem is not my problem\n")
		fileName := "cmd/httpserver/assets/req_badRequest.html"
		//status, body, err := html.ExtractTitleAndBodyFromFile(fileName)
		//if err != nil {
		//	log.Fatalf("failed to parse html")
		//}
		//code, err := html.ExtractCode(status)
		//if err != nil {
		//	log.Fatalf("failed to parse status code")
		//}
//
		//w.Response.Code = response.StatusCode(code)
		//w.Response.Message = body
		//w.Response.Headers = headers.Headers{}
		//w.Response.Headers.AddOverride("Content-Type", "text/html")
		//w.Response.Headers.AddOverride("Content-Length", fmt.Sprint(len(w.Response.Message)))
		err := html.WriteResponse(w, fileName)
		if err != nil {
			log.Fatalf("failed to write response: %v", err)
		}

	} else if req.RequestLine.RequestTarget == "/myproblem" {
		fileName := "cmd/httpserver/assets/req_internalErr.html"
		err := html.WriteResponse(w, fileName)
		if err != nil {
			log.Fatalf("failed to write response: %v", err)
		}

	} else {
		fileName := "cmd/httpserver/assets/req_success.html"
		err := html.WriteResponse(w, fileName)
		if err != nil {
			log.Fatalf("failed to write response: %v", err)
		}
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