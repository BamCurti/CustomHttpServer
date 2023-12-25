package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type HttpCode string

const (
	HTTP_1_1                    = "HTTP/1.1"
	CRLF                        = "\r\n"
	JUMP                        = "\r\n\r\n"
	OK_MSG             HttpCode = "200 OK"
	NOT_FOUND_MSG      HttpCode = "404 Not Found"
	INTERNAL_ERROR_MSG HttpCode = "500 Internal Error"
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

	buff := make([]byte, 1024)
	n, err := r.conn.Read(buff)
	if err != nil {
		msg := buildResponse(err.Error(), INTERNAL_ERROR_MSG)
		r.send(msg)
	}

	body := string(buff[:n])
	bodyLines := strings.Split(body, "\n")
	startLine := strings.Split(bodyLines[0], " ")
	path := startLine[1]

	var response string

	if path == "/" {
		response = buildResponse("", OK_MSG)
	} else {
		response = buildResponse("", NOT_FOUND_MSG)
	}

	r.send(response)
}

func (r *Response) send(msg string) {
	_, err := r.conn.Write([]byte(msg))
	if err != nil {
		log.Println("Failed to write data to connection", err.Error())
		return
	}
}

func buildResponse(payload string, statusCode HttpCode) string {
	return fmt.Sprintf("%s %s%s%s", HTTP_1_1, statusCode, payload, JUMP)
}
