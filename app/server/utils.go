package server

import "fmt"

type HttpStatusCode string
type ContentType string
type HttpMethod string

const (
	HTTP_1_1                          = "HTTP/1.1"
	CRLF                              = "\r\n"
	END                               = "\r\n\r\n"
	OK_MSG             HttpStatusCode = "200 OK"
	CREATED_MSG        HttpStatusCode = "201 Created"
	BAD_REQUEST_MSG    HttpStatusCode = "400 Bad Request"
	NOT_FOUND_MSG      HttpStatusCode = "404 Not Found"
	INTERNAL_ERROR_MSG HttpStatusCode = "500 Internal Error"
	TEXT_PLAIN         ContentType    = "text/plain"
	APP_OCTET_STREAM   ContentType    = "application/octet-stream"
	GET                HttpMethod     = "GET"
	POST               HttpMethod     = "POST"
)

// Variable to allows to read and send files from directory
var DirFlag *string

func ContentLength(payload string) string {
	return fmt.Sprintf("Content-Length: %d%s", len(payload), CRLF)
}
