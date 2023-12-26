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
	routes   map[string]func()
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
	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Println(err.Error())
			return
		}

		httpConn := NewHttpConnection(conn)
		go httpConn.Accept()
	}
}
