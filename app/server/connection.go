package server

import (
	"log"
	"net"
)

type HttpConnection struct {
	l net.Listener
}

func NewHttpConnection(l net.Listener) HttpConnection {
	return HttpConnection{
		l: l,
	}
}

func (c *HttpConnection) Accept() {
	_, err := c.l.Accept()
	if err != nil {
		log.Println(err.Error())
		return
	}

}
