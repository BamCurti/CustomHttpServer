package server

import (
	"net"
	"strings"
)

type HttpRequest struct {
	path    string
	method  HttpMethod
	Headers map[string]string
}

func newHttpRequest(conn net.Conn) (*HttpRequest, error) {
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)

	if err != nil {
		return nil, err
	}

	// fetch headers info
	body := string(buff[:n])
	bodyLines := strings.Split(body, CRLF)
	startLine := strings.Split(bodyLines[0], " ")
	method := startLine[0]
	path := startLine[1]
	headers := map[string]string{}

	for _, line := range bodyLines[1:] {
		if line == "" {
			break
		}

		parts := strings.Split(line, ": ")
		headers[parts[0]] = parts[1]
	}

	return &HttpRequest{
		path:    path,
		method:  HttpMethod(method),
		Headers: headers,
	}, nil
}
