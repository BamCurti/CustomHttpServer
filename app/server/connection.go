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

func (c *HttpConnection) Handle(paths PathHandler) {
	defer c.conn.Close()
	c.request, _ = newHttpRequest(c.conn)
	c.response = NewResponse(c.conn)
	f, found := paths[c.request.method][c.request.path]
	if !found {
		notFoundHandler(c.request, c.response)
	}

	f(c.request, c.response)
}

func (c *HttpConnection) Accept(routes PathHandler) {
	request, _ := newHttpRequest(c.conn)
	c.response = NewResponse(c.conn)
	c.request = request
	c.response.Handle()
}
