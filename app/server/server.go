package server

import (
	"fmt"
	"log"
	"net"
)

const tcp = "tcp"

type RouteFunction func(*HttpRequest, *HttpResponse)
type PathHandler map[HttpMethod]map[string]RouteFunction

type Server struct {
	network  string
	ip       string
	port     string
	listener net.Listener
	routes   PathHandler
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
		routes:   NewPathHandler(),
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
		go httpConn.Handle(s.routes)
	}
}

func (s *Server) Get(path string, f RouteFunction) {
	s.routes[GET][path] = f
}

func (s *Server) Post(path string, f RouteFunction) {
	s.routes[POST][path] = f
}

func NewPathHandler() PathHandler {
	ph := make(PathHandler)
	ph[GET] = make(map[string]RouteFunction)
	ph[POST] = make(map[string]RouteFunction)
	return ph
}
