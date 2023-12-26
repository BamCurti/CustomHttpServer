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
	conn     net.Conn
	path     string
	method   string
	Headers  map[string]string
	response string
}

func (r HttpResponse) String() string {
	return fmt.Sprintf("path: %s\nMethod: %s\nHeaders: %s", r.path, r.method, r.Headers)
}

func NewResponse(c net.Conn) *HttpResponse {
	return &HttpResponse{
		conn:    c,
		Headers: map[string]string{},
	}
}

func (r *HttpResponse) FetchRequestInfo() error {
	buff := make([]byte, 1024)
	n, err := r.conn.Read(buff)
	if err != nil {
		return err
	}

	// getting path
	body := string(buff[:n])
	bodyLines := strings.Split(body, CRLF)
	startLine := strings.Split(bodyLines[0], " ")
	r.method = startLine[0]
	r.path = startLine[1]

	//getting Headers info
	for _, line := range bodyLines[1:] {
		if line == "" {
			break
		}

		parts := strings.Split(line, ": ")
		r.Headers[parts[0]] = parts[1]
	}

	return nil
}

func (r *HttpResponse) Handle() {
	defer r.conn.Close()

	err := r.FetchRequestInfo()
	if err != nil {
		r.buildResponse(err.Error(), INTERNAL_ERROR_MSG, TEXT_PLAIN)
		r.send()
	}

	var payload string
	var hc HttpCode
	var ct ContentType

	if r.path == "/" {
		payload = ""
		hc = OK_MSG
		ct = TEXT_PLAIN
	} else if strings.HasPrefix(r.path, "/echo") {
		payload = strings.TrimPrefix(r.path, "/echo/")
		hc = OK_MSG
		ct = TEXT_PLAIN
	} else if r.path == "/user-agent" {
		payload = r.Headers["User-Agent"]
		hc = OK_MSG
		ct = TEXT_PLAIN
	} else if strings.HasPrefix(r.path, "/files") {
		fileName := strings.TrimPrefix(r.path, "/files/")
		filePath := filepath.Join(*DirFlag, fileName)

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			payload = ""
			hc = NOT_FOUND_MSG
			ct = TEXT_PLAIN
		} else {
			payload = string(fileContent)
			hc = OK_MSG
			ct = APP_OCTET_STREAM
		}

	} else {
		payload = ""
		hc = NOT_FOUND_MSG
		ct = TEXT_PLAIN
	}

	r.buildResponse(payload, hc, ct)
	r.send()
}

func (r *HttpResponse) send() {
	_, err := r.conn.Write([]byte(r.response))
	if err != nil {
		log.Println("Failed to write data to connection", err.Error())
		return
	}
}

func (r *HttpResponse) buildResponse(payload string, statusCode HttpCode, ct ContentType) {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s %s%s", HTTP_1_1, statusCode, CRLF))
	builder.WriteString(fmt.Sprintf("Content-Type: %s%s", ct, CRLF))

	if payload != "" {
		builder.WriteString(ContentLength(payload))
		formatted := fmt.Sprintf("%s%s%s", CRLF, payload, CRLF)
		builder.WriteString(formatted)
	} else {
		builder.WriteString(END)
	}

	r.response = builder.String()
}
