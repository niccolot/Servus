package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"Servus/internal/request"
)

func getLinesChannel(f io.ReadCloser) <- chan string {
	ch := make(chan string)

	go func ()  {
		defer f.Close()
		defer close(ch)
		
		byteSlice := make([]byte, 8)
		currLine := ""
		n, _ := f.Read(byteSlice)
		currLine += string(byteSlice[:n])
		for n == 8 {
			parts := strings.Split(currLine, "\n")
			if len(parts) > 1 {
				for i:= 0; i < len(parts) - 1; i++ {
					ch <- parts[i]
				}
				currLine = ""
				currLine += parts[len(parts) - 1]
			}
			n, _ = f.Read(byteSlice)
			currLine += string(byteSlice[:n])
		}

		if currLine != "" {
			ch <- currLine
		}
	} ()

	return ch
}

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
		
		//ch := getLinesChannel(connection)
		//for line := range ch {
		//	fmt.Printf("%s\n", line)
		//}

		req, err := request.RequestFromReader(connection)
		req.PrintRequestLine()

		fmt.Println("Connection closed...")
	}
}

