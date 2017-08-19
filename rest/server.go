package rest

import (
	"fmt"
	"github.com/fxnn/deadbox/model"
	"net"
	"net/http"
)

type Server struct {
	addr     string
	listener net.Listener
	router   *router
	stopped  chan error
}

func NewServer(addr string, drop model.Drop) *Server {
	return &Server{addr: addr, router: newRouter(drop)}
}

func (s *Server) Close() error {
	if s.listener != nil {
		var err error = s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
		return <-s.stopped
	}

	return nil
}

func (s *Server) StartServing() error {
	var (
		err error
	)

	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.stopped = make(chan error)
	go func() {
		defer close(s.stopped)
		if err := http.Serve(s.listener, s.router); err != nil {
			s.stopped <- fmt.Errorf("REST server on %s terminated: %s", s.addr, err)
		}
	}()
	return nil
}
