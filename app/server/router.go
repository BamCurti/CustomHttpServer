package server

import (
	"errors"
	"log"
	"strings"
)

type HandlerMap map[HttpMethod]RouteFunction

type Router struct {
	paths           map[string]*Router
	handler         HandlerMap
	urlParamHandler HandlerMap
	urlParam        string
}

func NewRouter() *Router {

	return &Router{
		paths:           make(map[string]*Router),
		handler:         newHandlerMap(),
		urlParamHandler: newHandlerMap(),
	}
}

// Allows the user to add new handlers to different paths.
func (r *Router) Method(method HttpMethod, path string, handler RouteFunction) {
	if path == "/" {
		r.handler[method] = handler
		return
	}

	key, subpath, err := trimPath(path)
	if err != nil {
		log.Panicln(err.Error())
	}

	// subroute is a Url Parameter, needs to be added as urlParamHandler
	urlParam := getUrlParam(key)
	if urlParam != "" {
		r.urlParam = urlParam
		r.urlParamHandler[method] = handler
		return
	}

	subrouter, exist := r.paths[key]
	if !exist {
		subrouter = NewRouter()
		r.paths[key] = subrouter
	}

	subrouter.Method(method, subpath, handler)
}

func (r *Router) Get(path string, handler RouteFunction) {
	r.Method(GET, path, handler)
}

func (r *Router) Post(path string, handler RouteFunction) {
	r.Method(POST, path, handler)
}

func (r *Router) GetHandler(method HttpMethod, path string) (RouteFunction, map[string]string) {
	if path == "/" {
		return r.handler[method], nil
	}

	key, subpath, err := trimPath(path)
	if err != nil {
		return BadRequestHandler, nil
	}

	subrouter, exists := r.paths[key]
	if !exists {
		if r.urlParamHandler[method] == nil {
			return NotFoundHandler, nil
		}

		return r.urlParamHandler[method], map[string]string{
			r.urlParam: path[1:],
		}
	}

	return subrouter.GetHandler(method, subpath)
}

// Get the first part of the path, which will be the key of the possible handler, and the
// subpath, which might be another path or a query param.
//
// Returns an error if it does not follow the convention of API or does not wraps query params
// in brackets.
func trimPath(path string) (string, string, error) {
	if len(path) == 0 || path[0] != '/' {
		return "", "", errors.New("not valid path")
	}

	path = path[1:]
	idx := strings.IndexRune(path, '/')

	if idx == -1 {
		return path, "/", nil
	}

	return path[:idx], path[idx:], nil
}

func getUrlParam(path string) string {
	length := len(path)
	if path[0] != '{' || path[length-1] != '}' {
		return ""
	}

	return path[1 : length-1]
}

func newHandlerMap() HandlerMap {
	return map[HttpMethod]RouteFunction{
		GET:  nil,
		POST: nil,
	}
}
