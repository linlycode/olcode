package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Service is the HTTP service
type Service interface {
	Serve() error
}

type service struct {
	server *http.Server
}

// NewService cretes a service
func NewService(port string) Service {
	h := newHandler()
	r := mux.NewRouter()

	r.HandleFunc("/ws", h.serveWS)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	return &service{
		server: server,
	}
}

// Serve starts the service
func (s *service) Serve() error {
	return s.server.ListenAndServe()
}
