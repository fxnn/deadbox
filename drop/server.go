package drop

import (
	"net/http"

	"github.com/gorilla/mux"
	"log"
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
	r.HandleFunc("/queue/{workerId}", s.handleGetWorkerQueue).
		Methods("GET")
	r.HandleFunc("/queue/{workerId}/{requestId}",
		s.handlePutRequestIntoWorkerQueue).Methods("PUT")
	r.HandleFunc("/worker/{workerId}", s.handlePutWorker).Methods("PUT")
	r.HandleFunc("/worker", s.handleGetAllWorkers).Methods("GET")
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
