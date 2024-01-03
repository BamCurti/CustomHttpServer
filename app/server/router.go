package server

import (
	"errors"
	"log"
	"strings"
)

type Router struct {
	paths   map[string]*Router
	handler map[HttpMethod]RouteFunction
}

func NewRouter() *Router {

	return &Router{
		paths: make(map[string]*Router),
		handler: map[HttpMethod]RouteFunction{
			GET:  nil,
			POST: nil,
		},
	}
}

// Allows the user to add new handlers to different paths.
func (r *Router) Method(method HttpMethod, path string, handler RouteFunction) {
	if path == "/" {
		r.handler[method] = handler
		return
	}

	key, subpath, err := trimPath(path)
	log.Println(key, subpath)
	if err != nil {
		log.Panicln(err.Error())
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

func (r *Router) GetHandler(method HttpMethod, path string) RouteFunction {
	if path == "/" {
		return r.handler[method]
	}

	key, subpath, err := trimPath(path)
	if err != nil {
		return BadRequestHandler
	}

	subrouter, exists := r.paths[key]
	if !exists {
		return NotFoundHandler
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

//
