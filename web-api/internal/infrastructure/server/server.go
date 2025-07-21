package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Start(handler http.Handler) error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("Server starting on %s...\n", addr)
	return http.ListenAndServe(addr, handler)
}

func (s *Server) Run(handler http.Handler) {
	if err := s.Start(handler); err != nil {
		log.Fatal(err)
	}
}