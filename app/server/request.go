package server

import (
	"net"
	"strings"
)

type HttpRequest struct {
	path      string
	method    HttpMethod
	Headers   map[string]string
	urlParams map[string]string
	Content   string
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
	var initBody int

	// Populate the headers
	for idx, line := range bodyLines[1:] {
		initBody = idx + 1
		if line == "" {
			break
		}

		parts := strings.Split(line, ": ")
		headers[parts[0]] = parts[1]
	}

	// save the remaining message
	var content string
	if initBody < len(bodyLines) {
		content = strings.Join(bodyLines[initBody+1:], "\n")
	}

	return &HttpRequest{
		path:      path,
		method:    HttpMethod(method),
		Headers:   headers,
		urlParams: nil,
		Content:   content,
	}, nil
}

func (r *HttpRequest) setUrlParams(params map[string]string) {
	r.urlParams = params
}

func (r *HttpRequest) UrlParam(param string) (string, bool) {
	val, exists := r.urlParams[param]
	return val, exists
}
