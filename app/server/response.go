package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type HttpResponse struct {
	conn        net.Conn
	path        string
	Headers     map[string]string
	response    string
	StatusCode  HttpStatusCode
	ContentType ContentType
}

func NewResponse(c net.Conn) *HttpResponse {
	return &HttpResponse{
		conn:    c,
		Headers: map[string]string{},
	}
}

func (r *HttpResponse) Handle() {
	defer r.conn.Close()

	var payload string

	if r.path == "/" {
		payload = ""
		r.StatusCode = OK_MSG
		r.ContentType = TEXT_PLAIN
	} else if strings.HasPrefix(r.path, "/echo") {
		payload = strings.TrimPrefix(r.path, "/echo/")
		r.StatusCode = OK_MSG
		r.ContentType = TEXT_PLAIN
	} else if r.path == "/user-agent" {
		payload = r.Headers["User-Agent"]
		r.StatusCode = OK_MSG
		r.ContentType = TEXT_PLAIN
	} else if strings.HasPrefix(r.path, "/files") {
		fileName := strings.TrimPrefix(r.path, "/files/")
		filePath := filepath.Join(*DirFlag, fileName)

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			payload = ""
			r.StatusCode = NOT_FOUND_MSG
			r.ContentType = TEXT_PLAIN
		} else {
			payload = string(fileContent)
			r.StatusCode = OK_MSG
			r.ContentType = APP_OCTET_STREAM
		}

	} else {
		payload = ""
		r.StatusCode = NOT_FOUND_MSG
		r.ContentType = TEXT_PLAIN
	}

	r.buildResponse(payload)
	r.send()
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
