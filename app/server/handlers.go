package server

func notFoundHandler(req *HttpRequest, res *HttpResponse) {
	res.StatusCode = NOT_FOUND_MSG
	res.ContentType = TEXT_PLAIN
	res.Write("")
}
