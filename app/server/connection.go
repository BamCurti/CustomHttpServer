package server

import (
	"net"
)

type HttpConnection struct {
	conn     net.Conn
	response *Response
}

func NewHttpConnection(conn net.Conn) HttpConnection {
	return HttpConnection{
		conn: conn,
	}
}

func (c *HttpConnection) Accept() {
	c.response = NewResponse(c.conn)
	c.response.Handle()
}
