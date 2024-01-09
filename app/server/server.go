package server

import (
	"fmt"
	"log"
	"net"
)

const tcp = "tcp"

type RouteFunction func(*HttpRequest, *HttpResponse)

type Server struct {
	network  string
	ip       string
	port     string
	listener net.Listener
	router   *Router
}

func New(ip, port string) (*Server, error) {
	addr := fmt.Sprintf("%s:%s", ip, port)
	l, err := net.Listen(tcp, addr)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	s := &Server{
		network:  tcp,
		ip:       ip,
		port:     port,
		listener: l,
		router:   NewRouter(),
	}
	return s, nil
}

func (s *Server) Serve() {
	defer s.listener.Close()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			log.Panicln(err.Error())
		}

		httpConn := NewHttpConnection(conn)
		go httpConn.Handle(s.router)
	}
}

func (s *Server) Get(path string, f RouteFunction) {
	r := s.router
	r.Get(path, f)
}

func (s *Server) Post(path string, f RouteFunction) {
	r := s.router
	r.Post(path, f)
}
