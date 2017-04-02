package drop

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	r.Path("/worker").
		HandlerFunc(s.handleGetAllWorkers).
		Methods("GET")
	r.Path("/worker/{workerId}").
		HandlerFunc(s.handlePutWorker).
		Methods("PUT")
	r.Path("/worker/{workerId}/request").
		HandlerFunc(s.handleGetAllWorkerRequests).
		Methods("GET")
	r.Path("/worker/{workerId}/request/{requestId}").
		HandlerFunc(s.handlePutWorkerRequest).
		Methods("PUT")
	return r
}

func (s *Server) handleGetAllWorkerRequests(http.ResponseWriter, *http.Request) {
	// FIXME Implement me!
}

func (s *Server) handlePutWorkerRequest(http.ResponseWriter,
	*http.Request) {
	// FIXME Implement me!
}

func (s *Server) handlePutWorker(http.ResponseWriter, *http.Request) {
	// FIXME Implement me!
}

func (s *Server) handleGetAllWorkers(http.ResponseWriter, *http.Request) {
	// FIXME Implement me!
}
