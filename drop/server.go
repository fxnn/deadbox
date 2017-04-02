package drop

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe(s.addr, s.newRouter()))
}

func (s *Server) newRouter() *mux.Router {
	r := mux.NewRouter()
	r.Path("/queue/{workerId}").
		HandlerFunc(s.handleGetWorkerQueue).
		Methods("GET")
	r.Path("/queue/{workerId}/{requestId}").
		HandlerFunc(s.handlePutRequestIntoWorkerQueue).
		Methods("PUT")
	r.Path("/worker/{workerId}").
		HandlerFunc(s.handlePutWorker).
		Methods("PUT")
	r.Path("/worker").
		HandlerFunc(s.handleGetAllWorkers).
		Methods("GET")
	return r
}

func (s *Server) handleGetWorkerQueue(http.ResponseWriter, *http.Request) {
	// FIXME Implement me!
}

func (s *Server) handlePutRequestIntoWorkerQueue(http.ResponseWriter,
	*http.Request) {
	// FIXME Implement me!
}

func (s *Server) handlePutWorker(http.ResponseWriter, *http.Request) {
	// FIXME Implement me!
}

func (s *Server) handleGetAllWorkers(http.ResponseWriter, *http.Request) {
	// FIXME Implement me!
}
