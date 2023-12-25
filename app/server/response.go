package server

import (
	"fmt"
	"log"
	"net"
)

const (
	HTTP_1_1 = "HTTP/1.1"
	CRLF     = "\r\n"
	JUMP     = "\r\n\r\n"
	OK_MSG   = "200 OK"
)

type Response struct {
	conn net.Conn
}

func NewResponse(c net.Conn) *Response {
	return &Response{
		conn: c,
	}
}

func (r *Response) Handle() {
	defer r.conn.Close()

	rawData := []byte{}
	_, err := r.conn.Read(rawData)
	if err != nil {
		msg := "500 INTERNAL ERROR"
		r.send(msg)
	}

	msg := fmt.Sprintf("%s %s%s", HTTP_1_1, OK_MSG, JUMP)

	r.send(msg)
}

func (r *Response) send(msg string) {
	log.Println(msg)
	_, err := r.conn.Write([]byte(msg))
	if err != nil {
		log.Println("Failed to write data to connection", err.Error())
		return
	}
}
