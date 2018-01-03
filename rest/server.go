package rest

import (
	"fmt"
	"net/http"

	"context"

	"github.com/fxnn/deadbox/model"
)

type Server struct {
	addr    string
	server  *http.Server
	tls     TLS
	router  *router
	stopped chan error
}

func NewServer(addr string, tls TLS, drop model.Drop) *Server {
	return &Server{addr: addr, tls: tls, router: newRouter(drop)}
}

func (s *Server) Close() error {
	if s.server != nil {
		err := s.server.Shutdown(context.Background())
		s.server = nil
		if err != nil {
			<-s.stopped
			return err
		}
		return <-s.stopped
	}

	return nil
}

func (s *Server) StartServing() (err error) {
	s.server = &http.Server{Addr: s.addr, Handler: s.router}
	if s.server.TLSConfig, err = s.tls.Config(); err != nil {
		return
	}

	s.stopped = make(chan error)

	go func() {
		defer close(s.stopped)
		if err := s.tls.ListenAndServe(s.server); err != nil && err != http.ErrServerClosed {
			s.stopped <- fmt.Errorf("REST server on %s terminated unexpectedly: %s", s.addr, err)
		}
	}()

	return
}
