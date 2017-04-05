package rest

import (
	"github.com/fxnn/deadbox/model"
	"net/http"
)

type Server struct {
	addr   string
	router *router
}

func NewServer(addr string, drop model.Drop) *Server {
	return &Server{addr: addr, router: newRouter(drop)}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(s.addr, s.router)
}
