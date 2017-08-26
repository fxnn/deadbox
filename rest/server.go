package rest

import (
	"fmt"
	"net/http"

	"github.com/fxnn/deadbox/model"
)

type Server struct {
	addr    string
	server  *http.Server
	router  *router
	stopped chan error
}

func NewServer(addr string, drop model.Drop) *Server {
	return &Server{addr: addr, router: newRouter(drop)}
}

func (s *Server) Close() error {
	if s.server != nil {
		var err error = s.server.Shutdown(nil)
		s.server = nil
		if err != nil {
			<-s.stopped
			return err
		}
		return <-s.stopped
	}

	return nil
}

func (s *Server) StartServing() error {
	s.server = &http.Server{Addr: s.addr, Handler: s.router}
	s.stopped = make(chan error)

	go func() {
		defer close(s.stopped)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.stopped <- fmt.Errorf("REST server on %s terminated unexpectedly: %s", s.addr, err)
		}
	}()

	return nil
}
