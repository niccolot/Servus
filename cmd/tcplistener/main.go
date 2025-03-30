package main

import (
	"fmt"
	"log"
	"net"

	"Servus/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Failed to instanciate listener: %v", err)
	}

	defer listener.Close()
	
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %v", err)
		}

		fmt.Println("Connection accepted...")

		req, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatalf("error while processing request: %v", err)
		}

		req.PrintRequest()
		connection.Close()
		fmt.Println("Connection closed...")
	}
}

