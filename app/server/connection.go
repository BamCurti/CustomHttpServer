package server

import (
	"net"
)

type HttpConnection struct {
	conn     net.Conn
	request  *HttpRequest
	response *HttpResponse
}

func NewHttpConnection(conn net.Conn) HttpConnection {
	return HttpConnection{
		conn: conn,
	}
}

func (c *HttpConnection) Accept(routes PathHandler) {
	request, _ := newHttpRequest(c.conn)
	c.response = NewResponse(c.conn)
	c.request = request
	c.response.Handle()
}
