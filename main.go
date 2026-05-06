package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Could not open messages.txt: %v", err)
	}
	linesCh := getLinesChannel(file)
	for line := range linesCh {
		fmt.Printf("read: %s\n", line)
	}
}
