package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	network  string
	ip       string
	port     string
	listener net.Listener
}

func New(network, ip, port string) (*Server, error) {
	addr := fmt.Sprintf("%s:%s", ip, port)
	l, err := net.Listen(network, addr)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	s := &Server{
		network:  network,
		ip:       ip,
		port:     port,
		listener: l,
	}
	return s, nil
}

func (s *Server) Listen() {
	conn := NewHttpConnection(s.listener)
	conn.Accept()
}
