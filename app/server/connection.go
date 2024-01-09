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

func (c *HttpConnection) Handle(r *Router) {
	defer c.Close()

	c.request, _ = newHttpRequest(c.conn)
	c.response = NewResponse(c.conn)

	method := c.request.method
	path := c.request.path

	f, urlParams := r.GetHandler(method, path)
	c.request.setUrlParams(urlParams)

	f(c.request, c.response)
}

func (c *HttpConnection) Close() {
	c.conn.Close()
}
