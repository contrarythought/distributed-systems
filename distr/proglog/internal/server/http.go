package server

import (
	"fmt"
	"net/http"
	"path"
)

// server that contains log of Records
type httpServer struct {
	Logs *Log
}

// map of client requests to a handler function
type pathResolver struct {
	handlers map[string]http.HandlerFunc
}

func newPathResolver() *pathResolver {
	return &pathResolver{make(map[string]http.HandlerFunc)}
}

// adds a handler to a path
func (p *pathResolver) Add(path string, handler http.HandlerFunc) {
	p.handlers[path] = handler
}

// iterates over handlers map, and calls the handler function associated with the path, if path and
// request method match
func (p *pathResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// method and path to check
	check := req.Method + " " + req.URL.Path

	// iterate across handlers and execute handler function
	for pattern, handlerFunc := range p.handlers {
		if ok, err := path.Match(pattern, check); ok && err != path.ErrBadPattern {
			handlerFunc(res, req)
			return
		} else if err == path.ErrBadPattern {
			fmt.Fprint(res, err)
		}
	}

	http.NotFound(res, req)
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Logs: NewLog(),
	}
}

// returns a server that handles set paths
func NewHTTPServer(addr string) *http.Server {
	s := newHTTPServer()
	p := newPathResolver()

	// add method and path
	p.Add("GET /", s.consumeHandle)
	p.Add("POST /", s.produceHandle)

	return &http.Server{
		Addr:    addr,
		Handler: p,
	}
}

// TODO
func (s *httpServer) consumeHandle(res http.ResponseWriter, req *http.Request) {

}

// TODO
func (s *httpServer) produceHandle(res http.ResponseWriter, req *http.Request) {

}
