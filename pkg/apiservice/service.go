package apiservice

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Service is the HTTP service
type Service struct {
	server *http.Server
}

// NewService cretes a service
func NewService(port string) *Service {
	h := newHandler()
	r := mux.NewRouter()

	r.HandleFunc("/api/create_hub", h.createHub)
	r.HandleFunc("/api/join_hub", h.joinHub)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	return &Service{
		server: server,
	}
}

// Serve starts the service
func (s *Service) Serve() error {
	return s.server.ListenAndServe()
}
