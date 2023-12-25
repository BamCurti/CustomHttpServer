package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type HttpCode string
type ContentType string

const (
	HTTP_1_1                       = "HTTP/1.1"
	CRLF                           = "\r\n"
	END                            = "\r\n\r\n"
	OK_MSG             HttpCode    = "200 OK"
	NOT_FOUND_MSG      HttpCode    = "404 Not Found"
	INTERNAL_ERROR_MSG HttpCode    = "500 Internal Error"
	TEXT_PLAIN         ContentType = "text/plain"
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
	} else if strings.HasPrefix(path, "/echo") {
		route := strings.TrimPrefix(path, "/echo/")
		response = buildResponse(route, OK_MSG)
	} else {
		response = buildResponse("", NOT_FOUND_MSG)
	}

	log.Println(response)

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
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s %s%s", HTTP_1_1, statusCode, CRLF))
	builder.WriteString(fmt.Sprintf("Content-Type: %s%s", TEXT_PLAIN, CRLF))

	if payload != "" {
		builder.WriteString(ContentLength(payload))
		formatted := fmt.Sprintf("%s%s%s", CRLF, payload, CRLF)
		builder.WriteString(formatted)
	} else {
		builder.WriteString(END)
	}

	return builder.String()
}

func ContentLength(payload string) string {
	return fmt.Sprintf("Content-Length: %d%s", len(payload), CRLF)
}
