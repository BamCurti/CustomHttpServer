package server

import (
	"fmt"
	"log"
	"net"
)

type RouteFunction func(*HttpRequest, *HttpResponse)
type PathHandler map[string]map[HttpMethod]RouteFunction

type Server struct {
	network  string
	ip       string
	port     string
	listener net.Listener
	routes   map[string]map[HttpMethod]RouteFunction
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

func (s *Server) Serve() {
	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Println(err.Error())
			continue
		}

		httpConn := NewHttpConnection(conn)
		go httpConn.Accept(s.routes)
	}
}

func (s *Server) Get(path string, f RouteFunction) {
	s.routes[path][GET] = f
}

func (s *Server) Post(path string, f RouteFunction) {
	s.routes[path][POST] = f
}
