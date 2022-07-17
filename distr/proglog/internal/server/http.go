package server

import (
	"encoding/json"
	"net/http"
	"path"
)

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func (s *httpServer) produceHandler(res http.ResponseWriter, req *http.Request) {
	var request ProduceRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	offset := s.Log.Append(request.Record)
	response := ProduceResponse{
		Offset: offset,
	}
	err = json.NewEncoder(res).Encode(&response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) consumeHandler(res http.ResponseWriter, req *http.Request) {
	var request ConsumeRequest
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := s.Log.Read(request.Offset)
	if err == ErrOutOfBounds {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	response := ConsumeResponse{
		Record: record,
	}
	err = json.NewEncoder(res).Encode(&response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}

// structs required to unmarshal json requests into
// request -> server
// server -> response
type ProduceRequest struct {
	Record Record `json:"record"` // append this record into the log
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"` // the response to caller that appended record
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"` // request to read record at offset
}

type ConsumeResponse struct {
	Record Record `json:"record"` // record to read
}

type pathResolver struct {
	handlers map[string]http.HandlerFunc
}

func newPathResolver() *pathResolver {
	return &pathResolver{
		handlers: make(map[string]http.HandlerFunc),
	}
}

func (p *pathResolver) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	check := req.Method + " " + req.URL.Path
	for reqPath, handler := range p.handlers {
		if matched, err := path.Match(check, reqPath); matched && err == nil {
			handler(res, req)
			return
		}
	}
	http.NotFound(res, req)
}

func (p *pathResolver) AddPath(path string, handler http.HandlerFunc) {
	p.handlers[path] = handler
}

func NewHTTPServer(addr string) *http.Server {
	server := newHTTPServer()
	paths := newPathResolver()
	paths.AddPath("GET /", server.consumeHandler)
	paths.AddPath("POST /", server.produceHandler)
	return &http.Server{
		Addr:    addr,
		Handler: paths,
	}
}
