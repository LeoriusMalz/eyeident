package server

import (
	"net/http"
)

type Server struct {
	router http.Handler
}

func NewServer() *Server {
	s := &Server{}
	s.registerRoutes()
	return s
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
