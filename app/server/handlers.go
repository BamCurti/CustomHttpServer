package server

func NotFoundHandler(req *HttpRequest, res *HttpResponse) {
	res.StatusCode = NOT_FOUND_MSG
	res.ContentType = TEXT_PLAIN
	res.Write("")
}

func internalErrorHandler(req *HttpRequest, res *HttpResponse) {
	res.StatusCode = INTERNAL_ERROR_MSG
	res.ContentType = TEXT_PLAIN
	res.Write("")
}

func BadRequestHandler(req *HttpRequest, res *HttpResponse) {
	res.StatusCode = BAD_REQUEST_MSG
	res.ContentType = TEXT_PLAIN
	res.Write("")
}
