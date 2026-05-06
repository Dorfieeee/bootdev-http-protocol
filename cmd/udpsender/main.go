package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("Error ResolveUDPAddr(), %v", err)
	}
	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		log.Fatalf("Error DialUDP(), %v", err)
	}
	defer conn.Close()
	buffer := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		str, err := buffer.ReadString(byte('\n'))
		if err != nil {
			log.Printf("Error reading from buffer, error: %v", err)
		}
		if _, err := conn.Write([]byte(str)); err != nil {
			log.Printf("Error writing to conn, error: %v", err)
		}
	}
}
