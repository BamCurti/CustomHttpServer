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
	GET                string      = "GET"
)

type Response struct {
	conn    net.Conn
	path    string
	method  string
	Headers map[string]string
}

func (r Response) String() string {
	return fmt.Sprintf("path: %s\nMethod: %s\nHeaders: %s", r.path, r.method, r.Headers)
}

func NewResponse(c net.Conn) *Response {
	return &Response{
		conn:    c,
		Headers: map[string]string{},
	}
}

func (r *Response) FetchRequestInfo() error {
	buff := make([]byte, 1024)
	n, err := r.conn.Read(buff)
	if err != nil {
		return err
	}

	// getting path
	body := string(buff[:n])
	log.Println(body)
	bodyLines := strings.Split(body, CRLF)
	startLine := strings.Split(bodyLines[0], " ")
	r.method = startLine[0]
	r.path = startLine[1]

	//getting userAgent
	for _, line := range bodyLines[1:] {
		if line == "" {
			break
		}

		parts := strings.Split(line, ": ")
		r.Headers[parts[0]] = parts[1]

	}

	return nil
}

func (r *Response) Handle() {
	defer r.conn.Close()

	err := r.FetchRequestInfo()
	if err != nil {
		msg := buildResponse(err.Error(), INTERNAL_ERROR_MSG)
		r.send(msg)
	}

	var response string

	if r.path == "/" {
		response = buildResponse("", OK_MSG)
	} else if strings.HasPrefix(r.path, "/echo") {
		route := strings.TrimPrefix(r.path, "/echo/")
		response = buildResponse(route, OK_MSG)
	} else if r.path == "/user-agent" {
		response = buildResponse(r.Headers["User-Agent"], OK_MSG)
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
