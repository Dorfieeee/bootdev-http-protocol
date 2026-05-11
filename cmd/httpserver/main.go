package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"google.com/Dorfieeee/bootdev-http-protocol/internal/request"
	"google.com/Dorfieeee/bootdev-http-protocol/internal/response"
	"google.com/Dorfieeee/bootdev-http-protocol/internal/server"
)

const port = 42069

func handleHTMLRoute(w *response.Writer, message string, status response.StatusCode) {
	header := response.GetDefaultHeaders(len(message))
	header.Set("Content-Type", "text/html")

	b := &bytes.Buffer{}
	fmt.Fprintf(b, "<html><head><title>")
	fmt.Fprintf(b, "%d %s", status, http.StatusText(int(status)))
	fmt.Fprintf(b, "<title><head>")
	fmt.Fprintf(b, "<body><h1>%s</h1><p>%s</p></body>", http.StatusText(int(status)), message)
	fmt.Fprintf(b, "</html>")

	header.Set("Content-Length", fmt.Sprintf("%d", b.Len()))

	w.WriteStatusLine(status)
	w.WriteHeaders(header)
	w.WriteBody(b.Bytes())
}

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handleHTMLRoute(w, "Your request honestly kinda sucked.", response.StatusBadRequest)
		return
	case "/myproblem":
		handleHTMLRoute(w, "Okay, you know what? This one is on me.", response.StatusInternalServerError)
		return
	}
	handleHTMLRoute(w, "Your request was an absolute banger.", response.StatusOK)
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
