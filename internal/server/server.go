package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	state    *atomic.Bool
	listener *net.Listener
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil {
		log.Fatalf("%s", err)
	}
	var started atomic.Bool
	started.Store(true)
	return &Server{
		state:    &started,
		listener: &listener,
	}, nil
}

func (s *Server) Close() error {
	return nil
}

func (s *Server) listen() {

}

func (s *Server) handle(conn net.Conn) {

}
