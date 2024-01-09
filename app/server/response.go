package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type HttpResponse struct {
	conn        net.Conn
	response    string
	StatusCode  HttpStatusCode
	ContentType ContentType
}

func NewResponse(c net.Conn) *HttpResponse {
	return &HttpResponse{
		conn: c,
	}
}

func (r *HttpResponse) Write(payload string) {
	r.buildResponse(payload)
	r.send()
}

func (r *HttpResponse) send() {
	_, err := r.conn.Write([]byte(r.response))
	if err != nil {
		log.Println("Failed to write data to connection", err.Error())
		return
	}
}

func (r *HttpResponse) buildResponse(payload string) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s %s%s", HTTP_1_1, r.StatusCode, CRLF))
	builder.WriteString(fmt.Sprintf("Content-Type: %s%s", r.ContentType, CRLF))

	if payload != "" {
		builder.WriteString(ContentLength(payload))
		formatted := fmt.Sprintf("%s%s%s", CRLF, payload, CRLF)
		builder.WriteString(formatted)
	} else {
		builder.WriteString(END)
	}

	r.response = builder.String()
}
