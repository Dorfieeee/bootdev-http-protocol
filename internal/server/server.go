package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"google.com/Dorfieeee/bootdev-http-protocol/internal/request"
	"google.com/Dorfieeee/bootdev-http-protocol/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		listener: listener,
		handler:  handler,
	}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}

	s.handler(w, req)
}
