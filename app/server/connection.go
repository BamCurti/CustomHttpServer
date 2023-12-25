package server

import (
	"log"
	"net"
)

type HttpConnection struct {
	listener net.Listener
}

func NewHttpConnection(l net.Listener) HttpConnection {
	return HttpConnection{
		listener: l,
	}
}

func (c *HttpConnection) Accept() {
	conn, err := c.listener.Accept()
	if err != nil {
		log.Println(err.Error())
		return
	}

	r := NewResponse(conn)
	r.Handle()

}
