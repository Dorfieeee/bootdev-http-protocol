package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"google.com/Dorfieeee/bootdev-http-protocol/internal/request"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string, 8)
	go func() {
		line := strings.Builder{}
		defer close(ch)
		defer f.Close()
		for {
			chunk := make([]byte, 8)
			n, err := f.Read(chunk)
			if err != nil {
				if line.Len() > 0 {
					ch <- line.String()
					line.Reset()
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("Error: %v\n", err)
				break
			}
			parts := strings.Split(string(chunk[:n]), "\n")
			for i, part := range parts {
				line.WriteString(part)
				if i < len(parts)-1 {
					ch <- line.String()
					line.Reset()
				}
			}
		}
	}()
	return ch
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("Could not open listener: %v", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}
		fmt.Println("Connection establised")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error %v\n", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %v\n- Target: %v\n- Version: %v\n", r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)
		fmt.Println("Connection closed")
	}
}
